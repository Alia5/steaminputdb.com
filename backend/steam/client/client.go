package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log/slog"
	mrand "math/rand/v2"
	"net/http"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	ProtocolVersion uint32 = 65580
	OSTypeLinux     uint32 = 203

	encryptProtoVersion uint32 = 1
	encryptKeySize      uint32 = 128
	sessionKeyLen              = 32
	eUniversePublic     uint32 = 1
	eResultOK           int32  = 1
	defaultHeartbeat    int32  = 5
	cmListURL                  = "https://api.steampowered.com/ISteamDirectory/GetCMList/v1/?cellId=0"

	steamIDAnon       = uint64(1)<<56 | uint64(10)<<52
	steamIDIndividual = uint64(1)<<56 | uint64(1)<<52 | uint64(1)<<32
)

var (
	ErrUnexpectedUniverse = errors.New("unexpected universe")
	ErrEncryptionFailed   = errors.New("encryption failed")
	ErrLoginFailed        = errors.New("login failed")
	ErrNoCMServers        = errors.New("no CM servers returned")
	ErrDisconnected       = errors.New("disconnected")
)

type Client interface {
	Connect(ctx context.Context) error
	Login(ctx context.Context, details LoginDetails) error
	Disconnect()
	SendMessage(ctx context.Context, reqEMsg, respEMsg EMsg, req, resp proto.Message) error
	Connected() bool
	LoggedIn() bool
}

type client struct {
	conn      *connection
	cancel    context.CancelFunc
	mu        sync.Mutex
	heartbeat *time.Ticker
	done      chan struct{}
	encrypted chan error
	pending   map[EMsg]chan *packet
	pendingMu sync.Mutex
	loggedIn  bool
}

func New() Client {
	return &client{
		done:      make(chan struct{}),
		encrypted: make(chan error, 1),
		pending:   make(map[EMsg]chan *packet),
	}
}

func (c *client) Connect(ctx context.Context) error {
	c.mu.Lock()
	c.done = make(chan struct{})
	c.encrypted = make(chan error, 1)
	c.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)

	addr, err := getRandomCM(ctx)
	if err != nil {
		cancel()
		return fmt.Errorf("discover CM: %w", err)
	}
	slog.Info("connecting to CM", "addr", addr)
	conn, err := dial(addr)
	if err != nil {
		cancel()
		return fmt.Errorf("dial: %w", err)
	}

	c.mu.Lock()
	c.cancel = cancel
	c.conn = conn
	c.mu.Unlock()

	go c.readLoop()
	go func() {
		select {
		case <-ctx.Done():
			c.Disconnect()
		case <-c.done:
		}
	}()
	select {
	case err := <-c.encrypted:
		if err != nil {
			c.Disconnect()
			return err
		}
	case <-ctx.Done():
		c.Disconnect()
		return ctx.Err()
	}
	return nil
}

type LoginDetails struct {
	Anonymous    bool
	RefreshToken string
	Language     string
}

func (c *client) Login(ctx context.Context, details LoginDetails) error {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()
	if conn == nil {
		return ErrDisconnected
	}
	req := &CMsgClientLogon{
		ProtocolVersion: proto.Uint32(ProtocolVersion),
		ClientOsType:    proto.Uint32(OSTypeLinux),
		ClientLanguage:  &details.Language,
	}
	if details.Anonymous {
		conn.steamID = steamIDAnon
		req.AnonUserTargetAccountName = proto.String("anonymous")
	} else {
		conn.steamID = steamIDIndividual
		req.AccessToken = &details.RefreshToken
		req.ShouldRememberPassword = proto.Bool(true)
	}
	pkt, err := c.sendAndWait(ctx, EMsg_k_EMsgClientLogon, req, EMsg_k_EMsgClientLogOnResponse)
	if err != nil {
		return err
	}
	body := &CMsgClientLogonResponse{}
	if err := proto.Unmarshal(pkt.body, body); err != nil {
		return fmt.Errorf("unmarshal logon response: %w", err)
	}
	if body.GetEresult() != eResultOK {
		return fmt.Errorf("%w: eresult=%d", ErrLoginFailed, body.GetEresult())
	}
	c.mu.Lock()
	if c.conn != nil {
		c.conn.sessionID = pkt.header.GetClientSessionid()
		c.conn.steamID = pkt.header.GetSteamid()
	}
	c.loggedIn = true
	c.mu.Unlock()
	go c.heartbeatLoop(time.Duration(max(body.GetHeartbeatSeconds(), defaultHeartbeat)) * time.Second)
	return nil
}

func (c *client) Disconnect() {
	c.mu.Lock()
	cancel := c.cancel
	c.cancel = nil
	c.loggedIn = false
	if c.heartbeat != nil {
		c.heartbeat.Stop()
		c.heartbeat = nil
	}
	if c.conn != nil {
		c.conn.close()
		c.conn = nil
	}
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	c.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (c *client) Connected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn != nil
}

func (c *client) LoggedIn() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.loggedIn
}

func (c *client) SendMessage(ctx context.Context, reqEMsg, respEMsg EMsg, req, resp proto.Message) error {
	pkt, err := c.sendAndWait(ctx, reqEMsg, req, respEMsg)
	if err != nil {
		return err
	}
	return proto.Unmarshal(pkt.body, resp)
}

func (c *client) sendAndWait(ctx context.Context, reqEMsg EMsg, req proto.Message, respEMsg EMsg) (*packet, error) {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()
	if conn == nil {
		return nil, ErrDisconnected
	}
	ch := make(chan *packet, 1)
	c.pendingMu.Lock()
	c.pending[respEMsg] = ch
	c.pendingMu.Unlock()
	if err := conn.writeProtoMsg(reqEMsg, req); err != nil {
		c.pendingMu.Lock()
		delete(c.pending, respEMsg)
		c.pendingMu.Unlock()
		return nil, err
	}
	select {
	case pkt := <-ch:
		return pkt, nil
	case <-ctx.Done():
		c.pendingMu.Lock()
		delete(c.pending, respEMsg)
		c.pendingMu.Unlock()
		return nil, ctx.Err()
	case <-c.done:
		return nil, ErrDisconnected
	}
}

func (c *client) readLoop() {
	defer c.Disconnect()
	for {
		c.mu.Lock()
		conn := c.conn
		c.mu.Unlock()
		if conn == nil {
			return
		}
		data, err := conn.read()
		if err != nil {
			return
		}
		if err := c.handlePacket(data); err != nil {
			slog.Error("packet error", "err", err)
			return
		}
	}
}

func (c *client) handlePacket(data []byte) error {
	pkt, err := parsePacket(data)
	if err != nil {
		return err
	}
	switch pkt.eMsg {
	case EMsg_k_EMsgChannelEncryptRequest:
		return c.handleEncryptRequest(pkt)
	case EMsg_k_EMsgChannelEncryptResult:
		return c.handleEncryptResult(pkt)
	case EMsg_k_EMsgMulti:
		return c.handleMulti(pkt)
	default:
		slog.Info("received message", "eMsg", pkt.eMsg)
		c.pendingMu.Lock()
		ch, ok := c.pending[pkt.eMsg]
		if ok {
			delete(c.pending, pkt.eMsg)
		}
		c.pendingMu.Unlock()
		if ok {
			ch <- pkt
		}
	}
	return nil
}

func (c *client) handleEncryptRequest(pkt *packet) error {
	r := bytes.NewReader(pkt.body)
	var skip, universe uint32
	binary.Read(r, binary.LittleEndian, &skip)
	binary.Read(r, binary.LittleEndian, &universe)
	if universe != eUniversePublic {
		return fmt.Errorf("%w: %d", ErrUnexpectedUniverse, universe)
	}
	sessionKey := make([]byte, sessionKeyLen)
	if _, err := rand.Read(sessionKey); err != nil {
		return err
	}
	c.mu.Lock()
	c.conn.tempSessionKey = sessionKey
	c.mu.Unlock()
	enc, err := rsaEncryptSessionKey(sessionKey)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint32(EMsg_k_EMsgChannelEncryptResponse))
	binary.Write(buf, binary.LittleEndian, jobIDNone)
	binary.Write(buf, binary.LittleEndian, jobIDNone)
	binary.Write(buf, binary.LittleEndian, encryptProtoVersion)
	binary.Write(buf, binary.LittleEndian, encryptKeySize)
	buf.Write(enc)
	binary.Write(buf, binary.LittleEndian, crc32.ChecksumIEEE(enc))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	return c.conn.write(buf.Bytes())
}

func (c *client) handleEncryptResult(pkt *packet) error {
	var result int32
	binary.Read(bytes.NewReader(pkt.body), binary.LittleEndian, &result)
	if result != eResultOK {
		err := fmt.Errorf("%w: %d", ErrEncryptionFailed, result)
		select {
		case c.encrypted <- err:
		default:
		}
		return err
	}
	c.mu.Lock()
	c.conn.setEncryptionKey(c.conn.tempSessionKey)
	c.conn.tempSessionKey = nil
	c.mu.Unlock()
	select {
	case c.encrypted <- nil:
	default:
	}
	return nil
}

func (c *client) heartbeatLoop(interval time.Duration) {
	c.mu.Lock()
	if c.heartbeat != nil {
		c.heartbeat.Stop()
	}
	c.heartbeat = time.NewTicker(interval)
	c.mu.Unlock()
	for {
		select {
		case <-c.done:
			return
		case <-c.heartbeat.C:
			c.mu.Lock()
			conn := c.conn
			c.mu.Unlock()
			if conn == nil {
				return
			}
			if err := conn.writeProtoMsg(EMsg_k_EMsgClientHeartBeat, &CMsgClientHeartBeat{}); err != nil {
				slog.Error("heartbeat failed", "err", err)
				return
			}
		}
	}
}

func (c *client) handleMulti(pkt *packet) error {
	body := &CMsgMulti{}
	if err := proto.Unmarshal(pkt.body, body); err != nil {
		return err
	}
	payload := body.GetMessageBody()
	if body.GetSizeUnzipped() > 0 {
		r, err := gzip.NewReader(bytes.NewReader(payload))
		if err != nil {
			return err
		}
		if payload, err = io.ReadAll(r); err != nil {
			return err
		}
	}
	for pr := bytes.NewReader(payload); pr.Len() > 0; {
		var length uint32
		if err := binary.Read(pr, binary.LittleEndian, &length); err != nil {
			return err
		}
		sub := make([]byte, length)
		if _, err := io.ReadFull(pr, sub); err != nil {
			return err
		}
		if err := c.handlePacket(sub); err != nil {
			slog.Warn("multi sub-packet error", "err", err)
		}
	}
	return nil
}

type cmListResponse struct {
	Response cmListResponseData `json:"response"`
}
type cmListResponseData struct {
	ServerList []string `json:"serverlist"`
	Result     uint32   `json:"result"`
}

func getRandomCM(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cmListURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result cmListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Response.Result != 1 || len(result.Response.ServerList) == 0 {
		return "", ErrNoCMServers
	}
	return result.Response.ServerList[mrand.IntN(len(result.Response.ServerList))], nil
}

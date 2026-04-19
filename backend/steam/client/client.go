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
	"sync/atomic"
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
	SendMessage(ctx context.Context, reqEMsg EMsg, req, resp proto.Message) error
	Connected() bool
	LoggedIn() bool
	EnableAutoReconnect(details LoginDetails, timeout time.Duration)
}

type client struct {
	conn        *connection
	cancel      context.CancelFunc
	mu          sync.Mutex
	done        chan struct{}
	encrypted   chan error
	pending     map[uint64]chan *packet
	pendingEMsg map[EMsg]chan *packet
	pendingMu   sync.Mutex
	nextJobID   atomic.Uint64
	loggedIn    bool

	reconnectMu   sync.Mutex
	autoReconnect atomic.Bool
	loginDetails  LoginDetails
	reconnTimeout time.Duration
}

func New() Client {
	c := &client{
		done:        make(chan struct{}),
		encrypted:   make(chan error, 1),
		pending:     make(map[uint64]chan *packet),
		pendingEMsg: make(map[EMsg]chan *packet),
	}
	c.nextJobID.Store(1)
	return c
}

func (c *client) Connect(ctx context.Context) error {
	c.mu.Lock()
	if c.conn != nil {
		c.mu.Unlock()
		return nil
	}
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
	done := c.done
	c.mu.Unlock()

	go c.readLoop(conn)
	go func() {
		select {
		case <-ctx.Done():
			c.Disconnect()
		case <-done:
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
		ProtocolVersion: new(ProtocolVersion),
		ClientOsType:    new(OSTypeLinux),
		ClientLanguage:  &details.Language,
	}
	if details.Anonymous {
		conn.steamID = steamIDAnon
		req.AnonUserTargetAccountName = new("anonymous")
	} else {
		conn.steamID = steamIDIndividual
		req.AccessToken = &details.RefreshToken
		req.ShouldRememberPassword = new(true)
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
	conn = c.conn
	done := c.done
	c.mu.Unlock()
	if conn != nil {
		go c.heartbeatLoop(conn, done, time.Duration(max(body.GetHeartbeatSeconds(), defaultHeartbeat))*time.Second)
	}
	return nil
}

func (c *client) Disconnect() {
	c.mu.Lock()
	cancel := c.cancel
	c.cancel = nil
	c.loggedIn = false
	if c.conn != nil {
		_ = c.conn.close()
		c.conn = nil
	}
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	c.pendingMu.Lock()
	clear(c.pending)
	clear(c.pendingEMsg)
	c.pendingMu.Unlock()
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

func (c *client) SendMessage(ctx context.Context, reqEMsg EMsg, req, resp proto.Message) error {
	pkt, err := c.sendAndWait(ctx, reqEMsg, req)
	if errors.Is(err, ErrDisconnected) {
		if c.autoReconnect.Load() {
			if rerr := c.doReconnect(ctx); rerr != nil {
				return rerr
			}
			pkt, err = c.sendAndWait(ctx, reqEMsg, req)
		}
	}
	if err != nil {
		return err
	}
	return proto.Unmarshal(pkt.body, resp)
}

func (c *client) EnableAutoReconnect(details LoginDetails, timeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.loginDetails = details
	c.reconnTimeout = timeout
	c.autoReconnect.Store(true)
}

func (c *client) doReconnect(ctx context.Context) error {
	c.reconnectMu.Lock()
	defer c.reconnectMu.Unlock()
	if c.Connected() {
		return nil
	}
	c.mu.Lock()
	details := c.loginDetails
	timeout := c.reconnTimeout
	c.mu.Unlock()
	rcCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	c.Disconnect()
	if err := c.Connect(rcCtx); err != nil {
		return fmt.Errorf("reconnect: %w", err)
	}
	return c.Login(rcCtx, details)
}

func (c *client) sendAndWait(ctx context.Context, reqEMsg EMsg, req proto.Message, respEMsg ...EMsg) (*packet, error) {
	c.mu.Lock()
	conn := c.conn
	done := c.done
	c.mu.Unlock()
	if conn == nil {
		return nil, ErrDisconnected
	}
	jobID := c.nextJobID.Add(1)
	ch := make(chan *packet, 1)
	c.pendingMu.Lock()
	c.pending[jobID] = ch
	for _, em := range respEMsg {
		c.pendingEMsg[em] = ch
	}
	c.pendingMu.Unlock()
	cleanup := func() {
		c.pendingMu.Lock()
		delete(c.pending, jobID)
		for _, em := range respEMsg {
			delete(c.pendingEMsg, em)
		}
		c.pendingMu.Unlock()
	}
	if err := conn.writeProtoMsg(reqEMsg, req, jobID); err != nil {
		cleanup()
		return nil, err
	}
	select {
	case pkt := <-ch:
		return pkt, nil
	case <-ctx.Done():
		cleanup()
		return nil, ctx.Err()
	case <-done:
		cleanup()
		return nil, ErrDisconnected
	}
}

func (c *client) readLoop(conn *connection) {
	defer func() {
		c.mu.Lock()
		stillOwner := c.conn == conn
		c.mu.Unlock()
		if stillOwner {
			c.Disconnect()
		}
	}()
	for {
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
		jobID := pkt.header.GetJobidTarget()
		c.pendingMu.Lock()
		var ch chan *packet
		var ok bool
		if jobID != jobIDNone {
			ch, ok = c.pending[jobID]
			if ok {
				delete(c.pending, jobID)
			}
		}
		if !ok {
			ch, ok = c.pendingEMsg[pkt.eMsg]
			if ok {
				delete(c.pendingEMsg, pkt.eMsg)
			}
		}
		c.pendingMu.Unlock()
		if ok {
			ch <- pkt
		} else {
			slog.Debug("unhandled message", "eMsg", pkt.eMsg, "jobID", jobID)
		}
	}
	return nil
}

func (c *client) handleEncryptRequest(pkt *packet) error {
	r := bytes.NewReader(pkt.body)
	var skip, universe uint32
	_ = binary.Read(r, binary.LittleEndian, &skip)
	_ = binary.Read(r, binary.LittleEndian, &universe)
	if universe != eUniversePublic {
		return fmt.Errorf("%w: %d", ErrUnexpectedUniverse, universe)
	}
	sessionKey := make([]byte, sessionKeyLen)
	if _, err := rand.Read(sessionKey); err != nil {
		return err
	}
	c.mu.Lock()
	conn := c.conn
	if conn == nil {
		c.mu.Unlock()
		return ErrDisconnected
	}
	conn.tempSessionKey = sessionKey
	c.mu.Unlock()
	enc, err := rsaEncryptSessionKey(sessionKey)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, uint32(EMsg_k_EMsgChannelEncryptResponse))
	_ = binary.Write(buf, binary.LittleEndian, jobIDNone)
	_ = binary.Write(buf, binary.LittleEndian, jobIDNone)
	_ = binary.Write(buf, binary.LittleEndian, encryptProtoVersion)
	_ = binary.Write(buf, binary.LittleEndian, encryptKeySize)
	buf.Write(enc)
	_ = binary.Write(buf, binary.LittleEndian, crc32.ChecksumIEEE(enc))
	_ = binary.Write(buf, binary.LittleEndian, uint32(0))
	return conn.write(buf.Bytes())
}

func (c *client) handleEncryptResult(pkt *packet) error {
	var result int32
	_ = binary.Read(bytes.NewReader(pkt.body), binary.LittleEndian, &result)
	if result != eResultOK {
		err := fmt.Errorf("%w: %d", ErrEncryptionFailed, result)
		select {
		case c.encrypted <- err:
		default:
		}
		return err
	}
	c.mu.Lock()
	conn := c.conn
	if conn == nil {
		c.mu.Unlock()
		return ErrDisconnected
	}
	conn.setEncryptionKey(conn.tempSessionKey)
	conn.tempSessionKey = nil
	c.mu.Unlock()
	select {
	case c.encrypted <- nil:
	default:
	}
	return nil
}

func (c *client) heartbeatLoop(conn *connection, done <-chan struct{}, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
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
	defer func() { _ = resp.Body.Close() }()
	var result cmListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Response.Result != 1 || len(result.Response.ServerList) == 0 {
		return "", ErrNoCMServers
	}
	return result.Response.ServerList[mrand.IntN(len(result.Response.ServerList))], nil
}

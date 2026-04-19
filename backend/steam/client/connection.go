package client

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	tcpMagic   uint32 = 0x31305456
	protoMask  uint32 = 0x80000000
	eMsgMask   uint32 = 0x7FFFFFFF
	msgHdrSize        = 20
	jobIDNone  uint64 = ^uint64(0)
)

var (
	errInvalidTCPMagic = errors.New("invalid tcp magic")
	errPacketTooShort  = errors.New("packet too short")
	errHeaderExceeds   = errors.New("proto header exceeds packet")
)

type packet struct {
	eMsg   EMsg
	header *CMsgProtoBufHeader
	body   []byte
}

func parsePacket(data []byte) (*packet, error) {
	if len(data) < 4 {
		return nil, errPacketTooShort
	}
	raw := binary.LittleEndian.Uint32(data[:4])
	eMsg := EMsg(raw & eMsgMask)
	if raw&protoMask == 0 {
		if len(data) < msgHdrSize {
			return nil, errPacketTooShort
		}
		return &packet{eMsg: eMsg, body: data[msgHdrSize:]}, nil
	}
	if len(data) < 8 {
		return nil, errPacketTooShort
	}
	hLen := binary.LittleEndian.Uint32(data[4:8])
	if len(data) < int(8+hLen) {
		return nil, errHeaderExceeds
	}
	hdr := &CMsgProtoBufHeader{}
	if err := proto.Unmarshal(data[8:8+hLen], hdr); err != nil {
		return nil, fmt.Errorf("%w: %w", errHeaderExceeds, err)
	}
	return &packet{eMsg: eMsg, header: hdr, body: data[8+hLen:]}, nil
}

type connection struct {
	conn           *net.TCPConn
	ciph           cipher.Block
	mu             sync.RWMutex
	writeMu        sync.Mutex
	tempSessionKey []byte
	sessionID      int32
	steamID        uint64
}

func dial(addr string) (*connection, error) {
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return nil, err
	}
	return &connection{conn: conn.(*net.TCPConn)}, nil
}

func (c *connection) read() ([]byte, error) {
	var hdr [8]byte
	if _, err := io.ReadFull(c.conn, hdr[:]); err != nil {
		return nil, err
	}
	pLen := binary.LittleEndian.Uint32(hdr[:4])
	if binary.LittleEndian.Uint32(hdr[4:]) != tcpMagic {
		return nil, fmt.Errorf("%w: %x", errInvalidTCPMagic, hdr[4:])
	}
	buf := make([]byte, pLen)
	if _, err := io.ReadFull(c.conn, buf); err != nil {
		return nil, err
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.ciph != nil {
		buf = symmetricDecrypt(c.ciph, buf)
	}
	return buf, nil
}

func (c *connection) write(data []byte) error {
	c.mu.RLock()
	if c.ciph != nil {
		var err error
		if data, err = symmetricEncrypt(c.ciph, data); err != nil {
			c.mu.RUnlock()
			return err
		}
	}
	c.mu.RUnlock()
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, uint32(len(data)))
	_ = binary.Write(buf, binary.LittleEndian, tcpMagic)
	buf.Write(data)
	c.writeMu.Lock()
	_, err := c.conn.Write(buf.Bytes())
	c.writeMu.Unlock()
	return err
}

func (c *connection) setEncryptionKey(key []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ciph, _ = aes.NewCipher(key)
}

func (c *connection) writeProtoMsg(eMsg EMsg, body proto.Message, jobID ...uint64) error {
	hdr := &CMsgProtoBufHeader{}
	if len(jobID) > 0 {
		hdr.JobidSource = &jobID[0]
	}
	if c.sessionID != 0 {
		hdr.ClientSessionid = new(c.sessionID)
	}
	if c.steamID != 0 {
		hdr.Steamid = new(c.steamID)
	}
	var bodyBytes []byte
	if body != nil {
		var err error
		if bodyBytes, err = proto.Marshal(body); err != nil {
			return err
		}
	}
	hdrBytes, err := proto.Marshal(hdr)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, uint32(eMsg)|protoMask)
	_ = binary.Write(buf, binary.LittleEndian, uint32(len(hdrBytes)))
	buf.Write(hdrBytes)
	buf.Write(bodyBytes)
	return c.write(buf.Bytes())
}

func (c *connection) close() error { return c.conn.Close() }

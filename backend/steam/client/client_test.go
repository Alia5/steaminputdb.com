package client

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func TestIsDisconnectError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "disconnected sentinel", err: ErrDisconnected, want: true},
		{name: "wrapped disconnected", err: fmt.Errorf("wrapped: %w", ErrDisconnected), want: true},
		{name: "net err closed", err: net.ErrClosed, want: true},
		{name: "wrapped net err closed", err: fmt.Errorf("write failed: %w", net.ErrClosed), want: true},
		{name: "io eof", err: io.EOF, want: true},
		{name: "epipe", err: syscall.EPIPE, want: true},
		{name: "econnreset", err: syscall.ECONNRESET, want: true},
		{name: "string closed conn", err: errors.New("write tcp 1.1.1.1:1234->2.2.2.2:27017: use of closed network connection"), want: true},
		{name: "string broken pipe", err: errors.New("write: broken pipe"), want: true},
		{name: "random error", err: errors.New("something else"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDisconnectError(tt.err); got != tt.want {
				t.Fatalf("isDisconnectError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestSteamClient_OutageAndRecoveryOnSameClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Steam integration test in short mode")
	}

	sc, ok := New().(*client)
	if !ok {
		t.Fatalf("client init failed")
	}
	t.Cleanup(sc.Disconnect)

	details := LoginDetails{Anonymous: true, Language: "english"}
	sc.EnableAutoReconnect(details, 10*time.Second)

	{
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := sc.Connect(ctx); err != nil {
			t.Fatalf("connect to steam failed: %v", err)
		}
		if err := sc.Login(ctx, details); err != nil {
			t.Fatalf("connect to steam failed: %v", err)
		}
	}

	{
		deadline := time.Now().Add(20 * time.Second)
		success := 0
		var lastErr error
		for success < 2 && time.Now().Before(deadline) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			appid := uint32(570)
			var resp CMsgClientPICSProductInfoResponse
			err := sc.SendMessage(
				ctx,
				EMsg_k_EMsgClientPICSProductInfoRequest,
				&CMsgClientPICSProductInfoRequest{
					Apps: []*CMsgClientPICSProductInfoRequest_AppInfo{{Appid: &appid}},
				},
				&resp,
			)
			cancel()
			if err != nil {
				lastErr = err
				time.Sleep(200 * time.Millisecond)
				continue
			}
			success++
		}
		if success < 2 {
			t.Fatalf("getAppInfo failed: %v", lastErr)
		}
	}

	sc.mu.Lock()
	conn := sc.conn
	sc.mu.Unlock()
	if conn == nil || conn.conn == nil {
		t.Fatalf("connect to steam failed: %v", ErrDisconnected)
	}
	if err := conn.conn.SetWriteDeadline(time.Now().Add(-1 * time.Second)); err != nil {
		t.Fatalf("failed to inject transport timeout: %v", err)
	}

	for range 3 {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		appid := uint32(570)
		var resp CMsgClientPICSProductInfoResponse
		err := sc.SendMessage(
			ctx,
			EMsg_k_EMsgClientPICSProductInfoRequest,
			&CMsgClientPICSProductInfoRequest{
				Apps: []*CMsgClientPICSProductInfoRequest_AppInfo{{Appid: &appid}},
			},
			&resp,
		)
		cancel()
		if err == nil {
			t.Fatalf("getAppInfo should fail")
		}
		var netErr net.Error
		if !(errors.As(err, &netErr) && netErr.Timeout()) {
			t.Fatalf("timeout error expected: %v", err)
		}
	}

	if err := conn.conn.SetWriteDeadline(time.Time{}); err != nil {
		t.Fatalf("failed to clear transport timeout: %v", err)
	}

	{
		deadline := time.Now().Add(30 * time.Second)
		success := 0
		var lastErr error
		for success < 3 && time.Now().Before(deadline) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			appid := uint32(570)
			var resp CMsgClientPICSProductInfoResponse
			err := sc.SendMessage(
				ctx,
				EMsg_k_EMsgClientPICSProductInfoRequest,
				&CMsgClientPICSProductInfoRequest{
					Apps: []*CMsgClientPICSProductInfoRequest_AppInfo{{Appid: &appid}},
				},
				&resp,
			)
			cancel()
			if err != nil {
				lastErr = err
				time.Sleep(200 * time.Millisecond)
				continue
			}
			success++
		}
		if success < 3 {
			t.Fatalf("getAppInfo failed: %v", lastErr)
		}
	}
}

func TestDoReconnectDoesNotDisconnectOnInternalTimeoutCancel(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()

		encryptResult := bytes.NewBuffer(nil)
		if err := binary.Write(encryptResult, binary.LittleEndian, int32(eResultOK)); err != nil {
			return
		}
		if err := writeProtoBodyPacket(conn, EMsg_k_EMsgChannelEncryptResult, encryptResult.Bytes()); err != nil {
			return
		}

		for {
			payload, err := readFramedPacket(conn)
			if err != nil {
				return
			}
			pkt, err := parsePacket(payload)
			if err != nil {
				return
			}
			if pkt.eMsg == EMsg_k_EMsgClientLogon {
				if err := writeProtoPacket(conn, EMsg_k_EMsgClientLogOnResponse, &CMsgProtoBufHeader{}, &CMsgClientLogonResponse{Eresult: new(eResultOK), HeartbeatSeconds: new(defaultHeartbeat)}); err != nil {
					return
				}
			}
		}
	}()
	t.Cleanup(wg.Wait)

	oldDiscover := discoverCM
	discoverCM = func(context.Context) (string, error) {
		return ln.Addr().String(), nil
	}
	t.Cleanup(func() { discoverCM = oldDiscover })

	sc, ok := New().(*client)
	if !ok {
		t.Fatalf("client init failed")
	}
	t.Cleanup(sc.Disconnect)

	details := LoginDetails{Anonymous: true, Language: "english"}
	sc.EnableAutoReconnect(details, 2*time.Second)

	if err := sc.doReconnect(context.Background()); err != nil {
		t.Fatalf("doReconnect failed: %v", err)
	}

	if !sc.Connected() {
		t.Fatalf("client should remain connected after doReconnect")
	}
	if !sc.LoggedIn() {
		t.Fatalf("client should remain logged in after doReconnect")
	}
}

func readFramedPacket(conn net.Conn) ([]byte, error) {
	var hdr [8]byte
	if _, err := io.ReadFull(conn, hdr[:]); err != nil {
		return nil, err
	}
	pLen := binary.LittleEndian.Uint32(hdr[:4])
	if binary.LittleEndian.Uint32(hdr[4:]) != tcpMagic {
		return nil, errInvalidTCPMagic
	}
	payload := make([]byte, pLen)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func writeProtoPacket(conn net.Conn, eMsg EMsg, hdr *CMsgProtoBufHeader, body proto.Message) error {
	hdrBytes, err := proto.Marshal(hdr)
	if err != nil {
		return err
	}
	bodyBytes, err := proto.Marshal(body)
	if err != nil {
		return err
	}
	payload := bytes.NewBuffer(nil)
	if err := binary.Write(payload, binary.LittleEndian, uint32(eMsg)|protoMask); err != nil {
		return err
	}
	if err := binary.Write(payload, binary.LittleEndian, uint32(len(hdrBytes))); err != nil {
		return err
	}
	payload.Write(hdrBytes)
	payload.Write(bodyBytes)

	framed := bytes.NewBuffer(nil)
	if err := binary.Write(framed, binary.LittleEndian, uint32(payload.Len())); err != nil {
		return err
	}
	if err := binary.Write(framed, binary.LittleEndian, tcpMagic); err != nil {
		return err
	}
	framed.Write(payload.Bytes())

	_, err = conn.Write(framed.Bytes())
	return err
}

func writeRawPacket(conn net.Conn, eMsg EMsg, body []byte) error {
	payload := bytes.NewBuffer(nil)
	if err := binary.Write(payload, binary.LittleEndian, uint32(eMsg)); err != nil {
		return err
	}
	payload.Write(body)

	framed := bytes.NewBuffer(nil)
	if err := binary.Write(framed, binary.LittleEndian, uint32(payload.Len())); err != nil {
		return err
	}
	if err := binary.Write(framed, binary.LittleEndian, tcpMagic); err != nil {
		return err
	}
	framed.Write(payload.Bytes())

	_, err := conn.Write(framed.Bytes())
	return err
}

func writeProtoBodyPacket(conn net.Conn, eMsg EMsg, body []byte) error {
	payload := bytes.NewBuffer(nil)
	if err := binary.Write(payload, binary.LittleEndian, uint32(eMsg)|protoMask); err != nil {
		return err
	}
	if err := binary.Write(payload, binary.LittleEndian, uint32(0)); err != nil {
		return err
	}
	payload.Write(body)

	framed := bytes.NewBuffer(nil)
	if err := binary.Write(framed, binary.LittleEndian, uint32(payload.Len())); err != nil {
		return err
	}
	if err := binary.Write(framed, binary.LittleEndian, tcpMagic); err != nil {
		return err
	}
	framed.Write(payload.Bytes())

	_, err := conn.Write(framed.Bytes())
	return err
}

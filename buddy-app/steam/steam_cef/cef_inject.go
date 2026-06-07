package steamcef

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"text/template"
	"time"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const DefaultInjectTab = "SharedJSContext"

type Executor[Args any, R any] interface {
	Execute(ctx context.Context, args Args) (R, error)
	ExecuteInTab(ctx context.Context, tab string, args Args) (R, error)
}

type executor[Args any, R any] struct {
	cfg  *appconfig.Steam
	tmpl *template.Template
}

func NewExecutor[Args any, R any](cfg *appconfig.Steam, tmpl *template.Template) Executor[Args, R] {
	return &executor[Args, R]{cfg: cfg, tmpl: tmpl}
}

func (e *executor[Args, R]) Execute(ctx context.Context, args Args) (R, error) {
	return e.ExecuteInTab(ctx, DefaultInjectTab, args)
}

func (e *executor[Args, R]) ExecuteInTab(ctx context.Context, tab string, args Args) (R, error) {
	var r R
	var buf bytes.Buffer
	if err := e.tmpl.Execute(&buf, args); err != nil {
		return r, fmt.Errorf("%w: %w", ErrCEFTemplateExecFailed, err)
	}
	js := buf.String()
	resStr, err := executeJs(ctx, e.cfg, tab, js)
	if err != nil {
		return r, err
	}
	if resStr == nil {
		return r, nil
	}
	if _, ok := any(r).(string); ok {
		return any(*resStr).(R), nil
	}
	if _, ok := any(r).(*string); ok {
		return any(resStr).(R), nil
	}
	if err := json.Unmarshal([]byte(*resStr), &r); err != nil {
		return r, err
	}
	return r, nil
}

func executeJs(ctx context.Context, cfg *appconfig.Steam, tab string, js string) (*string, error) {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		slog.Debug("executeJs: applied default 5s timeout")
	}

	slog.Debug("executeJs: starting", "tab", tab)

	tabs, err := steam.GetCEFTabs(ctx, cfg)
	if err != nil {
		slog.Debug("executeJs: GetCEFTabs failed", "err", err)
		return nil, err
	}
	slog.Debug("executeJs: fetched tabs", "count", len(tabs), "tabs", tabs)
	webSocketDebugURL := ""
	for _, t := range tabs {
		if strings.EqualFold(t.Title, tab) {
			webSocketDebugURL = t.WebSocketDebuggerURL
			break
		}
	}
	if webSocketDebugURL == "" {
		slog.Debug("executeJs: target tab not found", "tab", tab)
		return nil, ErrCEFTabNotFound
	}

	ws, _, err := websocket.Dial(ctx, webSocketDebugURL, &websocket.DialOptions{HTTPHeader: http.Header{}})
	if err != nil {
		slog.Debug("executeJs: websocket dial failed", "err", err)
		return nil, fmt.Errorf("%w: %w", ErrCEFRemoteDebugUnreachable, err)
	}
	defer ws.CloseNow()
	slog.Debug("executeJs: websocket connected")

	ws.SetReadLimit(10 * 1024 * 1024)

	go func() {
		<-ctx.Done()
		slog.Debug("executeJs: context cancelled, closing websocket")
		_ = ws.CloseNow()
	}()

	request := NewExecutePayload(js)
	slog.Debug("executeJs: sending payload", "id", request.ID)

	if err := wsjson.Write(ctx, ws, request); err != nil {
		slog.Debug("executeJs: failed to send payload", "err", err)
		return nil, fmt.Errorf("%w: %w", ErrFailedToExecuteJs, err)
	}
	slog.Debug("executeJs: payload sent, waiting for response", "id", request.ID)

	for {
		select {
		case <-ctx.Done():
			slog.Debug("executeJs: context timeout/cancellation", "err", ctx.Err())
			return nil, fmt.Errorf("%w: %w", ErrFailedToExecuteJs, ctx.Err())
		default:
		}

		var raw map[string]any
		if err := wsjson.Read(ctx, ws, &raw); err != nil {
			slog.Debug("executeJs: receive error", "err", err)
			return nil, fmt.Errorf("%w: %w", ErrCEFResponseReadFailed, err)
		}
		slog.Debug("executeJs: received message", "msg", raw)

		if protocolErr, ok := raw["error"]; ok {
			slog.Debug("executeJs: protocol error in response", "err", protocolErr)
			return nil, fmt.Errorf("%w: %v", ErrCEFProtocolResponseError, protocolErr)
		}

		id, hasID := raw["id"]
		if !hasID {
			slog.Debug("executeJs: message has no id, skipping")
			continue
		}
		idFloat, ok := id.(float64)
		if !ok {
			slog.Debug("executeJs: message id has unexpected type, skipping", "id", id)
			continue
		}
		if int(idFloat) != request.ID {
			slog.Debug("executeJs: message id mismatch, skipping", "want", request.ID, "got", id)
			continue
		}
		slog.Debug("executeJs: matching response found", "id", request.ID)

		resultMap, ok := raw["result"].(map[string]any)
		if !ok {
			slog.Debug("executeJs: malformed response - no result map")
			return nil, ErrCEFMalformedResponse
		}
		if exceptionDetails, ok := resultMap["exceptionDetails"]; ok {
			slog.Debug("executeJs: javascript exception", "exception", exceptionDetails)
			return nil, fmt.Errorf("%w: %v", ErrCEFJavaScriptException, exceptionDetails)
		}

		runtimeResult, ok := resultMap["result"].(map[string]any)
		if !ok {
			slog.Debug("executeJs: no runtime result")
			return nil, ErrCEFMissingRuntimeResult
		}
		value, ok := runtimeResult["value"]
		if !ok {
			if runtimeResult["type"] == "undefined" {
				slog.Debug("executeJs: void/undefined result")
				return nil, nil
			}
			slog.Debug("executeJs: no value in runtime result")
			return nil, ErrCEFMissingRuntimeValue
		}

		if s, ok := value.(string); ok {
			return &s, nil
		}
		if s, ok := value.(*string); ok {
			return s, nil
		}
		b, err := json.Marshal(value)
		if err != nil {
			slog.Debug("executeJs: failed to marshal value", "err", err)
			return nil, fmt.Errorf("%w: %w", ErrCEFEncodeResultFailed, err)
		}
		slog.Debug("executeJs: returning marshalled result")
		result := string(b)
		return &result, nil
	}
}

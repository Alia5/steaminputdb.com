package steamcef

import "errors"

var ErrCEFTabNotFound = errors.New("CEF tab not found")
var ErrCEFRemoteDebugUnreachable = errors.New("CEF remote debug websocket unreachable")
var ErrFailedToExecuteJs = errors.New("failed to execute JS in CEF")
var ErrCEFResponseReadFailed = errors.New("failed to read CEF websocket response")
var ErrCEFProtocolResponseError = errors.New("CEF protocol response error")
var ErrCEFMalformedResponse = errors.New("malformed CEF response")
var ErrCEFJavaScriptException = errors.New("javascript exception in CEF execution")
var ErrCEFMissingRuntimeResult = errors.New("missing runtime result in CEF response")
var ErrCEFMissingRuntimeValue = errors.New("missing runtime value in CEF response")
var ErrCEFEncodeResultFailed = errors.New("failed to encode CEF execution result")
var ErrCEFTemplateExecFailed = errors.New("failed to execute JS template")

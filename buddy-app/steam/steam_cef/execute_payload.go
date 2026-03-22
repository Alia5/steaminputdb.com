package steamcef

var idCounter = 0

type ExecutePayload struct {
	ID     int                  `json:"id"`
	Method string               `json:"method"`
	Params ExecutePayloadParams `json:"params"`
}

type ExecutePayloadParams struct {
	Expression    string `json:"expression"`
	ReturnByValue bool   `json:"returnByValue"`
	AwaitPromise  bool   `json:"awaitPromise"`
}

func NewExecutePayload(js string) ExecutePayload {
	return ExecutePayload{
		ID:     NextInjectID(),
		Method: "Runtime.evaluate",
		Params: ExecutePayloadParams{
			Expression:    js,
			ReturnByValue: true,
			AwaitPromise:  true,
		},
	}
}

func NextInjectID() int {
	idCounter++
	return idCounter
}

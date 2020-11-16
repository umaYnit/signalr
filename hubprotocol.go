package signalr

import (
	"io"
)

// HubProtocol interface
// ReadMessage() reads a message from buf and returns the message if the buf contained one completely.
// If buf does not contain the whole message, it returns a nil message and complete false
// WriteMessage writes a message to the specified writer
// UnmarshalArgument() unmarshals a raw message depending of the specified value type into value
type HubProtocol interface {
	ParseMessage(reader io.Reader) (interface{}, error)
	WriteMessage(message interface{}, writer io.Writer) error
	UnmarshalArgument(argument interface{}, value interface{}) error
	setDebugLogger(dbg StructuredLogger)
}

type HubAdapter interface {
	// target is the method name, arguments is a protocol specific slice
	// This func branches between protocol specific sub funcs
	// The sub funcs have switch which branches between methods
	Invoke(target string, arguments interface{}, streamIds []string, protocol HubProtocol) (result interface{})
	IntoChan(target string, chanIndex int, inChan interface{}, item []byte, protocol HubProtocol)
	FromChan(target string, outChan interface{})
}

//easyjson:json
type hubMessage struct {
	Type int `json:"type"`
}

// easyjson:json
type invocationMessage struct {
	Type         int           `json:"type"`
	Target       string        `json:"target"`
	InvocationID string        `json:"invocationId,omitempty"`
	Arguments    []interface{} `json:"arguments"`
	StreamIds    []string      `json:"streamIds,omitempty"`
}

//easyjson:json
type completionMessage struct {
	Type         int         `json:"type"`
	InvocationID string      `json:"invocationId"`
	Result       interface{} `json:"result,omitempty"`
	Error        string      `json:"error,omitempty"`
}

//easyjson:json
type streamItemMessage struct {
	Type         int         `json:"type"`
	InvocationID string      `json:"invocationId"`
	Item         interface{} `json:"item"`
}

//easyjson:json
type cancelInvocationMessage struct {
	Type         int    `json:"type"`
	InvocationID string `json:"invocationId"`
}

//easyjson:json
type closeMessage struct {
	Type           int    `json:"type"`
	Error          string `json:"error"`
	AllowReconnect bool   `json:"allowReconnect"`
}

//easyjson:json
type handshakeRequest struct {
	Protocol string `json:"protocol"`
	Version  int    `json:"version"`
}

//easyjson:json
type handshakeResponse struct {
	Error string `json:"error,omitempty"`
}

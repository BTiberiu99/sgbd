package response

import "encoding/json"

//Message ... is used to send messages and data back to the user interface
type Message struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (m *Message) String() string {
	byts, _ := json.Marshal(m)
	return string(byts)
}

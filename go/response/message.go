package response

import "encoding/json"

type Message struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (m *Message) String() string {
	byts, _ := json.Marshal(m)
	return string(byts)
}

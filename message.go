package wirepusher

// Message to send using WirePusher
type Message struct {
	ID        int
	Title     string
	Body      string
	Type      string
	ActionURL string
	ImageURL  string
}

// NewMsg returns a new message without a specific ID. This message cannot be cleared at a later time when no ID is provided.
func NewMsg(msgtype, title, body string) *Message {
	return &Message{Type: msgtype, Title: title, Body: body}
}

// MsgWithID returns a new message with a specific ID. This message can be cleared at a later time by referencing the ID.
func MsgWithID(msgid int, msgtype, title, body string) *Message {
	return &Message{ID: msgid, Type: msgtype, Title: title, Body: body}
}

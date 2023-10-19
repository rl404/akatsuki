package entity

type messageType string

// Available message type.
const (
	TypeParseAnime     messageType = "parse-anime"
	TypeParseUserAnime messageType = "parse-user-anime"
)

// Message is entity for message.
type Message struct {
	Type     messageType `json:"type"`
	ID       int64       `json:"id"`
	Username string      `json:"username"`
	Status   string      `json:"status"`
	Forced   bool        `json:"forced"`
}

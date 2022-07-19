package entity

type messageType string

// Available message type.
const (
	TypeParseAnime messageType = "parse-anime"
)

// Message is entity for message.
type Message struct {
	Type messageType `json:"type"`
	Data []byte      `json:"data"`
}

// ParseAnimeRequest is parse anime request model.
type ParseAnimeRequest struct {
	ID     int64 `json:"id"`
	Forced bool  `json:"forced"`
}

package entity

// GetUserAnimeRequest is get user anime request entity.
type GetUserAnimeRequest struct {
	Username string
	Limit    int
	Offset   int
}

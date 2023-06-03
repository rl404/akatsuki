package entity

// GetUserAnimeRequest is get user anime request entity.
type GetUserAnimeRequest struct {
	Username string
	Status   string
	Limit    int
	Offset   int
}

package entity

// Genre is entity for genre.
type Genre struct {
	ID    int64
	Name  string
	Count int
}

// GetRequest is get genre list request model.
type GetRequest struct {
	Name  string
	Page  int
	Limit int
}

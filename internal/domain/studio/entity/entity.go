package entity

// Studio is entity for studio.
type Studio struct {
	ID    int64
	Name  string
	Count int
}

// GetRequest is get studio list request model.
type GetRequest struct {
	Name  string
	Page  int
	Limit int
}

package entity

// Studio is entity for studio.
type Studio struct {
	ID     int64
	Name   string
	Count  int
	Mean   float64
	Member int
}

// GetRequest is get studio list request model.
type GetRequest struct {
	Name  string
	Sort  Sort
	Page  int
	Limit int
}

// History is entity for studio history.
type History struct {
	Year       int
	Month      int
	Mean       float64
	Rank       int
	Popularity int
	Member     int
	Voter      int
	Count      int
}

// GetHistoriesRequest is get histories request model.
type GetHistoriesRequest struct {
	StudioID  int64
	StartYear int
	EndYear   int
	Group     HistoryGroup
}

package entity

// Genre is entity for genre.
type Genre struct {
	ID     int64
	Name   string
	Count  int
	Mean   float64
	Member int
}

// GetRequest is get genre list request model.
type GetRequest struct {
	Name  string
	Sort  Sort
	Page  int
	Limit int
}

// History is entity for genre history.
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
	GenreID   int64
	StartYear int
	EndYear   int
	Group     HistoryGroup
}

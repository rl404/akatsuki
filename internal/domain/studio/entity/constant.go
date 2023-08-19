package entity

// HistoryGroup is studio history group.
type HistoryGroup string

// Available studio history group.
const (
	Yearly  HistoryGroup = "YEARLY"
	Monthly HistoryGroup = "MONTHLY"
)

// Sort is studio sorting.
type Sort string

// Available studio sorting.
const (
	SortName   Sort = "NAME"
	SortCount  Sort = "COUNT"
	SortMean   Sort = "MEAN"
	SortMember Sort = "MEMBER"
)

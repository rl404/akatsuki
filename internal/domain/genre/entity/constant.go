package entity

// HistoryGroup is genre history group.
type HistoryGroup string

// Available genre history group.
const (
	Yearly  HistoryGroup = "YEARLY"
	Monthly HistoryGroup = "MONTHLY"
)

// Sort is genre sorting.
type Sort string

// Available genre sorting.
const (
	SortName   Sort = "NAME"
	SortCount  Sort = "COUNT"
	SortMean   Sort = "MEAN"
	SortMember Sort = "MEMBER"
)

package sql

import "time"

// EmptyID is empty_id database model.
type EmptyID struct {
	AnimeID   int64 `gorm:"primaryKey"`
	CreatedAt time.Time
}

package sql

import (
	"time"

	"github.com/rl404/akatsuki/internal/domain/studio/entity"
	"gorm.io/gorm"
)

// Studio is studio database model.
type Studio struct {
	ID        int64 `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (sql *SQL) fromEntities(data []entity.Studio) []Studio {
	s := make([]Studio, len(data))
	for i, ss := range data {
		s[i] = Studio{
			ID:   ss.ID,
			Name: ss.Name,
		}
	}
	return s
}

func (s *Studio) toEntity() *entity.Studio {
	return &entity.Studio{
		ID:   s.ID,
		Name: s.Name,
	}
}

func (sql *SQL) toEntities(data []Studio) []*entity.Studio {
	s := make([]*entity.Studio, len(data))
	for i, ss := range data {
		s[i] = ss.toEntity()
	}
	return s
}

type studioHistory struct {
	Year       int
	Month      int
	Mean       float64
	Rank       int
	Popularity int
	Member     int
	Voter      int
	Count      int
}

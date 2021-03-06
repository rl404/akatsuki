package sql

import (
	"time"

	"github.com/rl404/akatsuki/internal/domain/genre/entity"
	"gorm.io/gorm"
)

// Genre is genre database model.
type Genre struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (sql *SQL) fromEntities(data []entity.Genre) []Genre {
	g := make([]Genre, len(data))
	for i, gg := range data {
		g[i] = Genre{
			ID:   gg.ID,
			Name: gg.Name,
		}
	}
	return g
}

func (g *Genre) toEntity() *entity.Genre {
	return &entity.Genre{
		ID:   g.ID,
		Name: g.Name,
	}
}

func (sql *SQL) toEntities(data []Genre) []*entity.Genre {
	g := make([]*entity.Genre, len(data))
	for i, gg := range data {
		g[i] = gg.toEntity()
	}
	return g
}

package sql

import (
	"fmt"
	"strings"
	"time"

	"github.com/rl404/akatsuki/internal/domain/genre/entity"
	"gorm.io/gorm"
)

// Genre is genre database model.
type Genre struct {
	ID        int64 `gorm:"primaryKey"`
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

type genreHistory struct {
	Year       int
	Month      int
	Mean       float64
	Rank       int
	Popularity int
	Member     int
	Voter      int
	Count      int
}

func (sql *SQL) convertSort(sort entity.Sort) string {
	if sort == "" {
		sort = entity.SortName
	}

	suffix := "asc"
	if sort[0] == '-' {
		sort, suffix = sort[1:], "desc"
	}

	switch sort {
	case entity.SortName:
		return fmt.Sprintf("lower(g.name) %s", suffix)
	case entity.SortMean:
		return fmt.Sprintf("avg(nullif(a.mean, 0)) = 0 nulls last, %s %s", strings.ToLower(string(sort)), suffix)
	case entity.SortMember:
		return fmt.Sprintf("sum(a.member) = 0 nulls last, %s %s", strings.ToLower(string(sort)), suffix)
	default:
		return fmt.Sprintf("%s %s", strings.ToLower(string(sort)), suffix)
	}
}

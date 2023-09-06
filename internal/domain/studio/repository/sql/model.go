package sql

import (
	"fmt"
	"strings"
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

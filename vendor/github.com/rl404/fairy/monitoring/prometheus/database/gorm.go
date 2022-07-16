// Package database is a prometheus wrapper for database.
package database

import (
	"context"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

const gormReqName = "database_requests_total"
const gormLatencyName = "database_request_duration_seconds"

type gormTimeNowCtxKey struct{}

var gormTimeNowCtx = gormTimeNowCtxKey(struct{}{})

type gormPrometheus struct {
	req *prometheus.CounterVec
	lat *prometheus.HistogramVec
}

var gp gormPrometheus

func init() {
	gp = gormPrometheus{
		req: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: gormReqName,
				Help: "How many database queries processed, partitioned by database, operation and table.",
			},
			[]string{"database", "operation", "table"},
		),
		lat: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: gormLatencyName,
				Help: "How long it took to process the query, partitioned by database, operation and table.",
			},
			[]string{"database", "operation", "table"},
		),
	}

	prometheus.MustRegister(gp.req)
	prometheus.MustRegister(gp.lat)
}

// RegisterGORM to register prometheus plugin for gorm.
func RegisterGORM(dbName string, db *gorm.DB) {
	db.Callback().Query().Before("gorm:query").Register("prometheus:before_query", gp.start)
	db.Callback().Query().After("gorm:query").Register("prometheus:after_query", gp.end(dbName))

	db.Callback().Create().Before("gorm:create").Register("prometheus:before_create", gp.start)
	db.Callback().Create().After("gorm:create").Register("prometheus:after_create", gp.end(dbName))

	db.Callback().Delete().Before("gorm:delete").Register("prometheus:before_delete", gp.start)
	db.Callback().Delete().After("gorm:delete").Register("prometheus:after_delete", gp.end(dbName))

	db.Callback().Update().Before("gorm:update").Register("prometheus:before_update", gp.start)
	db.Callback().Update().After("gorm:update").Register("prometheus:after_update", gp.end(dbName))

	db.Callback().Row().Before("gorm:row").Register("prometheus:before_row", gp.start)
	db.Callback().Row().After("gorm:row").Register("prometheus:after_row", gp.end(dbName))
}

func (gp *gormPrometheus) start(db *gorm.DB) {
	db.Statement.Context = context.WithValue(db.Statement.Context, gormTimeNowCtx, time.Now())
}

func (gp *gormPrometheus) end(name string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		start, ok := db.Statement.Context.Value(gormTimeNowCtx).(time.Time)
		if !ok {
			return
		}

		operation := "OTHER"
		ops := strings.Split(strings.TrimSpace(db.Statement.SQL.String()), " ")
		if len(ops) > 0 {
			operation = strings.ToUpper(ops[0])
		}

		table := db.Statement.Table

		gp.req.WithLabelValues(name, operation, table).Inc()
		gp.lat.WithLabelValues(name, operation, table).Observe(float64(time.Since(start).Seconds()))
	}
}

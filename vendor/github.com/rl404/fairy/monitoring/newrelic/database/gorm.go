package database

import (
	"context"
	"strings"

	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm"
)

type gormTimeNowCtxKey struct{}

var gormTimeNowCtx = gormTimeNowCtxKey(struct{}{})

// RegisterGORM to register newrelic plugin for gorm.
func RegisterGORM(dialect, dbName string, db *gorm.DB) {
	db.Callback().Query().Before("gorm:query").Register("newrelic:before_query", start)
	db.Callback().Query().After("gorm:query").Register("newrelic:after_query", end(dialect, dbName))

	db.Callback().Create().Before("gorm:create").Register("newrelic:before_create", start)
	db.Callback().Create().After("gorm:create").Register("newrelic:after_create", end(dialect, dbName))

	db.Callback().Delete().Before("gorm:delete").Register("newrelic:before_delete", start)
	db.Callback().Delete().After("gorm:delete").Register("newrelic:after_delete", end(dialect, dbName))

	db.Callback().Update().Before("gorm:update").Register("newrelic:before_update", start)
	db.Callback().Update().After("gorm:update").Register("newrelic:after_update", end(dialect, dbName))

	db.Callback().Row().Before("gorm:row").Register("newrelic:before_row", start)
	db.Callback().Row().After("gorm:row").Register("newrelic:after_row", end(dialect, dbName))
}

func start(db *gorm.DB) {
	s := newrelic.FromContext(db.Statement.Context).StartSegmentNow()
	db.Statement.Context = context.WithValue(db.Statement.Context, gormTimeNowCtx, s)
}

func end(dialect, name string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		seg, ok := db.Statement.Context.Value(gormTimeNowCtx).(newrelic.SegmentStartTime)
		if !ok {
			return
		}

		operation := "OTHER"
		ops := strings.Split(strings.TrimSpace(db.Statement.SQL.String()), " ")
		if len(ops) > 0 {
			operation = strings.ToUpper(ops[0])
		}

		s := newrelic.DatastoreSegment{
			StartTime:          seg,
			Product:            newrelic.DatastoreProduct(dialect),
			Collection:         db.Statement.Table,
			Operation:          operation,
			ParameterizedQuery: db.Statement.SQL.String(),
			DatabaseName:       name,
		}

		s.End()
	}
}

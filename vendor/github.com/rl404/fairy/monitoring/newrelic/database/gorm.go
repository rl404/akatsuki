package database

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm"
)

type gormPlugin struct {
	host   string
	port   string
	dbName string
}

// NewGORM to create new newrelic plugin for gorm.
func NewGORM(address, dbName string) *gormPlugin {
	host, port := splitHostPort(address)
	return &gormPlugin{
		host:   host,
		port:   port,
		dbName: dbName,
	}
}

// Name is plugin name.
func (gp *gormPlugin) Name() string {
	return "newrelic"
}

// Initialize to init newrelic plugin.
func (gp *gormPlugin) Initialize(db *gorm.DB) error {
	db.Callback().Query().Before("gorm:query").Register("newrelic:before_query", gp.start("SELECT"))
	db.Callback().Query().After("gorm:query").Register("newrelic:after_query", gp.end("SELECT"))

	db.Callback().Create().Before("gorm:create").Register("newrelic:before_create", gp.start("INSERT"))
	db.Callback().Create().After("gorm:create").Register("newrelic:after_create", gp.end("INSERT"))
	db.Callback().Create().Before("gorm:begin_transaction").Register("newrelic:begin_transaction_create", gp.start("TRANSACTION"))
	db.Callback().Create().After("gorm:commit_or_rollback_transaction").Register("newrelic:commit_or_rollback_transaction_create", gp.end("TRANSACTION"))

	db.Callback().Delete().Before("gorm:delete").Register("newrelic:before_delete", gp.start("DELETE"))
	db.Callback().Delete().After("gorm:delete").Register("newrelic:after_delete", gp.end("DELETE"))
	db.Callback().Delete().Before("gorm:begin_transaction").Register("newrelic:begin_transaction_delete", gp.start("TRANSACTION"))
	db.Callback().Delete().After("gorm:commit_or_rollback_transaction").Register("newrelic:commit_or_rollback_transaction_delete", gp.end("TRANSACTION"))

	db.Callback().Update().Before("gorm:update").Register("newrelic:before_update", gp.start("UPDATE"))
	db.Callback().Update().After("gorm:update").Register("newrelic:after_update", gp.end("UPDATE"))
	db.Callback().Update().Before("gorm:begin_transaction").Register("newrelic:begin_transaction_update", gp.start("TRANSACTION"))
	db.Callback().Update().After("gorm:commit_or_rollback_transaction").Register("newrelic:commit_or_rollback_transaction_update", gp.end("TRANSACTION"))

	db.Callback().Row().Before("gorm:row").Register("newrelic:before_row", gp.start("ROW"))
	db.Callback().Row().After("gorm:row").Register("newrelic:after_row", gp.end("ROW"))

	db.Callback().Raw().Before("gorm:raw").Register("newrelic:before_raw", gp.start("RAW"))
	db.Callback().Raw().After("gorm:raw").Register("newrelic:after_raw", gp.end("RAW"))

	return nil
}

func (gp *gormPlugin) key(operation string) string {
	return fmt.Sprintf("newrelic:startKey:%s", operation)
}

func (gp *gormPlugin) start(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if tx := newrelic.FromContext(db.Statement.Context); tx != nil {
			db.Set(gp.key(operation), tx.StartSegmentNow())
		}
	}
}

func (gp *gormPlugin) end(operation string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		startSegment, ok := db.Get(gp.key(operation))
		if !ok {
			return
		}

		s := newrelic.DatastoreSegment{
			StartTime:          startSegment.(newrelic.SegmentStartTime),
			Product:            newrelic.DatastoreProduct(db.Name()),
			Collection:         db.Statement.Table,
			Operation:          operation,
			ParameterizedQuery: db.Statement.SQL.String(),
			QueryParameters:    gp.parseVars(db.Statement.Vars),
			DatabaseName:       gp.dbName,
			Host:               gp.host,
			PortPathOrID:       gp.port,
		}

		s.End()
	}
}

func (gp *gormPlugin) parseVars(vars []interface{}) map[string]interface{} {
	queryParameters := make(map[string]interface{})
	for i, v := range vars {
		queryParameters[strconv.Itoa(i+1)] = fmt.Sprintf("%v", v)
	}
	return queryParameters
}

func splitHostPort(address string) (host string, port string) {
	host = address

	colon := strings.LastIndexByte(host, ':')
	if colon != -1 && validOptionalPort(host[colon:]) {
		host, port = host[:colon], host[colon+1:]
	}

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	return
}

func validOptionalPort(port string) bool {
	if port == "" {
		return true
	}
	if port[0] != ':' {
		return false
	}
	for _, b := range port[1:] {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

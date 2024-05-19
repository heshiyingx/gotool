package gormdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"
)

type DBType string

const (
	DBTYPE_MySQL      DBType = "MYSQL"
	DBTYPE_SQLite     DBType = "SQLite"
	DBTYPE_PostgreSQL DBType = "PostgreSQL"
	DBTYPE_SQL_Server DBType = "SQL_Server"
	DBTYPE_SQL_TiDB   DBType = "TiDB"
	DBTYPE_ClickHouse DBType = "ClickHouse"
)

type (
	QueryCtxFn func(ctx context.Context, result any, db *gorm.DB) error
	ExecCtxFn  func(ctx context.Context, db *gorm.DB) (int64, error)
	CacheFn    func(result string) error
	Config     struct {
		DSN               string
		DBType            DBType
		GormConfig        gorm.Config
		Rdb               redis.UniversalClient
		NotFoundExpireSec int
		CacheExpireSec    int
		RandSec           int
	}

)

func getDialector(c Config) gorm.Dialector {
	var dialetor gorm.Dialector
	switch c.DBType {
	case DBTYPE_SQL_TiDB:
		fallthrough
	case DBTYPE_MySQL:
		dialetor = mysql.Open(c.DSN)
	case DBTYPE_SQLite:
		dialetor = sqlite.Open(c.DSN)
	case DBTYPE_SQL_Server:
		dialetor = sqlserver.Open(c.DSN)
	case DBTYPE_PostgreSQL:
		dialetor = postgres.Open(c.DSN)
	case DBTYPE_ClickHouse:
		dialetor = clickhouse.Open(c.DSN)
	default:
		log.Fatalf("DBType %s not supported", c.DBType)
	}
	return dialetor
}

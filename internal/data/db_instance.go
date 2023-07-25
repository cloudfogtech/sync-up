package data

import (
	"errors"

	"github.com/cloudfogtech/sync-up/internal/common"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func GetDatabase(dbType, dsn string) (gorm.Dialector, error) {
	switch dbType {
	case common.DBTypeSqlite:
		return sqlite.Open(dsn), nil
	case common.DBTypeMySQL:
		return mysql.Open(dsn), nil
	case common.DBTypePostgres:
		return postgres.Open(dsn), nil
	case common.DBTypeSQLServer:
		return sqlserver.Open(dsn), nil
	}
	return nil, errors.New("unsupported database type")
}

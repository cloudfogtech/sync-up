package common

const (
	DefaultDebug            = "false"
	DefaultPort             = "13003"
	DefaultDBType           = "sqlite"
	DefaultUsername         = "admin"
	DefaultPassword         = "admin"
	DefaultReportTypes      = "local"
	DefaultReportLocalPath  = "/data/report.json"
	DefaultFESalt           = "bed10de517c68606cb44453f38832d38f332fb6b" // sha-1 of 'syncup_fe'
	DefaultBESalt           = "d02f52421fa700fe5b82bdcf31c2f9cd1774cdf8" // sha-1 of 'syncup_be'
	DefaultSecretKey        = "b8555b47784f612c36c7ca5561c9a3cfd806800a" // sha-1 of 'syncup'
	DefaultDBDsnSqlite      = "sqlite.db"
	DefaultRedisRDBDumpPath = "/data/dump.rdb"
)

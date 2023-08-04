package env

const (
	Split  = "_"
	Prefix = "SYNCUP" + Split
)

const (
	PostfixPath         = Split + "PATH"
	PostfixFilePath     = Split + "FILE" + PostfixPath
	PostfixDirPath      = Split + "DIR" + PostfixPath
	PostfixType         = Split + "TYPE"
	PostfixContainer    = Split + "CONTAINER"
	PostfixLocalDirPath = Split + "LOCAL" + PostfixDirPath
	PostfixCron         = Split + "CRON"
)

const (
	PostfixRClone               = Split + "RC"
	PostfixRCloneRemote         = PostfixRClone + Split + "REMOTE"
	PostfixRCloneDirPath        = PostfixRClone + PostfixDirPath
	PostfixRCloneBandwidthLimit = PostfixRClone + "Split+BANDWIDTH_LIMIT"
	PostfixRCloneRemoteName     = PostfixRCloneRemote + Split + "NAME"
	PostfixRCloneRemotePath     = PostfixRCloneRemote + PostfixPath
)

const (
	VarUsername     = "USERNAME"
	PostfixUsername = Split + VarUsername
	VarPassword     = "PASSWORD"
	PostfixPassword = Split + VarPassword
)

const (
	Placeholder    = "%s"
	ServiceIdRegEx = "([0-9a-zA-Z_-]+)"
)

// Common

const (
	DebugTpl                = Prefix + "DEBUG"
	PortTpl                 = Prefix + "PORT"
	DBTypeTpl               = Prefix + "DB"
	DBDsnTpl                = DBTypeTpl + Split + "DSN"
	UsernameTpl             = Prefix + "USERNAME"
	PasswordTpl             = Prefix + "PASSWORD"
	SecretKeyTpl            = Prefix + "SECRET_KEY"
	ReportTpl               = Prefix + "REPORT" + Split
	ReportTypesTpl          = ReportTpl + "TYPES"
	ReportLocalPath         = ReportTpl + "LOCAL" + Split + "PATH"
	ReportEmailTpl          = ReportTpl + "EMAIL" + Split
	ReportEmailSmtpHost     = ReportEmailTpl + "HOST"
	ReportEmailSmtpPort     = ReportEmailTpl + "PORT"
	ReportEmailSmtpUsername = ReportEmailTpl + "USERNAME"
	ReportEmailSmtpPassword = ReportEmailTpl + "PASSWORD"
	ReportEmailSmtpFrom     = ReportEmailTpl + "FROM"
	ReportEmailSmtpTo       = ReportEmailTpl + "To"
	TypeTpl                 = Prefix + Placeholder + PostfixType
	ContainerTpl            = Prefix + Placeholder + PostfixContainer
	LocalDirPathTpl         = Prefix + Placeholder + PostfixLocalDirPath
	CronTpl                 = Prefix + Placeholder + PostfixCron
)

// Common addition

const (
	ServiceUsernameTpl = Prefix + Placeholder + PostfixUsername
	ServicePasswordTpl = Prefix + Placeholder + PostfixPassword
)

// RClone

const (
	RCTpl               = Prefix + Placeholder + PostfixRClone
	RCDirPathTpl        = Prefix + Placeholder + PostfixRCloneDirPath
	RCBandwidthLimitTpl = Prefix + Placeholder + PostfixRCloneBandwidthLimit
	RCRemoteNameTpl     = Prefix + Placeholder + PostfixRCloneRemoteName
	RCRemotePathTpl     = Prefix + Placeholder + PostfixRCloneRemotePath
)

// Redis

const RedisRDBFilePathTpl = Prefix + Placeholder + "_RDB_FILE_PATH"

// SQLite

const SQLiteFilePathTpl = Prefix + Placeholder + "_SQLITE_FILE_PATH"

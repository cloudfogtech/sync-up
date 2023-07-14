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
	PostfixCorn         = Split + "CORN"
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
	ServicePlaceholder = "%s"
	ServiceIdRegEx     = "([0-9a-zA-Z_-]+)"
)

// Common

const (
	PortTpl         = Prefix + "PORT"
	DBTypeTpl       = Prefix + "DB"
	DBDsnTpl        = DBTypeTpl + Split + "DSN"
	UsernameTpl     = Prefix + "USERNAME"
	PasswordTpl     = Prefix + "PASSWORD"
	SecretKeyTpl    = Prefix + "SECRET_KEY"
	TypeTpl         = Prefix + ServicePlaceholder + PostfixType
	ContainerTpl    = Prefix + ServicePlaceholder + PostfixContainer
	LocalDirPathTpl = Prefix + ServicePlaceholder + PostfixLocalDirPath
	CronTpl         = Prefix + ServicePlaceholder + PostfixCorn
)

// Common addition

const (
	ServiceUsernameTpl = Prefix + ServicePlaceholder + PostfixUsername
	ServicePasswordTpl = Prefix + ServicePlaceholder + PostfixPassword
)

// RClone

const (
	RCTpl               = Prefix + ServicePlaceholder + PostfixRClone
	RCDirPathTpl        = Prefix + ServicePlaceholder + PostfixRCloneDirPath
	RCBandwidthLimitTpl = Prefix + ServicePlaceholder + PostfixRCloneBandwidthLimit
	RCRemoteNameTpl     = Prefix + ServicePlaceholder + PostfixRCloneRemoteName
	RCRemotePathTpl     = Prefix + ServicePlaceholder + PostfixRCloneRemotePath
)

// Redis

const RedisRDBFilePathTpl = Prefix + ServicePlaceholder + "_RDB_FILE_PATH"

package sync

type SqliteContainerBackupper struct {
	BackupperOptions
	DbFilePath string
}

type SqliteContainerBackupperOptions struct {
	ContainerName string
	DbFilePath    string
}

func NewSqliteContainer(options SqliteContainerBackupperOptions) Backupper {
	return &SqliteContainerBackupper{
		BackupperOptions{
			ContainerName: options.ContainerName,
		},
		options.DbFilePath,
	}
}

func (p *SqliteContainerBackupper) Type() string {
	return "SqliteContainerBackupper"
}

func (p *SqliteContainerBackupper) ContainerName() string {
	return p.BackupperOptions.ContainerName
}

func (p *SqliteContainerBackupper) BackupCommand() []string {
	return []string{}
}

func (p *SqliteContainerBackupper) CheckResultError(_ string) bool {
	return false
}

func (p *SqliteContainerBackupper) ContainerDumpFilePath() string {
	return p.DbFilePath
}

func (p *SqliteContainerBackupper) ContainerDumpFileAutoRemove() bool {
	return false
}

func (p *SqliteContainerBackupper) CheckVersionCommands() []string {
	return []string{
		"file",
		p.DbFilePath,
	}
}

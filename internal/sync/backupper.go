package sync

type Backupper interface {
	Type() string
	ContainerName() string
	BackupCommand() []string
	ContainerDumpFilePath() string
	ContainerDumpFileAutoRemove() bool
	CheckResultError(string) bool
	CheckVersionCommands() []string
}

type BackupperOptions struct {
	ContainerName string
}

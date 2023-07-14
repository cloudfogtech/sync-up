package sync

const PostgresDumpFilePath = "/root/postgres.dump"

type PostgresContainerBackupper struct {
	BackupperOptions
	Username string
}

type PostgresContainerBackupperOptions struct {
	ContainerName string
	Username      string
}

func NewPostgresContainer(options PostgresContainerBackupperOptions) Backupper {
	return &PostgresContainerBackupper{
		BackupperOptions{
			ContainerName: options.ContainerName,
		},
		options.Username,
	}
}

func (p *PostgresContainerBackupper) Type() string {
	return "PostgresContainerBackupper"
}

func (p *PostgresContainerBackupper) ContainerName() string {
	return p.BackupperOptions.ContainerName
}

func (p *PostgresContainerBackupper) BackupCommand() []string {
	commands := []string{
		"pg_dumpall",
	}
	if p.Username != "" {
		commands = append(commands, "-U", p.Username)
	}
	commands = append(commands, "-f", PostgresDumpFilePath)
	return commands
}
func (p *PostgresContainerBackupper) CheckResultError(_ string) bool {
	return false
}

func (p *PostgresContainerBackupper) ContainerDumpFilePath() string {
	return PostgresDumpFilePath
}

func (p *PostgresContainerBackupper) ContainerDumpFileAutoRemove() bool {
	return true
}

func (p *PostgresContainerBackupper) CheckVersionCommands() []string {
	return []string{
		"pg_dumpall",
		"-V",
	}
}

package common

import (
	"fmt"
	"runtime"
)

type StartCommand struct {
}
type Command struct {
	StartCommand *StartCommand `arg:"subcommand:start" help:"start application"`
}

func (*Command) Version() string {
	return fmt.Sprintf("%s %s for %s", Name, Version, runtime.GOOS)
}

func (*Command) Epilogue() string {
	return "For more information visit https://github.com/catfishlty/syncup"
}

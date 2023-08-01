package common

import (
	"fmt"
	"runtime"
)

type Command struct{}

func (*Command) Version() string {
	return fmt.Sprintf("%s %s for %s", Name, Version, runtime.GOOS)
}

func (*Command) Epilogue() string {
	return "For more information visit https://github.com/cloudfogtech/syncup"
}

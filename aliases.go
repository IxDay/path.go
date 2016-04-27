package path

import (
	"os"
	"syscall"
)

var (
	ENOENT     = syscall.ENOENT
	IsNotExist = os.IsNotExist
)

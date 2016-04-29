package path

import (
	"errors"
	"fmt"
	"os"
)

var NotAFileError = errors.New("not a file")

type WarningError []Path

func (self *WarningError) Error() string {
	return fmt.Sprintf("%d directories skipped", len(*self))
}

func IsAFileError(err error) bool {
	switch pe := err.(type) {
	case nil:
		return false
	case *os.PathError:
		err = pe.Err
	case *os.LinkError:
		err = pe.Err
	case *os.SyscallError:
		err = pe.Err
	}
	return err == NotAFileError
}

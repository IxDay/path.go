package path

import (
	"errors"
	"fmt"
)

var NotAFileError = errors.New("not a file")

type WarningError []Path

func (self *WarningError) Error() string {
	return fmt.Sprintf("%d directories skipped", len(*self))
}

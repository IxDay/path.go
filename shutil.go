package path

import (
	"os"
)

func (self Path) Remove() error {
	b, err := self.isDir()
	if err != nil {
		return err
	}
	if b {
		return &os.PathError{"remove", string(self), NotAFileError}
	}
	return os.Remove(string(self))
}

func (self Path) RemoveP() error {
	err, ok := self.Remove().(*os.PathError)

	if err == nil || ok && err.Err == ENOENT {
		return nil
	}
	return err
}

func (self Path) RemoveTree() error {
	_, err := self.Stat()
	if err != nil {
		if e, ok := err.(*os.PathError); ok && e.Err == ENOENT {
			e.Op = "remove"
			return e
		}
	}
	return os.RemoveAll(string(self))
}

func (self Path) RemoveTreeP() error {
	return os.RemoveAll(string(self))
}

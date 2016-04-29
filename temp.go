package path

import (
	"io/ioutil"
	"os"
)

type File struct {
	*os.File
	Path
}

// convenient alias from io/ioutil.TempFile
func TempFile(dir, prefix string) (*File, error) {
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return nil, err
	}
	return &File{f, Path(f.Name())}, nil
}

func TempFile_() (*File, error) {
	return TempFile("", "")
}

func (self Path) TempFile(prefix string) (*File, error) {
	return TempFile(string(self), prefix)
}

func (self Path) TempFile_() (*File, error) {
	return self.TempFile("")
}

func (self Path) TmpDir_(cb func(Path)) error {
	return TmpDir(string(self), "", cb)
}

func (self Path) TmpDir(prefix string, cb func(Path)) error {
	return TmpDir(string(self), prefix, cb)
}

func (self Path) TempDir_() (Path, error) {
	return self.TempDir("")
}

func (self Path) TempDir(prefix string) (Path, error) {
	return TempDir(string(self), prefix)
}

func TmpDir(dir, prefix string, cb func(Path)) (err error) {
	var path Path

	if path, err = TempDir(dir, prefix); err != nil {
		return err
	}

	if err = path.Cd(); err != nil {
		return err
	}
	defer path.RemoveTreeP()
	cb(path)
	return nil
}

func TmpDir_(cb func(Path)) error {
	return TmpDir("", "", cb)
}

func TempDir_() (Path, error) {
	return TempDir("", "")
}

func TempDir(dir, prefix string) (Path, error) {
	path, err := ioutil.TempDir(dir, prefix)
	return Path(path), err
}

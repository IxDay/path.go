package path

import (
	"io/ioutil"
	"os"
)

type File struct {
	*os.File
	Path Path
}

// convenient alias from io/ioutil.TempFile
func TempFile(dir, prefix string) (*File, error) {
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return nil, err
	}
	return &File{f, Path(f.Name())}, nil
}

func (self Path) TempFile() (*File, error) {
	return self.TempFileP("")
}

func (self Path) TempFileP(prefix string) (*File, error) {
	return TempFile(string(self), prefix)
}

func TempDir(cb func(Path)) error {
	return TempDirNamed("", "", cb)
}

func TempDirNamed(dir, prefix string, cb func(Path)) error {
	name, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return err
	}
	p := Path(name)
	defer p.RemoveTreeP()
	cb(p)
	return nil
}

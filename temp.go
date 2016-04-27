package path

import (
	"io/ioutil"
)

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

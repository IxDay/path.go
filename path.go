package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

var POSIX = runtime.GOOS != "windows"

type Path string

func createPointer(path string) *Path {
	p_path := new(Path)
	*p_path = Path(path)
	return p_path
}

func (self *Path) Abs() (*Path, error) {
	path, err := filepath.Abs(string(*self))
	if err != nil {
		return nil, err
	}
	return createPointer(path), nil
}

func (self *Path) NormCase() *Path {
	if !POSIX {
		return createPointer(strings.ToLower(string(*self)))
	}
	return self
}

func (self *Path) NormPath() *Path {
	return createPointer(filepath.Clean(string(*self)))
}

func (self *Path) RealPath() (*Path, error) {
	path, err := filepath.EvalSymlinks(string(*self))
	if err != nil {
		return nil, err
	}
	return createPointer(path), nil
}

func expandPosixUser(path string) (*Path, error) {
	i := strings.IndexRune(path, os.PathSeparator)
	switch {
	case i == -1:
		i = len(path)
	case i == 1:
		home := os.Getenv("HOME")
		if home == "" {
			u, err := user.Current()
			if err != nil {
				return nil, err
			}
			home = u.HomeDir
		}
	default:
	}
	return nil, nil
}

func expandNtUser(path string) (*Path, error) {
	return nil, nil
}

func (self *Path) ExpandUser() (*Path, error) {
	path := string(*self)
	if path[0] != '~' {
		return self, nil
	}
	if POSIX {
		return expandPosixUser(path)
	} else {
		return expandNtUser(path)
	}
}

func main() {
	p := Path(".")
	v, _ := p.Abs()
	fmt.Printf("%v\n", *v)

	p = Path("./Toto")
	v = p.NormCase()
	fmt.Printf("%v\n", *v)

	p = Path("..//streams/.")
	v = p.NormPath()
	fmt.Printf("%v\n", *v)

	p = Path("./toto")
	v, _ = p.RealPath()
	fmt.Printf("%v\n", *v)
	// RealPath
	// ExpandUser
}

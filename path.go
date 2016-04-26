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

func (self Path) Abs() (Path, error) {
	path, err := filepath.Abs(string(self))
	return Path(path), err
}

func (self Path) NormCase() Path {
	if !POSIX {
		return Path(strings.ToLower(string(self)))
	}
	return self
}

func (self Path) NormPath() Path {
	return Path(filepath.Clean(string(self)))
}

func (self Path) RealPath() (Path, error) {
	path, err := filepath.EvalSymlinks(string(self))
	return Path(path), err
}

func expandPosixUser(path string) string {
	var home string
	i := strings.IndexRune(path, os.PathSeparator)
	switch {
	case i == -1:
		i = len(path)
		fallthrough
	case i == 1:
		home = os.Getenv("HOME")
		if home == "" {
			u, err := user.Current()
			if err != nil {
				return path
			}
			home = u.HomeDir
		}
	default:
		u, err := user.Lookup(path[1:i])
		if err != nil {
			return path
		}
		home = u.HomeDir
	}
	return home + path[i:]
}

func expandNtUser(path string) string {
	return path
}

func (self Path) ExpandUser() Path {
	path := string(self)
	if path[0] != '~' {
	} else if POSIX {
		path = expandPosixUser(path)
	} else {
		path = expandNtUser(path)
	}

	return Path(path)
}

func (self Path) ExpandVars() Path {
	return self
}

func (self Path) DirName() Path {
	return Path(filepath.Dir(string(self)))
}

func (self Path) BaseName() Path {
	return Path(filepath.Base(string(self)))
}

func (self Path) Expand() Path {
	return self.ExpandVars().ExpandUser().NormPath()
}

func (self Path) SplitExt() (Path, string) {
	filePath := string(self)
	ext := filepath.Ext(filePath)
	length := len(filePath) - len(ext)
	return Path(filePath[:length]), ext
}

func (self Path) NameBase() string {
	base, _ := self.BaseName().SplitExt()
	return string(base)
}

func (self Path) Ext() string {
	return filepath.Ext(string(self))
}

func (self Path) Drive() Path {
	return filepath.VolumeName(string(self))
}

func (self Path) Parent() Path {
	return self.DirName()
}

func main() {
	v, _ := Path(".").Abs()
	fmt.Printf("%s\n", v)

	fmt.Printf("%s\n", Path("./Toto").NormCase())
	fmt.Printf("%s\n", Path("..//streams/.").NormPath())

	v, _ = Path("./toto").RealPath()
	fmt.Printf("%s\n", v)

	fmt.Printf("%s\n", Path("~/toto").ExpandUser())

	// ExpandVars
	fmt.Printf("%s\n", Path("/toto/titi.go").DirName())
	fmt.Printf("%s\n", Path("/toto/titi.go").BaseName())
	fmt.Printf("%s\n", Path("~/toto/titi.go").Expand())
	v, ext := Path("~/toto/titi.go").SplitExt()
	fmt.Printf("base: %s, ext: %s\n", v, ext)

	fmt.Printf("%s\n", Path("~/toto/titi.go").NameBase())
}

// +build !windows

package path

import (
	"os"
	"os/user"
	"strings"
)

func (self Path) NormCase() Path {
	return self
}

func expandUser(path string) string {
	var home string

	if path[0] != '~' {
		return path
	}
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

func (self Path) ExpandUser() Path {
	return Path(expandUser(string(self)))
}

func (self Path) ExpandVars() Path {
	return self
}

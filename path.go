package path

import (
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

var (
	IsNotExist    = os.IsNotExist
	NotAFileError = errors.New("not a file")
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

func (self Path) NameBase() string {
	base, _ := self.BaseName().SplitExt()
	return string(base)
}

func (self Path) Ext() string {
	return filepath.Ext(string(self))
}

func (self Path) Drive() Path {
	return Path(filepath.VolumeName(string(self)))
}

func (self Path) Parent() Path {
	return self.DirName()
}

func (self Path) Name() Path {
	return self.BaseName()
}

func (self Path) SplitPath() (Path, string) {
	dir, file := filepath.Split(string(self))
	return Path(filepath.Dir(dir)), file
}

func (self Path) SplitDrive() (Path, string) {
	drive := self.Drive()
	return drive, string(self)[len(string(drive)):]
}

func (self Path) SplitExt() (Path, string) {
	filePath := string(self)
	ext := filepath.Ext(filePath)
	length := len(filePath) - len(ext)
	return Path(filePath[:length]), ext
}

func (self Path) StripExt() Path {
	f, _ := self.SplitExt()
	return f
}

func (self Path) Join(paths ...string) Path {
	paths = append([]string{string(self)}, paths...)
	return Path(filepath.Join(paths...))
}

func (self Path) JoinPath(paths ...Path) Path {
	stringPaths := make([]string, len(paths))
	for i, path := range paths {
		stringPaths[i] = string(path)
	}
	return self.Join(stringPaths...)
}

func splitAll(path string) []string {
	dir, file := filepath.Split(path)
	if file == "" {
		return []string{dir}
	}
	return append(splitAll(filepath.Dir(dir)), file)
}

func (self Path) SplitAll() (Path, []string) {
	parts := splitAll(string(self))
	curDir, err := os.Getwd()

	if err == nil {
		loc := parts[0]
		for i, part := range parts {
			if loc == curDir {
				return Path(loc), parts[i:]
			}
			loc = filepath.Join(loc, part)
		}
	}

	return Path(parts[0]), parts[1:]
}

func (self Path) RelPath(start string) (Path, error) {
	path, err := filepath.Rel(start, string(self))
	return Path(path), err
}

func (self Path) RelPathTo(dest string) (Path, error) {
	path, err := filepath.Rel(string(self), dest)
	return Path(path), err
}

type filterFunc func(file os.FileInfo) bool

func (self Path) FilterDir(filter filterFunc) (paths []Path, err error) {
	files, err := ioutil.ReadDir(string(self))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filter == nil || filter(file) {
			paths = append(paths, Path(file.Name()))
		}
	}
	return
}

func (self Path) ListDir() ([]Path, error) {
	return self.FilterDir(nil)
}

func (self Path) ListDirPattern(pattern string) ([]Path, error) {
	_, err := filepath.Match(pattern, "")
	if err != nil {
		return nil, err
	}
	return self.FilterDir(func(file os.FileInfo) bool {
		matched, _ := filepath.Match(pattern, file.Name())
		return matched
	})
}

func (self Path) Files() ([]Path, error) {
	return self.FilterDir(func(file os.FileInfo) bool {
		return !file.IsDir()
	})
}

func (self Path) FilesPattern(pattern string) ([]Path, error) {
	_, err := filepath.Match(pattern, "")
	if err != nil {
		return nil, err
	}
	return self.FilterDir(func(file os.FileInfo) bool {
		matched, _ := filepath.Match(pattern, file.Name())
		return matched && !file.IsDir()
	})
}

func (self Path) Dirs() ([]Path, error) {
	return self.FilterDir(func(file os.FileInfo) bool {
		return file.IsDir()
	})
}

func (self Path) DirsPattern(pattern string) ([]Path, error) {
	_, err := filepath.Match(pattern, "")
	if err != nil {
		return nil, err
	}
	return self.FilterDir(func(file os.FileInfo) bool {
		matched, _ := filepath.Match(pattern, file.Name())
		return matched && file.IsDir()
	})
}

func (self Path) Stat() (os.FileInfo, error) {
	return os.Stat(string(self))
}

func (self Path) LStat() (os.FileInfo, error) {
	return os.Stat(string(self))
}

func (self Path) isDir() (bool, error) {
	fileInfo, err := self.Stat()
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

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

	if err == nil || ok && err.Err == syscall.ENOENT {
		return nil
	}
	return err
}

func (self Path) RemoveTree() error {
	_, err := self.Stat()
	if err != nil {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			e.Op = "remove"
			return e
		}
	}
	return os.RemoveAll(string(self))
}

func (self Path) RemoveTreeP() error {
	return os.RemoveAll(string(self))
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

func Getwd() (Path, error) {
	p, err := os.Getwd()
	return Path(p), err
}

func (self Path) Chdir() error {
	return os.Chdir(string(self))
}

func (self Path) Cd() error {
	return self.Chdir()
}

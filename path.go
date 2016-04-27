package path

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Path string

func (self Path) Abs() (Path, error) {
	path, err := filepath.Abs(string(self))
	return Path(path), err
}

func (self Path) NormPath() Path {
	return Path(filepath.Clean(string(self)))
}

func (self Path) RealPath() (Path, error) {
	path, err := filepath.EvalSymlinks(string(self))
	return Path(path), err
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

const (
	Ignore = iota
	Warn
	Strict
)

func (self Path) Walk(errors int) error {
	warningError := WarningError([]Path{})
	var cb WalkFunc

	switch errors {
	case Ignore:
		cb = func(_ Path, _ error) error { return filepath.SkipDir }
	case Strict:
		cb = func(_ Path, err error) error { return err }
	case Warn:
		cb = func(path Path, _ error) error {
			warningError = append(warningError, path)
			return filepath.SkipDir
		}
	}
	err := self.WalkCb(cb)

	if errors == Warn {
		return &warningError
	}
	return err
}

type WalkFunc func(path Path, err error) error

func (self Path) WalkCb(walkFn WalkFunc) error {
	return filepath.Walk(
		string(self),
		func(path string, _ os.FileInfo, err error) error {
			return walkFn(Path(path), err)
		},
	)
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

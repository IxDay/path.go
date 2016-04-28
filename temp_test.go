package path

import (
	"os"
	"syscall"
	"testing"
)

func TestTempDir(t *testing.T) {
	called := false

	TempDir(func(p Path) {
		if _, err := p.Stat(); IsNotExist(err) {
			t.Errorf("TempDir() => %q, does not exist", p)
		}
		called = true
	})
	if !called {
		t.Errorf("TempDir() => cb, not called")
	}
}

func TestTempDirNamed(t *testing.T) {
	dir, prefix := "/tmp", "prefix"
	prefixLen := len(prefix)

	TempDirNamed(dir, prefix, func(p Path) {
		cbDir, cbPrefix := p.DirName(), p.BaseName()[:prefixLen]
		if string(cbDir) != dir || string(cbPrefix) != prefix {
			t.Errorf(
				"TempDirNamed(%q, %q) => cb, got prefix: %q dirname: %q",
				dir, prefix, cbDir, cbPrefix,
			)
		}
	})
}

func TestTempDirNamedPermError(t *testing.T) {
	err := TempDirNamed("/root", "", func(_ Path) {}).(*os.PathError)
	if err.Err != syscall.EACCES {
		t.Errorf("TempDir() => %q, want: %q", err.Err, syscall.EACCES)
	}
}

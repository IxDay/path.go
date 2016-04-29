package path

import (
	"os"
	"testing"
)

func TestTempFile(t *testing.T) {
	// test prefixed tempfile
	TmpDir_(func(p Path) {
		if f, _ := p.TempFile("foo"); string(f.BaseName())[:3] != "foo" {
			t.Errorf(
				"TempFile(foo) => %q, must return a file prefixed with 'foo'", f.Path,
			)
		}
	})
}

func TestTempFile_(t *testing.T) {
	// test normal behaviour
	TmpDir_(func(p Path) {
		if _, err := p.TempFile_(); err != nil {
			t.Errorf("TempFile(foo) => f, %q, must not return an error", err)
		}
	})

	// test access error
	if _, err := Path("/root").TempFile_(); !os.IsPermission(err) {
		t.Errorf(
			"Path('/root').TempFile() must return permission error, got %q", err,
		)
	}
}
func TestPathTempFile(t *testing.T) {
	// test prefixed tempfile
	TmpDir_(func(p Path) {
		if f, _ := p.TempFile("foo"); string(f.BaseName())[:3] != "foo" {
			t.Errorf(
				"TempFile(foo) => %q, must return a file prefixed with 'foo'", f.Path,
			)
		}
	})
}

func TestPathTempFile_(t *testing.T) {
	// test normal behaviour
	TmpDir_(func(p Path) {
		if _, err := p.TempFile_(); err != nil {
			t.Errorf("TempFile(foo) => f, %q, must not return an error", err)
		}
	})

	// test access error
	if _, err := Path("/root").TempFile_(); !os.IsPermission(err) {
		t.Errorf(
			"Path('/root').TempFile() must return permission error, got %q", err,
		)
	}
}

func TestTmpDir(t *testing.T) {
	// test normal behavior
	dir, prefix := "/tmp", "prefix"
	prefixLen := len(prefix)

	TmpDir(dir, prefix, func(p Path) {
		cbDir, cbPrefix := p.DirName(), p.BaseName()[:prefixLen]
		if string(cbDir) != dir || string(cbPrefix) != prefix {
			t.Errorf(
				"TmpDir(%q, %q) => cb, got prefix: %q dirname: %q",
				dir, prefix, cbDir, cbPrefix,
			)
		}
	})

	// test permission denied
	err := TmpDir("/root", "", func(_ Path) {}).(*os.PathError)
	if err == nil || !os.IsPermission(err.Err) {
		t.Errorf("TempDir() => %q, must return no access error", err)
	}
}

func TestTmpDir_(t *testing.T) {
	called := false

	TmpDir_(func(p Path) {
		if _, err := p.Stat(); os.IsNotExist(err) {
			t.Errorf("TempDir() => %q, does not exist", p)
		}
		called = true
	})
	if !called {
		t.Errorf("TempDir() => cb, not called")
	}
}

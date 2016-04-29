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
			t.Errorf("Path().TempFile() => p, %q, must not return an error", err)
		}
	})

	// test access error
	if _, err := Path("/root").TempFile_(); !os.IsPermission(err) {
		t.Errorf(
			"Path('/root').TempFile() must return permission error, got %q", err,
		)
	}

	// test TempFile alias
	if p, err := TempFile_(); err != nil {
		t.Errorf("TempFile_() => p, %q, must not return an error", err)
	} else {
		p.Remove()
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

func TestPathTmpDir(t *testing.T) {
	called := false

	TmpDir_(func(p Path) {
		// test prefix
		p.TmpDir("foo", func(p Path) {
			p = p.BaseName()
			if string(p)[:3] != "foo" {
				t.Errorf("must return a dir prefixed with 'foo', got %q", p)
			}
		})

		// test no prefix
		p.TmpDir_(func(_ Path) { called = true })
	})
	if !called {
		t.Errorf("cb, not called")
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
	err := TmpDir("/root", "", func(_ Path) {})
	if err == nil || !os.IsPermission(err) {
		t.Errorf("TmpDir() => %q, must return no access error", err)
	}
}

func TestTmpDir_(t *testing.T) {
	called := false

	TmpDir_(func(p Path) {
		if _, err := p.Stat(); os.IsNotExist(err) {
			t.Errorf("TmpDir_() => %q, does not exist", p)
		}
		called = true
	})
	if !called {
		t.Errorf("TmpDir_() => cb, not called")
	}
}

func TestTempDir_(t *testing.T) {
	if p, err := TempDir_(); err != nil {
		t.Errorf("TempDir_() => %q, must not return an error", err)
	} else {
		p.RemoveTree()
	}
}

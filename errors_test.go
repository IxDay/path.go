package path

import (
	"os"
	"syscall"
	"testing"
)

func TestWarningError(t *testing.T) {
	warns := WarningError([]Path{
		Path("."),
		Path("foo"),
	})
	errStr := warns.Error()
	if errStr != "2 directories skipped" {
		t.Errorf("Skipping directory error not accurate, got: %q", errStr)
	}
}

func TestIsAFileError(t *testing.T) {
	errors := []struct {
		in  error
		out bool
	}{
		{nil, false},
		{&os.PathError{"foo", "baz", syscall.ENOENT}, false},
		{&os.PathError{"foo", "baz", NotAFileError}, true},
		{&os.LinkError{"foo", "bar", "baz", syscall.ENOENT}, false},
		{&os.LinkError{"foo", "bar", "baz", NotAFileError}, true},
		{os.NewSyscallError("foo", syscall.ENOENT), false},
		{os.NewSyscallError("foo", NotAFileError), true},
	}

	for _, err := range errors {
		if res := IsAFileError(err.in); res != err.out {
			t.Errorf("IsAFileError(%q) => %q, want %q", err.in, res, err.out)
		}
	}
}

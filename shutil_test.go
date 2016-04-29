package path

import (
	"os"
	"testing"
)

func TestRemove(t *testing.T) {

	// test normal behavior
	TmpDir_(func(p Path) {
		f, _ := p.TempFile_()
		if err := f.Remove(); err != nil {
			t.Errorf("Remove() => %q, must not return an error", err)
		}
	})

	// test no access behavior
	if err := Path("/root/foo").Remove(); !os.IsPermission(err) {
		t.Errorf(
			"Path('/root/foo').Remove() must return permission error, got %q", err,
		)
	}

	// test delete dir behavior
	TmpDir_(func(p Path) {
		d, _ := p.TempDir_()
		if err := d.Remove(); !IsAFileError(err) {
			t.Errorf(
				"Path('/some/dir').Remove() must return a is not file error, got %q",
				err,
			)
		}
	})
}

func TestRemoveP(t *testing.T) {
	var f *File
	TmpDir_(func(p Path) {
		// test normal behavior
		f, _ = p.TempFile_()
		if err := f.RemoveP(); err != nil {
			t.Errorf("Removing file must not failed, got %q", err)
		}

		// test with non existing file
		f, _ = p.TempFile_()
		f.Remove()
		if err := f.RemoveP(); err != nil {
			t.Errorf("Removing non existing file must not return error, got %q", err)
		}

		// test with no access privilege
		if err := Path("/root/foo").RemoveP(); !os.IsPermission(err) {
			t.Errorf(
				"Path('/root/foo').RemoveP() must return permission error, got %q",
				err,
			)
		}
	})
}

func TestRemoveTree(t *testing.T) {
	TmpDir_(func(p Path) {
		// test normal behavior
		d, _ := TempDir_()
		if err := d.RemoveTree(); err != nil {
			t.Errorf("Removing a dir must not return an error, got %q", err)
		}

		// test removing non existing dir
		if err := d.RemoveTree(); !os.IsNotExist(err) {
			t.Errorf("Removing a non existing dir must return an error, got %q", err)
		}
	})
}

func TestRemoveTreeP(t *testing.T) {
	TmpDir_(func(p Path) {
		// test normal behavior
		d, _ := TempDir_()
		if err := d.RemoveTreeP(); err != nil {
			t.Errorf("Removing a dir must not return an error, got %q", err)
		}

		// test removing non existing dir
		if err := d.RemoveTreeP(); err != nil {
			t.Errorf(
				"Removing a non existing dir must not return an error, got %q", err,
			)
		}
	})
}

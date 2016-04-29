package path

import (
	"os"
	"reflect"
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
		err := d.Remove().(*os.PathError)
		if !reflect.DeepEqual(err.Err, NotAFileError) {
			t.Errorf(
				"Path('/some/dir').Remove() must return a is not file error, got %q",
				err,
			)
		}
	})
}

package path

import (
	"os"
	"reflect"
	"testing"
)

var stattests = []string{
	".",
	"/root/foo",
}

func TestStat(t *testing.T) {
	for _, tt := range stattests {

		osStatFileInfo, osStatError := os.Stat(tt)
		pathStatFileInfo, pathStatError := Path(tt).Stat()

		isValid := reflect.DeepEqual(osStatFileInfo, pathStatFileInfo)
		isValid = isValid && reflect.DeepEqual(osStatError, pathStatError)

		if !isValid {
			t.Errorf(
				"Path(%q).Stat() => \n%q %q, \nwant:\n%q %q", tt,
				pathStatFileInfo, pathStatError,
				osStatFileInfo, osStatError,
			)
		}
	}
}

func TestGetwd(t *testing.T) {
	osGetwdDir, osGetwdErr := os.Getwd()
	pathGetwdDir, pathGetwdErr := Getwd()

	isValid := osGetwdDir == string(pathGetwdDir)
	isValid = isValid && reflect.DeepEqual(osGetwdErr, pathGetwdErr)

	if !isValid {
		t.Errorf(
			"Getwd() => \n%q %q \nwant: \n%q %q",
			pathGetwdDir, pathGetwdErr, osGetwdDir, osGetwdErr,
		)
	}
}

func TestCd(t *testing.T) {
	TmpDir_(func(p Path) {
		if cwd, _ := Getwd(); p != cwd {
			t.Errorf("Path(%q).Cd(), do not move to correct directory", p)
		}
	})
}

func TestAbs(t *testing.T) {
	TmpDir_(func(p Path) {
		abs, _ := Path(".").Abs()
		if p != abs {
			t.Errorf("Path(.).Abs() => %q, want: %q", abs, p)
		}
	})
}

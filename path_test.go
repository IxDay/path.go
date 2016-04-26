package path

import (
	"fmt"
	"os"
	"testing"
)

func TestAbs(t *testing.T) {
	_cwd, err := os.Getwd()

	cwd := Path(_cwd)
	if err != nil {
		t.Errorf("An error occured: %s", err)
		return
	}

	fmt.Printf("%s\n", cwd)
}

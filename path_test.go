package path

import (
	"fmt"
	//"os"
	"testing"
)

//func TestAbs(t *testing.T) {
//	_cwd, err := os.Getwd()
//
//	cwd := Path(_cwd)
//	if err != nil {
//		t.Errorf("An error occured: %s", err)
//		return
//	}
//
//	fmt.Printf("%s\n", cwd)
//}

func TestAbs(t *testing.T) {
	//_cwd, _ := os.Getwd()

	//cwd := Path(_cwd)
	//cwd = cwd.Join("/toto")
	//
	//err := cwd.RemoveTree()
	//fmt.Printf("%v\n", err)

	TempDir(func(p Path) {
		fmt.Printf("%s\n", p)
	})
}

package util

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestPwdMd5(t *testing.T) {
	s := GenPasswordString("aaa", "aaa")
	fmt.Println("result=>", s)
}

func TestPath(t *testing.T) {
	fileName := filepath.Base("./aaa/ccc/abc.txt")
	fileName1 := filepath.Base("/aaa/bbb")
	t.Log("is file:", FileExists("./utils.go"))
	t.Log("base:", fileName, fileName1)
}
func TestBreakLoop(t *testing.T) {
end:
	for {
		fmt.Println("loop.")
		break end
	}
	fmt.Println("end ...")
}

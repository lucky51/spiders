package test

import (
	"errors"
	"fmt"
	"testing"
)

func TestRtnErrorWhenPicnic(t *testing.T) {
	_, err := invoke()
	t.Logf("TestRtnErrorWhenPicnic error:%v", err)
}

func invoke() (str string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	panic(errors.New("panic error"))
}

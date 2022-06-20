package util_test

import (
	"github.com/g-knife/common/util"
	"testing"
)

func TestGetGoPath(t *testing.T) {
	p, err := util.GetGoPath()
	if err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(p + "wxf")
	}
}

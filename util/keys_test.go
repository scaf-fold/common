package util_test

import (
	"github.com/scaf-fold/common/util"
	"testing"
	"time"
)

func TestRSAGenerateKey(t *testing.T) {
	t1 := time.Now()
	defer func() {
		t.Log(time.Since(t1))
	}()
	path := "./"
	err := util.GenerateRSAPPr(512, path)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success", path)
}

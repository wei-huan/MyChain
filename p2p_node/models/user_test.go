package models

import (
	"fmt"
	"testing"
)

func TestLogin(t *testing.T) {
	u, err := Login("luda")
	if err != nil {
		t.FailNow()
	}
	// 打印结构体
	fmt.Printf("%v", u)
}

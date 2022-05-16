package models

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenerateBlock(t *testing.T) {
	b := GenerateBlock("0", "TEST_0", 0)
	fmt.Println(b)
}

//测试 Block json 格式转换
func TestFormatBlock(t *testing.T) {
	b := GenerateBlock("0", "TEST_0", 0)
	// Marshal 即将结构体格式的 block 转换成 json 格式的 block
	jb, _ := json.Marshal(b)
	// FormatBlock 即将 json 格式的 block 转换成结构体格式的 block
	bf, err := FormatBlock(jb)
	if err != nil {
		fmt.Println(err)
	}
	if *bf != *b {
		t.Fail()
	}
	fmt.Println("bf == b")
}

//测试 Block 是否合法
func TestBlock_IsValid(t *testing.T) {
	b := GenerateBlock("0", "TEST_0", 0)
	b1 := GenerateBlock(b.Hash, "TEST_1", 1)
	b2 := GenerateBlock(b.Hash, "TEST_2", 2)
	fmt.Println("b1 is behind b:", b1.IsValid(b))
	fmt.Println("b2 is behind b1:", b2.IsValid(b1))
	fmt.Println("b2 is behind b:", b2.IsValid(b))
}

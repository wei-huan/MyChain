package models

import (
	"encoding/json"

	"github.com/Blockchain-CN/keygen"
)

// Trans transaction struct
type Trans struct {
	// Account public key
	Account string `json:"account"`
	// Cipher encrypt result
	Cipher string `json:"cipher"`
	// Transaction result
	Transaction string `json:"transaction"`
}

// IsVaild return if a trans is legal.
func (t *Trans) IsVaild() error {
	return keygen.Verify(t.Account, t.Cipher, []byte(t.Transaction))
}

// FormatTrans format []byte to a trans object.
// 格式化为json
func FormatTrans(b []byte) (*Trans, error) {
	trans := &Trans{}
	err := json.Unmarshal(b, trans)
	if err != nil {
		return nil, err
	}
	return trans, nil
}

// GenerateTransWithID generate a trans using user's ID.
// 根据登陆用户名 id 查找公钥加密data信息，并返回 base64 编码后结果
func GenerateTransWithID(id, data string) (*Trans, error) {
	// a 是账号公钥 base64 编码后结果
	// c 是data经过公钥加密后并 base64 编码后结果
	a, c, err := keygen.Signature(id, []byte(data))
	if err != nil {
		return nil, err
	}
	return &Trans{
		Account:     a,
		Cipher:      c,
		Transaction: data}, nil
}

// GenerateTransWithKey generate a trans using the key.
// 根据公钥 pb 直接加密data信息，并返回 base64 编码后结果
func GenerateTransWithKey(pb, pv, data string) (*Trans, error) {
	c, err := keygen.Signature2(pv, []byte(data))
	if err != nil {
		return nil, err
	}
	return &Trans{
		Account:     pb,
		Cipher:      c,
		Transaction: data}, nil
}

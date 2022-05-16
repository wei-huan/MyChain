package models

import (
	"path"

	"github.com/Blockchain-CN/keygen"

	"MyChain/common"
)

// User struct.
type User struct {
	Name    string
	Path    string
	Public  string
	Private string
}

// Login allow user to login and get their key.
func Login(name string) (*User, error) {
	// 根据用户名找用户目录
	uPath := keygen.GetUserPath(name)
	// 从用户目录的私钥文件和公钥文件获取公钥私钥
	pvKeyPath := path.Join(uPath, "private.pem")
	pbKeyPath := path.Join(uPath, "public.pem")
	pv, errv := keygen.GetKeyMd5(pvKeyPath)
	pb, errb := keygen.GetKeyMd5(pbKeyPath)
	// 生成结构体
	if errv == nil && errb == nil {
		return &User{
			Name:    name,
			Path:    uPath,
			Public:  pb,
			Private: pv}, nil
	}
	// 错误检查，RSA 部分上次陈诚讲过了，我就不赘述
	if err := keygen.GenRsaKey(common.RSADefaultLenth, name); err != nil {
		return nil, err
	}
	pv, err := keygen.GetKeyMd5(pvKeyPath)
	if err != nil {
		return nil, err
	}
	pb, err = keygen.GetKeyMd5(pbKeyPath)
	if err != nil {
		return nil, err
	}
	return &User{
		Name:    name,
		Path:    uPath,
		Public:  pb,
		Private: pv}, nil
}

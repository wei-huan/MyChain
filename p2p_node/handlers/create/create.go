package create

import (
	"MyChain/common"
	idl "MyChain/idls/create"
	"MyChain/single"
)

// GenerateBlock create a new block and spread it.
func GenerateBlock(req *idl.CRequest) *idl.CResponse {
	resp := idl.NewCResponseIDL()
	resp.Errno = common.Success
	resp.Msg = common.ErrMap[common.Success]
	single.DataQueue <- req
	return resp
}

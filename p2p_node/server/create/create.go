package create

import (
	handler "MyChain/handlers/create"
	idl "MyChain/idls/create"
)

// CController ...
type CController struct {
}

// GenIdl ...
func (c *CController) GenIdl() interface{} {
	return idl.NewCRequestIDL()
}

// Do ...
func (c *CController) Do(req interface{}) interface{} {
	r := req.(*idl.CRequest)
	return handler.GenerateBlock(r)
}

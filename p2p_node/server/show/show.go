package show

import (
	handler "MyChain/handlers/show"
	idl "MyChain/idls/show"
)

// JController ...
type SController struct {
}

// GenIdl ...
func (c *SController) GenIdl() interface{} {
	return idl.NewJRequest()
}

// Do ...
func (c *SController) Do(req interface{}) interface{} {
	r := req.(*idl.SRequest)
	return handler.Show(r)
}

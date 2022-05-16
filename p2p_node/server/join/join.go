package join

import (
	idl "MyChain/idls/join"
)

// JController ...
type JController struct {
}

// GenIdl ...
func (c *JController) GenIdl() interface{} {
	return idl.NewJRequest()
}

// Do ...
func (c *JController) Do(req interface{}) interface{} {
	return nil
}

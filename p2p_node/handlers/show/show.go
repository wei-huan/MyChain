package show

import (
	idl "MyChain/idls/show"
	"MyChain/models"
	"MyChain/single"
)

type result struct {
}

// Show join to the blockchain system by connect to a peer
func Show(req *idl.SRequest) *idl.SResponse {
	resp := idl.NewJResponse()
	single := single.GetProtocal()
	if req.Chain {
		resp.Chain = models.FetchChain()
	}
	if req.Peer {
		resp.Peer = single.Proto.GetRouter().FetchPeers()
	}
	return resp
}

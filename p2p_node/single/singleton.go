package single

import (
	"encoding/json"
	"time"

	idl "MyChain/idls/create"
	"MyChain/models"
	"MyChain/p2p"
)

var (
	singleton *p2p.Server
	// DataQueue data channel
	DataQueue chan *idl.CRequest
)

func GetProtocal() *p2p.Server {
	return singleton
}

// InitPto init the default protocal object
func InitPto(to time.Duration) {
	r1 := p2p.NewRouter(to)
	p1 := p2p.NewProtocal(*r1, to)
	s1 := p2p.NewServer("luda", *p1, to)
	singleton = s1
	DataQueue = make(chan *idl.CRequest, 100)
	go BlockPublisher()
	go s1.PeerAndServe()
}

// BlockPublisher loop to publish the block
func BlockPublisher() {
	for {
		select {
		case ud := <-DataQueue:
			// get user object
			user, err := models.Login(ud.Name)
			if err != nil {
				return
			}

			// get trans object
			trans, err := models.GenerateTransWithKey(user.Public, user.Private, ud.Data)
			if err != nil {
				return
			}
			transStr, err := json.Marshal(trans)
			if err != nil {
				return
			}

			// append a block to the chain until succeed
			for {
				// get block object
				block := models.GenerateBlock(models.GetChainTail().Hash, string(transStr), models.GetChainLen())

				// add to blockchain
				err = models.AppendChain(block)
				if err == nil {
					singleton.Proto.Spread(block, singleton.HostName, singleton.HostAddr)
					break
				}
			}
		}
	}
}

package p2p

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"MyChain/common"
	"MyChain/models"

	dhash "github.com/Blockchain-CN/sha256"
)

const (
	// 连接请求
	ConnectReq = "connectreq"
	// 获取一个
	GetReq = "getreq"
	// 批量获取
	FetchReq = "fetchreq"
	// 同步更新
	NoticeReq = "noticereq"
	// 连接请求返回
	ConnectResp = "connectresp"
	// 获取一个返回
	GetResp = "getresp"
	// 批量获取返回
	FetchResp = "fetchresp"
	// 同步更新返回
	NoticeResp = "noticeresp"
	// 未知操作
	UnknownOp = "unknownop"
)

// MsgPto 协议数据格式
type MsgPto struct {
	Name      string `json:"name"`
	Operation string `json:"operation"`
	// 子协议json
	Data []byte `json:"data"`
}

type MsgGreetingReq struct {
	Addr    string `json:"add"`
	Account int    `json:"account"`
}

type Protocal struct {
	router Router
	to     time.Duration
}

func NewProtocal(r Router, to time.Duration) *Protocal {
	return &Protocal{r, to}
}

func (p *Protocal) Handle(c net.Conn, hostname, hostaddr string, msg []byte) ([]byte, error) {
	if msg == nil {
		return nil, nil
	}
	req := &MsgPto{}
	resp := &MsgPto{}
	err := json.Unmarshal(msg, req)
	if err != nil {
		return nil, Error(ErrMismatchProtocalReq)
	}
	resp.Name = hostname
	switch req.Operation {
	case RequireBlock:
		err = p.router.AddRoute(req.Name, req.Name)
		if err != nil {
			fmt.Println(err)
		}
		c, _ := json.Marshal(models.GetChainTail())
		resp.Operation = DeliveryBlock
		resp.Data = c
	case DeliveryBlock:
		dhash.StopHash()
		defer dhash.StartHash()
		block, err := models.FormatBlock(req.Data)
		if err != nil {
			return nil, Error(ErrMismatchProtocalResp)
		}
		// if the block's index is shorter or invalidate
		tailBlock := models.GetChainTail()
		if *block == *tailBlock {
			return nil, nil
		}
		if !block.IsTempValid() || block.Index <= tailBlock.Index {
			return nil, common.Error(common.ErrInvalidBlock)
		}
		// if the block can append to the tail
		if block.IsValid(tailBlock) {
			models.AppendChain(block)
			// 并需要向外广播
			go p.Spread(block, hostname, hostaddr)
			return nil, nil
		}
		// if the block's index is longer
		resp.Operation = RequireChain
	case RequireChain:
		c, _ := json.Marshal(models.FetchChain())
		resp.Operation = DeliveryChain
		resp.Data = c
	case DeliveryChain:
		dhash.StopHash()
		defer dhash.StartHash()
		chain, err := models.FormatChain(req.Data)
		if err != nil {
			return nil, Error(ErrMismatchProtocalResp)
		}
		err = models.ReplaceChain(chain)
		if err != nil {
			return nil, common.Error(common.ErrInvalidChain)
		}
		// 向外广播 models.GetChainTail()
		go p.Spread(models.GetChainTail(), hostname, hostaddr)
		return nil, nil
	default:
		fmt.Printf("@%s@report: %s operation from @%s@ finished\n", hostname, req.Operation, req.Name)
		return nil, nil
	}
	ret, err := json.Marshal(resp)
	fmt.Printf("@%s@report: %s operation from @%s@ succeed\n", hostname, req.Operation, req.Name)
	return ret, nil
}

func (p *Protocal) read(r io.Reader) ([]byte, error) {
	buf := make([]byte, defultByte)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	// read读出来的是[]byte("abcdefg"+0x00)，带一个结束符，需要去掉
	return buf[:n], nil
}

func (p *Protocal) GetRouter() Router {
	return p.router
}

func (p *Protocal) DispatchAll(msg []byte) map[string][]byte {
	return p.router.DispatchAll(msg)
}

func (p *Protocal) Dispatch(name string, msg []byte) ([]byte, error) {
	return p.router.Dispatch(name, msg)
}

func (p *Protocal) Delete(name string) error {
	return p.router.Delete(name)
}

// spread the latest block to all peers
func (p *Protocal) Spread(block *models.Block, hostname, hostaddr string) {
	blockStr, err := json.Marshal(block)
	if err != nil {
		return
	}
	req := &MsgPto{
		Name:      hostaddr,
		Operation: DeliveryBlock,
		Data:      blockStr,
	}
	reqStr, err := json.Marshal(req)
	if err != nil || reqStr == nil {
		return
	}
	peerList := p.GetRouter().FetchPeers()
	// 同步等待和所有peer通信完毕
	for k, _ := range peerList {
		wg.Add(1)
		go func(name string) {
			for reqStr != nil {
				b, err := p.Dispatch(name, reqStr)
				if err != nil {
					println("操作失败", err.Error())
					return
				}
				reqStr = nil
				reqStr, err = p.Handle(nil, hostname, hostaddr, b)
				fmt.Println(string(reqStr), err)
			}
			wg.Done()
		}(k)
	}
	wg.Wait()
}

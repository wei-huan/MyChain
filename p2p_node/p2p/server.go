package p2p

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// RequireBlock 请求最新block
	RequireBlock = "require a block"

	// Publish 发布最新block
	DeliveryBlock = "delivery a block"

	// DeliveryChain 发送整条链
	DeliveryChain = "delivery the block"

	// RequireChain 请求整条链
	RequireChain = "require the block"

	// 未知命令
	UnknownCmd = "unknown cmd"
)

var wg sync.WaitGroup

// Server p2p监听连接server
type Server struct {
	HostName string
	HostAddr string // 外网地址, 即在公网显示的地址
	Proto    Protocal
	to       time.Duration
}

// NewServer ...
func NewServer(name string, p Protocal, to time.Duration) *Server {
	return &Server{name, "", p, to}
}

func (s *Server) handler(c net.Conn) {
	defer c.Close()
	msg, err := s.read(c, s.to)
	fmt.Printf("收到请求, localhost=%s||remotehost=%s||msg=%s\n", c.LocalAddr(), c.RemoteAddr(), string(msg))
	if err != nil {
		return
	}
	resp, err := s.Proto.Handle(c, s.HostName, s.HostAddr, msg)
	if err != nil || resp == nil {
		resp = nil
	}
	c.SetWriteDeadline(time.Now().Add(s.to))
	for i := 0; i < 3; i++ {
		_, err = c.Write(resp)
		if err != nil {
			continue
		}
		fmt.Printf("发送回复, localhost=%s||remotehost=%s||msg=%s\n", c.LocalAddr(), c.RemoteAddr(), string(resp))
		return
	}
}

func (s *Server) read(r io.Reader, to time.Duration) ([]byte, error) {
	buf := make([]byte, defultByte)
	messnager := make(chan int)
	go func() {
		n, _ := r.Read(buf[:])
		messnager <- n
		close(messnager)
	}()
	select {
	case n := <-messnager:
		return buf[:n], nil
	case <-time.After(to):
		return nil, Error(ErrLocalSocketTimeout)
	}
}

// PeerAndServe 监听 peer 的链接请求
// 配对用的端口和进行P2P服务的端口是两个不同的端口
func (s *Server) PeerAndServe() error {
	rname, raddr, err := s.linkpeer()
	if err != nil {
		fmt.Println(err)
		return err
	}
	s.Proto.router.AddRoute(rname, raddr)
	fmt.Println(s.HostName + "Peers:")
	fmt.Println(s.Proto.router.FetchPeers())
	ln, err := net.Listen("tcp", s.HostAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				runtime.Gosched()
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
			}
			break
		}
		go s.handler(c)
	}
	return nil
}

/// 私有打洞服务器 IP + Port，运行了配对程序
var server_addr = "39.105.139.143:9527"

/// P2P 结点通过固定的打洞服务器寻找配对
/// 成功返回配对的用户名和 IP，注意返回的是进行 P2P 连接的IP，不是打洞用的 UDP 的 IP
func (s *Server) linkpeer() (string, string, error) {
	// conn 代表一个 udp 网络连接，实现了 Conn 和 Packet Conn 接口
	conn, err := net.Dial("udp", server_addr)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	defer conn.Close()
	if _, err = conn.Write([]byte(s.HostName)); err != nil {
		fmt.Println(err)
		return "", "", err
	}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
		return "", "", err
	}
	hostaddr, rname, raddr := parseNameAddrID(string(data[:n]))
	s.HostAddr = hostaddr
	fmt.Printf("remote peer:%s %s\n", rname, raddr)
	return rname, raddr, nil
}

/// return laddr, rname and raddr
func parseNameAddrID(addr string) (string, string, string) {
	t := strings.Split(addr, " ")
	return t[0], t[1], t[2]
}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

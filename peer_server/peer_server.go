package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

var wg = sync.WaitGroup{}
var para_num = 1
var port = 9527

/// 1个端口10个conn并发
func main() {
	//ListenUDP 创建一个接收目的地是本地地址 laddr 的 UDP 数据包的网络连接。net 必须是 "udp"、"udp4"、"udp6"；如果 laddr 端口为0，函数将选择一个当前可用的端口，可以用Listener的Addr方法获得该端口。返回的*UDPConn的ReadFrom和WriteTo方法可以用来发送和接收UDP数据包（每个包都可获得来源地址或设置目标地址）。
	//IPv4zero: 本地地址，只能作为源地址（曾用作广播地址）
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.IPv4zero, Port: port})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	//LocalAddr 返回本地网络地址
	log.Printf("本地地址：<%s> \n", listener.LocalAddr().String())
	wg.Add(para_num)
	for i := 0; i < para_num; i++ {
		go find_peers(listener)
	}
	wg.Wait()
}

type Peer struct {
	Name string
	Addr string
}

func find_peers(listener *net.UDPConn) {
	peers := make([]Peer, 0, 2)
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Println("err during read: %s", err)
			break
		}
		log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])
		peers = append(peers, Peer{Addr: remoteAddr.String(), Name: string(data[:n])})
		if len(peers) == 2 {
			log.Printf("进行UDP打洞，建立 %s <--> %s 的链接\n", peers[0], peers[1])
			//WriteToUDP通过c向地址addr发送一个数据包，b为包的有效负载，返回写入的字节。
			//WriteToUDP方***在超过一个固定的时间点之后超时，并返回一个错误。在面向数据包的连接上，写入超时是十分罕见的。
			peer0_addr, err := net.ResolveUDPAddr("udp", peers[0].Addr)
			if err != nil {
				fmt.Println(err)
				return
			}
			peer1_addr, err := net.ResolveUDPAddr("udp", peers[1].Addr)
			if err != nil {
				fmt.Println(err)
				return
			}
			listener.WriteTo([]byte(peers[0].Addr+" "+peers[1].Name+" "+peers[1].Addr), peer0_addr)
			listener.WriteTo([]byte(peers[1].Addr+" "+peers[0].Name+" "+peers[0].Addr), peer1_addr)
			log.Println("中转服务器重新开始新的配对，不影响peers间通信")
			// 把当服务器的留下继续配对
			peers = peers[:1]
			fmt.Println(peers)
		}
	}
	wg.Done()
}

/// 比较两个 ipv4 外网地址端口，如果本地的大就返回 true，否则返回 false, 优先比较端口
/// 两个地址必须是合法的 IPv4:Port 地址
func isLocalAddrBigger(local, remote string) bool {
	l := strings.Split(local, ":")
	lport, _ := strconv.Atoi(l[1])
	r := strings.Split(remote, ":")
	rport, _ := strconv.Atoi(r[1])
	if lport > rport {
		return true
	} else if lport < rport {
		return false
	}
	lip := strings.Split(l[0], ".")
	rip := strings.Split(r[0], ".")
	for i := 0; i < 4; i++ {
		li, _ := strconv.Atoi(lip[i])
		ri, _ := strconv.Atoi(rip[i])
		if li > ri {
			return true
		} else if li < ri {
			return false
		}
	}
	return false
}

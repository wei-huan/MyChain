package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"

	"MyChain/common"
	"MyChain/server"
	"MyChain/single"
)

func main() {
	var ServerPort string

	flag.StringVar(&ServerPort, "server", "", "-server=:10024")
	flag.Parse()
	if ServerPort == "" {
		useage()
		printAndDie(errors.New("Unable to get a avilable port for p2p node"))
	}

	ip := getIP()
	if ip == "" {
		printAndDie(errors.New("Unable to get a avilable ip"))
	}

	// init protocal
	single.InitPto(common.P2PTimeOut)

	// call this func will block current goroutine
	if err := server.Serve(ip + ServerPort); err != nil {
		printAndDie(err)
		return
	}
}

func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func printAndDie(err error) {
	fmt.Fprintf(os.Stderr, "init failed, err:%s", err)
	os.Exit(-1)
}

func useage() {
	fmt.Fprintf(os.Stdout, "please run \"%s --help\" and get help info\n", os.Args[0])
	os.Exit(-1)
}

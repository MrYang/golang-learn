package rpc

import (
	"github.com/MrYang/golang-learn/conf"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func StartJsonRpc() {
	addr := conf.Config().Server.JsonRpc.Listen
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("net.ResolveTCPAddr fail: %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Println("json rpc listening", addr)
	}

	server := rpc.NewServer()
	server.Register(new(Echo))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept occur error:", err)
			continue
		}

		// go jsonrpc.ServeConn(conn)
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

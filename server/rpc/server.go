package rpc

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"zz.com/go-study/conf"
)

type Server int

func (this *Server) Ping(args *string, reply *int) error {
	log.Println("recevice rpc call args:", *args)
	*reply = 2
	return nil
}

func StartRpc() {
	addr := conf.Config().Server.JsonRpc.Listen
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("net.ResolveTCPAddr fail: %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Println("rpc listening", addr)
	}

	server := rpc.NewServer()
	server.Register(new(Server))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept occur error:", err)
			continue
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

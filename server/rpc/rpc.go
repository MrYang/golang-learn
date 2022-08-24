package rpc

import (
	"github.com/MrYang/golang-learn/conf"
	"log"
	"net"
	"net/rpc"
)

type Echo int

// Ping 函数必须是导出的(首字母大写)
// 必须有两个参数，并且是导出类型或者内建类型
// 第二个参数必须是指针类型的
// 函数还要有一个返回值 error
func (echo *Echo) Ping(args *string, reply *int) error {
	log.Println("receive rpc call args:", *args)
	*reply = 2
	return nil
}

func StartRpc() {
	addr := conf.Config().Server.Rpc.Listen

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Println("rpc listening", addr)
	}

	rpc.Register(new(Echo))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept occur error:", err)
			continue
		}

		go rpc.ServeConn(conn)
	}
}

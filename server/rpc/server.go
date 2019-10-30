package rpc

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"bufio"
	"go-study/conf"
	pb "go-study/protos"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Echo int

// 函数必须是导出的(首字母大写)
// 必须有两个参数，并且是导出类型或者内建类型
// 第二个参数必须是指针类型的
// 函数还要有一个返回值 error
func (echo *Echo) Ping(args *string, reply *int) error {
	log.Println("recevice rpc call args:", *args)
	*reply = 2
	return nil
}

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

type gRpcServer struct {
}

func (s *gRpcServer) Hello(context.Context, *pb.Req) (*pb.Resp, error) {
	return &pb.Resp{Msg: "grpc resp"}, nil
}

func StartGRpc() {
	addr := conf.Config().Server.GRpc.Listen
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Println("grpc listening", addr)
	}

	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &gRpcServer{})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// StartTcp 可使用telnet 127.0.0.1 6060 测试
func StartTcp() {
	addr := conf.Config().Server.Tcp.Listen

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Println("tcp listening", addr)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept occur error:", err)
			continue
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	buf := bufio.NewReader(conn)
	log.Printf("conn: local:%s,remote:%s", conn.LocalAddr().String(), conn.RemoteAddr().String())

	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			break
		}
		if "quit" == string(line) {
			break
		}
		writer := bufio.NewWriter(conn)
		writer.WriteString("echo " + string(line) + "\n")
		writer.Flush()
	}

	log.Println(conn.RemoteAddr().String(), "quit")
}

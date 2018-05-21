package rpc

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"zz.com/go-study/conf"
	"bufio"
	"google.golang.org/grpc"
	pb "zz.com/go-study/protos"
	"google.golang.org/grpc/reflection"
	"golang.org/x/net/context"
)

type Echo int

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

func StartTcp() {
	addr := conf.Config().Server.Rpc.Listen

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
	log.Println("conn: ", conn.LocalAddr().String(), conn.RemoteAddr().String())

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

	log.Println(conn.LocalAddr().String(), conn.RemoteAddr().String(), "quit")
}

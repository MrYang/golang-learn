package rpc

import (
	"github.com/MrYang/golang-learn/conf"
	pb "github.com/MrYang/golang-learn/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

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

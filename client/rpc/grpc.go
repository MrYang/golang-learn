package rpc

import (
	pb "github.com/MrYang/golang-learn/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func CallGRpc(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil
	}
	defer conn.Close()
	c := pb.NewHelloClient(conn)

	req := &pb.Req{
		Id:     1,
		Name:   "yxb",
		Age:    0,
		Gender: pb.Req_MALE,
	}

	resp, err := c.Hello(context.Background(), req)
	if err != nil {
		return nil
	}
	log.Printf("hello: %s", resp.Msg)

	return nil
}

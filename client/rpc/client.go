package rpc

import (
	"log"
	"math"
	"net/rpc"
	"sync"
	"time"
	"net/rpc/jsonrpc"
	"net"
	"google.golang.org/grpc"
	pb "zz.com/go-study/protos"
	"context"
)

type ConnRpcClient struct {
	sync.Mutex
	rpcClient        *rpc.Client
	RpcServerAddress string
	Timeout          time.Duration
}

func (crc *ConnRpcClient) close() {
	if crc.rpcClient != nil {
		crc.rpcClient.Close()
		crc.rpcClient = nil
	}
}

func (crc *ConnRpcClient) servConn() error {
	if crc.rpcClient != nil {
		return nil
	}

	retry := 1

	for {
		if crc.rpcClient != nil {
			return nil
		}

		conn, err := net.DialTimeout("tcp", crc.RpcServerAddress, crc.Timeout)
		if err != nil {
			return err
		}

		crc.rpcClient = jsonrpc.NewClient(conn)

		if err != nil {
			log.Printf("dial %s fail: %v", crc.RpcServerAddress, err)
			if retry > 3 {
				return err
			}

			time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
			retry++
			continue
		}
		return err
	}
}

func (crc *ConnRpcClient) Call(method string, args interface{}, reply interface{}) error {
	crc.Lock()
	defer crc.Unlock()

	err := crc.servConn()
	if err != nil {
		return err
	}

	timeout := time.Duration(10 * time.Second)
	done := make(chan error, 1)

	go func() {
		err := crc.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		log.Printf("[WARN] rpc call timeout %v => %v", crc.rpcClient, crc.RpcServerAddress)
		crc.close()
	case err := <-done:
		if err != nil {
			crc.close()
			return err
		}
	}
	return nil
}

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

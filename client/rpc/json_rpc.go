package rpc

import (
	"log"
	"math"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"
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

func (crc *ConnRpcClient) getRpcClient() error {
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

func (crc *ConnRpcClient) CallJsonRpc(method string, args interface{}, reply interface{}) error {
	crc.Lock()
	defer crc.Unlock()

	err := crc.getRpcClient()
	if err != nil {
		return err
	}

	timeout := 10 * time.Second
	done := make(chan error, 1)

	go func() {
		err := crc.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		log.Printf("[WARN] json rpc call timeout %v => %v", crc.rpcClient, crc.RpcServerAddress)
		crc.close()
	case err := <-done:
		if err != nil {
			crc.close()
			return err
		}
	}
	return nil
}

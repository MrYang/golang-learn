package rpc

import (
	"log"
	"math"
	"net/rpc"
	"sync"
	"time"

	"github.com/toolkits/net"
)

type ConnRpcClient struct {
	sync.Mutex
	rpcClient        *rpc.Client
	RpcServerAddress string
	Timeout          time.Duration
}

func (this *ConnRpcClient) close() {
	if this.rpcClient != nil {
		this.rpcClient.Close()
		this.rpcClient = nil
	}
}

func (this *ConnRpcClient) servConn() error {
	if this.rpcClient != nil {
		return nil
	}

	var err error
	retry := 1

	for {
		if this.rpcClient != nil {
			return nil
		}

		this.rpcClient, err = net.JsonRpcClient("tcp", this.RpcServerAddress, this.Timeout)

		if err != nil {
			log.Printf("dial %s fail: %v", this.RpcServerAddress, err)
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

func (this *ConnRpcClient) Call(method string, args interface{}, reply interface{}) error {
	this.Lock()
	defer this.Unlock()

	err := this.servConn()
	if err != nil {
		return err
	}

	timeout := time.Duration(50 * time.Second)
	done := make(chan error, 1)

	go func() {
		err := this.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		log.Printf("[WARN] rpc call timeout %v => %v", this.rpcClient, this.RpcServerAddress)
		this.close()
	case err := <-done:
		if err != nil {
			this.close()
			return err
		}
	}
	return nil
}

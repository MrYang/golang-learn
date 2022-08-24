package rpc

import (
	"log"
	"net/rpc"
	"time"
)

func CallRpc(addr string, method string, args interface{}, reply interface{}) error {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return err
	}

	timeout := 10 * time.Second
	done := make(chan error, 1)

	go func() {
		err := client.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		log.Printf("[WARN] rpc call timeout %v => %v", client, addr)
		client.Close()
	case err := <-done:
		if err != nil {
			client.Close()
			return err
		}
	}

	return nil
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	crpc "zz.com/go-study/client/rpc"
	"zz.com/go-study/conf"
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
// ./go-study -c cfg.json 2>app.log
// go run main.go

const v = "0.0.1"

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	cfg := flag.String("c", "../conf/cfg.json", "cfg json")
	flag.Parse()

	if *version {
		fmt.Println(v)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	conf.ParseConfig(*cfg)

	rpcAddr := conf.Config().Client.JsonRpc
	gRpcAddr := conf.Config().Client.GRpc

	client := &crpc.ConnRpcClient{
		RpcServerAddress: rpcAddr,
		Timeout:          time.Duration(5 * time.Second),
	}

	args := "query"
	var reply int

	err := client.Call("Echo.Ping", &args, &reply)

	if err != nil {
		log.Println("rpc call error %v", err)
	} else {
		log.Println("rpc call result:", reply)
	}

	crpc.CallGRpc(gRpcAddr)
}

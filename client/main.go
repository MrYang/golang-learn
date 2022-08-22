package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	crpc "github.com/MrYang/golang-learn/client/rpc"
	"github.com/MrYang/golang-learn/conf"
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
// ./golang-learn -c cfg.json 2>app.log
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

	jsonRpcAddr := conf.Config().Client.JsonRpc
	gRpcAddr := conf.Config().Client.GRpc
	rpcAddr := conf.Config().Client.Rpc

	client := &crpc.ConnRpcClient{
		RpcServerAddress: jsonRpcAddr,
		Timeout:          5 * time.Second,
	}

	args := "json rpc query"
	var reply int
	err := client.Call("Echo.Ping", &args, &reply)
	if err != nil {
		log.Printf("json rpc call error %v", err)
	} else {
		log.Println("json rpc call result:", reply)
	}

	args2 := "rpc query"
	var reply2 int
	err = crpc.CallRpc(rpcAddr, "Echo.Ping", &args2, &reply2)
	if err != nil {
		log.Printf("rpc call error %v", err)
	} else {
		log.Println("rpc call result:", reply2)
	}

	crpc.CallGRpc(gRpcAddr)
}

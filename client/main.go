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

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	cfg := flag.String("c", "../conf/cfg.json", "cfg json")
	flag.Parse()

	if *version {
		fmt.Println("0.0.1")
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	conf.ParseConfig(*cfg)

	rpcAddr := conf.Config().Client.JsonRpc

	client := &crpc.ConnRpcClient{
		RpcServerAddress: rpcAddr,
		Timeout:          time.Duration(5 * time.Second),
	}

	args := "query"
	var reply int

	err := client.Call("Server.Ping", &args, &reply)

	if err != nil {
		log.Fatalf("rpc call error %v", err)
	} else {
		log.Println("rpc call result:", reply)
	}
}

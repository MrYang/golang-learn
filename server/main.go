package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"zz.com/go-study/conf"
	"zz.com/go-study/server/db"
	"zz.com/go-study/server/http"
	srpc "zz.com/go-study/server/rpc"
	"os/signal"
	"syscall"
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
// ./go-study -c cfg.json 2>app.log
// go run main.go

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())

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

	db.Init(conf.Config().Common.Database)

	go srpc.StartJsonRpc()
	go srpc.StartTcp()
	go http.Start()
	go srpc.StartGRpc()

	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)

	select {
	case s := <-sg:
		log.Println("got signal", s)
	}

	log.Println("server is stopping...")
}

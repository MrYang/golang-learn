package main

import (
	"flag"
	"fmt"
	"os"

	"zz.com/go-study/conf"
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
// ./go-study -c cfg.json 2>app.log
// go run main.go

func main() {
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	cfg := flag.String("c", "conf/cfg.json", "cfg json")
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

	print(conf.Config().Redis.Addr)
}

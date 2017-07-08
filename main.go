package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"zz.com/go-study/conf"
	"zz.com/go-study/server/db"
	"zz.com/go-study/server/http"
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
// ./go-study -c cfg.json 2>app.log
// go run main.go

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

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

	log.Println(conf.Config().Redis.Addr)
	for _, addr := range conf.Config().JsonRpc.Addrs {
		log.Println(addr)
	}

	db.Init(conf.Config().Database)
	users, _ := db.Query()
	for _, u := range users {
		log.Println(u.ID, u.Username, u.Password, u.CreateDate.Format("2006-01-02 15:04:05"))
	}

	http.Init()
}

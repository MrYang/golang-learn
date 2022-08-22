package tcp

import (
	"bufio"
	"github.com/MrYang/golang-learn/conf"
	"log"
	"net"
)

// StartTcp 可使用telnet 127.0.0.1 6060 测试
func StartTcp() {
	addr := conf.Config().Server.Tcp.Listen

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Println("tcp listening", addr)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept occur error:", err)
			continue
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	buf := bufio.NewReader(conn)
	log.Printf("conn: local:%s,remote:%s", conn.LocalAddr().String(), conn.RemoteAddr().String())

	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			break
		}
		if "quit" == string(line) {
			break
		}
		writer := bufio.NewWriter(conn)
		writer.WriteString("echo " + string(line) + "\n")
		writer.Flush()
	}

	log.Println(conn.RemoteAddr().String(), "quit")
}

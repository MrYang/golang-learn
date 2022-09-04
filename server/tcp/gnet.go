package tcp

import (
	"fmt"
	"github.com/MrYang/golang-learn/conf"
	"github.com/panjf2000/gnet/v2"
	"log"
)

type EchoServer struct {
	gnet.BuiltinEventEngine

	eng gnet.Engine
}

func (s *EchoServer) OnBoot(engine gnet.Engine) (action gnet.Action) {
	s.eng = engine
	log.Println("echo server start")
	return
}

func (s *EchoServer) OnOpen(conn gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("conn %s open", conn.RemoteAddr())
	return
}

func (s *EchoServer) OnTraffic(conn gnet.Conn) gnet.Action {
	buf, _ := conn.Next(-1)
	conn.Write(buf)
	return gnet.None
}

func (s *EchoServer) OnClose(conn gnet.Conn, err error) (action gnet.Action) {
	log.Printf("conn %s close", conn.RemoteAddr())
	return
}

func (s *EchoServer) OnShutdown(engine gnet.Engine) {
	log.Println("echo server shutdown")
}

func StartGnet() {
	addr := conf.Config().Server.Gnet.Listen
	echo := &EchoServer{}
	log.Fatal(gnet.Run(echo, fmt.Sprintf("tcp://%s", addr)))
}

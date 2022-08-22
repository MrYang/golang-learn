package tcp

import "github.com/panjf2000/gnet/v2"

type EchoServer struct {
	*gnet.BuiltinEventEngine
}

func (s *EchoServer) OnOpen(conn gnet.Conn) (out []byte, action gnet.Action) {
	return
}

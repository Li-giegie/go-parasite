package main

import (
	"log"
	"net"
	"sync"
)

type Server struct {
	conn               *net.UDPConn
	cache              sync.Map
	conf
}

func NewServer(options ...Option) (*Server, error) {
	var srv = new(Server)
	option := margeOption(options,isServer)
	option.initOption(srv)
	return srv, nil
}

func (s *Server) ListenAndServer() (_err error) {
	s.conn, _err = net.ListenUDP(ProtocolUDP, s.localAddr)
	if _err != nil {
		return _err
	}
	defer s.conn.Close()
	for {
		buf, addr, err := Read(s.conn, s.byteBufferSize)
		if err != nil {
			log.Fatalln("server read data err:", err)
		}
		go s.Receive(buf, addr)
	}

}

func (s *Server) Receive(buf []byte, addr *net.UDPAddr) {
	_pack,ok := UnmarshalPack(buf)
	if !ok {
		log.Println("丢弃的包------", _pack.String())
		return
	}
	if _pack.Length == 1 && _pack.SN == 1 {
		log.Println("接收包1次完成")
		return
	}

}

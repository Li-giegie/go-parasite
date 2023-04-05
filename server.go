package main

import (
	"log"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	conn               *net.UDPConn
	localAddr          *net.UDPAddr
	bufferSize         int
	MaxReceivePackSize int
	cache              sync.Map
}

func NewServer(addr ...string) (*Server, error) {
	localAddr, lErr := net.ResolveUDPAddr(ProtocolUDP, DefaultServerAddress)
	if lErr != nil {
		return nil, lErr
	}

	return &Server{localAddr: localAddr, bufferSize: DefaultBufferSize, MaxReceivePackSize: DefaultBufferSize}, nil
}

func (s *Server) ListenAndServer() (_err error) {
	s.conn, _err = net.ListenUDP(ProtocolUDP, s.localAddr)
	if _err != nil {
		return _err
	}
	defer s.conn.Close()
	for {
		buf, addr, err := Read(s.conn, s.bufferSize)
		if err != nil {
			log.Fatalln("server read data err:", err)
		}
		go s.Receive(buf, addr)
	}

}

func (s *Server) Receive(buf []byte, addr *net.UDPAddr) {
	_pack := new(Pack)
	_pack.Unmarshal(buf)
	if !_pack.Cheek() {
		log.Println("丢弃的包------", _pack.String())
		return
	}
	key := addr.String() + strconv.Itoa(int(_pack.ID))
	add_PackCache(&s.cache, key, _pack)

}

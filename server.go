package main

import (
	"encoding/binary"
	"log"
	"net"
	"sync"
)

type Server struct {
	conn *net.UDPConn
	localAddr *net.UDPAddr
	bufferSize int
	data sync.Map
}

func NewServer(addr ...string) (*Server,error) {
	localAddr,lErr := net.ResolveUDPAddr(ProtocolUDP,DefaultServerAddress)
	if lErr != nil {
		return nil, lErr
	}

	return &Server{localAddr: localAddr,bufferSize: DefaultBufferSize},nil
}

func (s *Server) ListenAndServer() (_err error) {
	s.conn,_err = net.ListenUDP(ProtocolUDP,s.localAddr)
	if _err != nil {
		return _err
	}
	defer s.conn.Close()

	var buf []byte
	for {
		buf = make([]byte,s.bufferSize)
		_len,addr,err:=s.conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln("server read data err:",err)
		}
		go s.AssembleData(buf[:_len],addr)
	}

}
// 组装报文
func (s *Server) AssembleData(data []byte,addr *net.UDPAddr)  {
	dataAllLen := len(data)
	dataLen := binary.LittleEndian.Uint32(data[:4])
	if dataAllLen != int(4+dataLen) {

	}
}

// 处理报文
func (s *Server) Process(data []byte,addr *net.UDPAddr)  {

}
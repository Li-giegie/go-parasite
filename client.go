package main

import (
	"log"
	"net"

	"github.com/Li-giegie/errors"
)

type Client struct {
	conn       *net.UDPConn
	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr
}

func NewClient(remoteAddr string, localAddr ...string) (*Client, error) {
	if localAddr == nil || len(localAddr) < 1 {
		localAddr = []string{DefaultClientAddress}
	}
	lAddr, lErr := net.ResolveUDPAddr(ProtocolUDP, localAddr[0])
	rAddr, rErr := net.ResolveUDPAddr(ProtocolUDP, remoteAddr)
	if lErr != nil || rErr != nil {
		return nil, errors.NewErrors(lErr, rErr)
	}
	return &Client{localAddr: lAddr, remoteAddr: rAddr}, nil
}

func (c *Client) Connect() error {
	conn, err := net.DialUDP(ProtocolUDP, c.localAddr, c.remoteAddr)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.Read()
	return nil
}

func (c *Client) Read() {
	defer c.conn.Close()
	for {
		buf, addr, err := Read(c.conn, DefaultBufferSize)
		if err != nil {
			log.Fatalln(err)
		}
		go Merge(buf, addr)
	}
}

func (c *Client) Send(data []byte) error {
	var buf []byte
	packs, err := Disassembly(data)
	for _, v := range packs {

		buf, err = v.Marshal()
		if err != nil {
			log.Println("write 1", err)
			continue
		}
		if _, err = c.conn.Write(buf); err != nil {
			log.Println("write 2", err)
		}
	}

	return err
}

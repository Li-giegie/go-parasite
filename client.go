package main

import (
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/Li-giegie/errors"
)

type Client struct {
	conn       *net.UDPConn
	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr
	PassiveMessage sync.Map		//被动消息 （请求后的响应）
	activeMessage	chan *Pack	//主动消息
}

type Reply struct {
	ID uint32
	passiveMessage *sync.Map
	Wait chan int
}

func (r *Reply) Reply() *Pack {
	<- r.Wait
	r.passiveMessage.Load(r)
	return nil
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
	return &Client{localAddr: lAddr, remoteAddr: rAddr,activeMessage: make(chan *Pack)}, nil
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
		go c.Receive(buf, addr)
	}
}

func (c *Client) Send(data []byte) (Reply,error) {
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

	return Reply{ID: packs[0].ID,passiveMessage: &c.PassiveMessage,Wait: make(chan int)},err
}

func (c *Client) Receive(buf []byte, addr *net.UDPAddr) {
	_pack := new(Pack)
	_pack.Unmarshal(buf)
	if !_pack.Cheek() {
		log.Println("丢弃的包------", _pack.String())
		return
	}
	key := addr.String() + strconv.Itoa(int(_pack.ID))
	vals,ok := c.PassiveMessage.Load(key)

}

func (c *Client) GetPassiveMessage(Key string) ([]*Pack,bool) {
	v,ok := c.PassiveMessage.Load(Key)
	if !ok {
		return nil,false
	}

	packs,ok := v.([]*Pack)
	if !ok {
		return nil,false
	}
	return packs,true
}
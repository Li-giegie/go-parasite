package main

import (
	"log"
	"net"
	"sync"
)

type Reply struct {
	ID uint32
	passiveMessage *sync.Map
	cachePack cachePack
}
func (r *Reply) Reply() *Pack {
	r.cachePack.Init(r.ID)
	r.passiveMessage.Store(r.ID,&r.cachePack)
	sp := <- r.cachePack.wait

	r.passiveMessage.Store(r.ID,&r.cachePack)
	return sp
}

type cachePack struct {
	id uint32
	wait chan *Pack
	state uint8
	pack Packs
}

func (c *cachePack) Init (id uint32) {
	c.wait = make(chan *Pack)
	c.pack = make([]*Pack,0)
	c.id = id
}

type Client struct {
	conn       *net.UDPConn
	conf
	responsePack chan *Pack
	pushPack chan *Pack

	taskCache sync.Map
}

func NewClient(remoteAddr string,options ...Option) (*Client, error) {
	var c =new(Client)
	c.responsePack = make(chan *Pack)
	ops := margeOption(options,isClient)
	ops.initOption(&c.conf)
	return c,nil
}

func (c *Client) Connect() error {
	conn, err := net.DialUDP(ProtocolUDP, c.localAddr, c.remoteAddr)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.Receive()
	return nil
}

func (c *Client) Receive() {
	_Receive(c.conn,c.byteBufferSize,c.responsePack,&c.taskCache)
}

func (c *Client) Send(data []byte) (Reply,error) {
	var buf []byte
	packs, err := Disassembly(data)
	go func() {
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
	}()


	return Reply{ID: packs[0].ID},err
}

func (c *Client) SetLocalAddr(addr *net.UDPAddr)()  {
	c.localAddr = addr
}

func (c *Client) SetRemoteAddr(addr net.UDPAddr)()  {
	c.remoteAddr = &addr
}
func (c *Client) SetResponsePackReceiveSize(n uint)()  {
	c.responsePack = make(chan *Pack,n)
}
func (c *Client) SetByteBufferSize(n uint)()  {
	c.byteBufferSize = int(n)
}
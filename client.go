package main

import (
	"log"
	"net"
	"sync"
	"time"
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
	defer c.conn.Close()
	for {
		buf, _, err := Read(c.conn, c.byteBufferSize)
		if err != nil {
			log.Fatalln(err)
		}
		pack,ok := UnmarshalPack(buf)
		if !ok {
			log.Println("验证失败的消息")
			continue
		}
		if pack.Length == 1 && pack.SN == 1 {
			c.responsePack <- pack
			log.Println("接收包1次完成")
			continue
		}
		v,ok := c.taskCache.Load(pack.ID)
		if !ok {
			log.Println("不存在消息任务新建一个")
			tmpTask :=NewTask(pack)
			c.taskCache.Store(pack.ID,tmpTask)
			go tmpTask.run()
			continue
		}
		tmpTask := v.(*task)
		if !tmpTask.isClose {
			log.Println("消息任务：向管道添加消息")
			tmpTask.packChan <- pack
			continue
		}
		log.Println("向已经关闭的消息发送------",pack.ID)
	}
}

type task struct {
	id uint32
	packs Packs
	isClose bool
	packChan chan *Pack
	updateTime int64
}
func NewTask(pack *Pack) *task {
	return &task{
		id:         pack.ID,
		packs:      Packs{pack},
		updateTime: time.Now().UnixMilli(),
	}
}
func (t *task) run()  {
	var ti = time.NewTicker(time.Millisecond*300)
	for  {
		select {
		case pack := <- t.packChan:
			d := time.Now().UnixMilli()-t.updateTime
			if d <= 30 {
				d = 100
			}
			t.updateTime = d
			t.packs.Append(pack)
			if t.packs.CheckIntegrality() {
				t.isClose = true
				close(t.packChan)
				ti.Stop()
				return
			}
			ti.Reset(time.Millisecond * time.Duration(d))

		case <-ti.C:
			log.Println("超时任务：",t.id)
		}
	}
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
package main

import (
	"bytes"
	"sync"
)

var _PackAutoID uint32
var sock sync.Mutex
type Pack struct {
	buf   *bytes.Buffer
	ID    uint32       `json:"id"`
	State byte         `json:"state"`
	Buf   []byte       `json:"buf"`
}

func NewPack(data []byte) *Pack {
	sock.Lock()
	defer sock.Unlock()
	_PackAutoID++
	var buf = new(bytes.Buffer)
	buf.Write(data)
	return &Pack{buf: buf,ID: _PackAutoID}
}

func (p *Pack) Marshal() ([][]byte,error) {
	return nil,nil
}

func (p *Pack) Unmarshal()  {

}

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/Li-giegie/errors"
	"io"
	"log"
	"math"
	"sync"
)

var nullPackLength int

func init()  {
	countNullPackJsonLength()
}

var _PackAutoID uint32
var lock sync.Mutex

type Pack struct {
	buf  *bytes.Buffer
	ID uint32
	subPackBufSize int
}

func NewPack(data []byte,subPackBufSize int) (_pack *Pack) {
	lock.Lock()
	defer lock.Unlock()
	_PackAutoID++
	var buf = new(bytes.Buffer)
	buf.Write(data)
	_pack = &Pack{
		buf: buf,
		ID: _PackAutoID,
		subPackBufSize: subPackBufSize,
	}
	return
}

func (p *Pack) Marshal() ([]*SubPack,error) {
	var subPacks = make([]*SubPack,0)
	var err error
	bufLen := p.buf.Len()
	var subPackSize = p.subPackBufSize - 12
	if subPackSize < 1 {
		return nil, errors.NewErrors("subPackSize < 13 !")
	}
	if bufLen <= p.subPackBufSize - 12 {
		subPacks = append(subPacks, &SubPack{p.ID,1,0,p.buf.Bytes()})
		return subPacks,nil
	}

	var count = int(math.Ceil(float64(bufLen) / float64(p.subPackBufSize - 12)))
	subPacks = make([]*SubPack,count)
	for  i:=0;i<count;i++ {
		var tmpSP SubPack
		tmpSP.ID = p.ID
		tmpSP.SN = uint32(i+1)
		tmpSP.Length = uint32(count)
		if p.buf.Len() < subPackSize {
			subPackSize = p.buf.Len()
		}
		tmpSP.Msg = make([]byte,subPackSize)
		if _,err = p.buf.Read(tmpSP.Msg); err != nil && err != io.EOF {
			return nil, errors.NewErrors("pack Marshal err:",err)
		}
		subPacks[i] = &tmpSP
	}

	return subPacks,nil
}

func (p *Pack) Unmarshal(data chan []byte) (*SubPack,error) {

	return nil,nil
}

type SubPack struct {
	ID     uint32        `json:"id"`	//消息ID
	SN     uint32        `json:"sn"`	//消息序号
	Length uint32		`json:"length"`	//消息长度
	Msg    []byte        `json:"msg"`	//消息体
}

func (p *SubPack) Marshal() ([]byte,error)  {
	var buf =new(bytes.Buffer)
	idE:=binary.Write(buf,binary.LittleEndian,p.ID)
	snE:=binary.Write(buf,binary.LittleEndian,p.SN)
	lE:=binary.Write(buf,binary.LittleEndian,p.Length)
	_,err := buf.Write(p.Msg)
	if idE != nil || snE != nil || lE != nil || err != nil{
		return nil,errors.NewErrors("subPack marshal err: ",idE,snE,lE)
	}
	return buf.Bytes(),nil
}

func (p *SubPack) Unmarshal(data []byte)  {
	p.ID = binary.LittleEndian.Uint32(data[:4])
	p.SN = binary.LittleEndian.Uint32(data[4:8])
	p.Length = binary.LittleEndian.Uint32(data[8:12])
	p.Msg = data[12:]
}

func countNullPackJsonLength(){
	buf,err := json.Marshal(SubPack{})
	if err != nil {
		log.Fatalln("init pack json length err:",err)
	}
	nullPackLength = len(buf) - 2-1-1
	//log.Println("null pack length:",nullPackLength,"\npack:",string(buf))
}
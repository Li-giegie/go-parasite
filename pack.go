package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"sync"

	"github.com/Li-giegie/errors"
)

var _PackAutoID uint32
var _lock sync.Mutex
var PackBufSize int = DefaultBufferSize

type PackAndAddress struct {
	pack *Pack
	addr *net.UDPAddr
}

func (pack *PackAndAddress) getKey() string {
	return pack.addr.String() + strconv.Itoa(int(pack.pack.ID))
}

func getPackAutoID() uint32 {
	_lock.Lock()
	defer _lock.Unlock()
	_PackAutoID++
	return _PackAutoID
}

// 分解包成指定大小
func Disassembly(data []byte) (Packs []*Pack, err error) {
	var _Buf = new(bytes.Buffer)
	_Buf.Write(data)
	_ID := getPackAutoID()

	bufLen := _Buf.Len()
	var tmpPackSize = PackBufSize - 12
	if tmpPackSize < 1 {
		return nil, errors.NewErrors("PackSize < 13 !")
	}
	if bufLen <= PackBufSize-12 {
		Packs = append(Packs, NewPack(_ID, 1, 1, _Buf.Bytes()))
		return Packs, nil
	}

	var count = int(math.Ceil(float64(bufLen) / float64(PackBufSize-12)))
	Packs = make([]*Pack, count)
	for i := 0; i < count; i++ {
		if _Buf.Len() < tmpPackSize {
			tmpPackSize = _Buf.Len()
		}
		tmpMsg := make([]byte, tmpPackSize)
		if _, err = _Buf.Read(tmpMsg); err != nil && err != io.EOF {
			return nil, errors.NewErrors("pack Marshal err:", err)
		}
		Packs[i] = NewPack(_ID, uint32(i+1), uint32(count), tmpMsg)
	}

	return Packs, nil
}

type Pack struct {
	ID     uint32 `json:"id"`     //消息ID
	SN     uint32 `json:"sn"`     //消息序号
	Length uint32 `json:"length"` //消息长度
	Md5Sum uint16 `json:"md5Sum"` //消息指纹
	Msg    []byte `json:"msg"`    //消息体
}

func NewPack(id, sn, length uint32, msg []byte) *Pack {
	Pack := &Pack{
		ID:     id,
		SN:     sn,
		Length: length,
		Msg:    msg,
	}
	Pack.Md5Sum = Pack.SumMd5()
	return Pack
}

func (p *Pack) Marshal() ([]byte, error) {
	var buf = new(bytes.Buffer)
	idE := binary.Write(buf, binary.LittleEndian, p.ID)
	snE := binary.Write(buf, binary.LittleEndian, p.SN)
	lE := binary.Write(buf, binary.LittleEndian, p.Length)
	md5Err := binary.Write(buf, binary.LittleEndian, p.Md5Sum)
	_, msgErr := buf.Write(p.Msg)
	if idE != nil || snE != nil || lE != nil || md5Err != nil || msgErr != nil {
		return nil, errors.NewErrors("Pack marshal err: ", idE, snE, lE)
	}
	return buf.Bytes(), nil
}

func (p *Pack) Unmarshal(data []byte) {
	p.ID = binary.LittleEndian.Uint32(data[:4])
	p.SN = binary.LittleEndian.Uint32(data[4:8])
	p.Length = binary.LittleEndian.Uint32(data[8:12])
	p.Md5Sum = binary.LittleEndian.Uint16(data[12:14])
	p.Msg = data[14:]
}

func (p *Pack) SumMd5() uint16 {
	var buf = new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, p.ID)
	_ = binary.Write(buf, binary.LittleEndian, p.SN)
	_ = binary.Write(buf, binary.LittleEndian, p.Length)
	_, _ = buf.Write(p.Msg)
	return sumMd5(getMd5(buf.Bytes()))
}

func (p *Pack) Cheek() bool {
	return p.Md5Sum == p.SumMd5()
}

func (p *Pack) String() string {
	return fmt.Sprint("ID:", p.ID, ", SN:", p.SN, ", Length:", p.Length, ", Md5Sum:", p.Md5Sum, ", MSG:", string(p.Msg))
}

type Packs []*Pack

func (p Packs) Len() int {
	return len(p)
}

// 实现sort.Interface接口的比较元素方法
func (p Packs) Less(i, j int) bool {
	return p[i].SN < p[j].SN
}

// 实现sort.Interface接口的交换元素方法
func (p Packs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

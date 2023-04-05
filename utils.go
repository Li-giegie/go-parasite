package main

import (
	"crypto/md5"
	"net"
	"strconv"
	"sync"
)

func getMd5(src []byte) []byte {
	_md5 := md5.New()
	_, err := _md5.Write(src)
	if err != nil {
		return src[:16]
	}
	return _md5.Sum(nil)
}

func sumMd5(b []byte) (u16 uint16) {
	for _, v := range b {
		u16 += uint16(v)
	}
	return
}

// 计算数字转字符串长度
func countNumLenUint32(n uint32) int {
	return len(strconv.Itoa(int(n)))
}

func Sum(i int) int {
	if i == 1 {
		return 1
	}
	return i + Sum(i-1)
}

func Read(conn *net.UDPConn, BufferSize int) ([]byte, *net.UDPAddr, error) {

	var data = make([]byte, BufferSize)

	n, uaddr, err := conn.ReadFromUDP(data)
	if err != nil {
		return nil, nil, err
	}

	return data[:n], uaddr, nil
}

func add_PackCache(smap *sync.Map, key string, _val *Pack) bool {
	val, ok := smap.Load(key)
	if !ok {
		smap.Store(key, []*Pack{_val})
	}
	packs, ok := val.([]*Pack)
	if !ok {
		smap.Store(key, []*Pack{_val})
	}
	for _, v := range packs {
		if v.SN == _val.SN {
			return false
		}
	}
	return true
}

func get_PackCache(smap *sync.Map, key string) ([]*Pack, bool) {
	val, ok := smap.Load(key)
	if !ok {
		return nil, false
	}
	_pack, ok := val.([]*Pack)
	if !ok {
		return nil, false
	}

	return _pack, ok
}

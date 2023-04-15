package main

import "net"

const (
	localAddr uint8 = iota
	remoteAddr
	responsePackReceiveSize
	byteBufferSize
	isServer
	isClient
)

type conf struct {
	localAddr *net.UDPAddr
	remoteAddr *net.UDPAddr
	responsePackReceiveSize int
	byteBufferSize int
}

func (o *conf) SetLocalAddr(addr *net.UDPAddr) {
	o.localAddr = addr
}

func (o *conf) SetRemoteAddr(addr *net.UDPAddr) {
	o.remoteAddr = addr
}

func (o *conf) SetResponsePackReceiveSize(n int) {
	o.responsePackReceiveSize = n
}

func (o *conf) SetByteBufferSize(n int) {
	o.byteBufferSize = n
}

type OptionI interface {
	SetLocalAddr(addr *net.UDPAddr)
	SetRemoteAddr(addr *net.UDPAddr)
	SetResponsePackReceiveSize(n int)
	SetByteBufferSize(n int)
}

type Option map[uint8]interface{}

func (o Option) initOption(obj OptionI){
	for k, v := range o {
		switch k {
		case localAddr:
			obj.SetLocalAddr(v.(*net.UDPAddr))
		case remoteAddr:
			obj.SetRemoteAddr(v.(*net.UDPAddr))
		case responsePackReceiveSize:
			obj.SetResponsePackReceiveSize(v.(int))
		case byteBufferSize:
			obj.SetByteBufferSize(v.(int))
		}
	}

}

func SetLocalAddr(addr string) Option{
	_addr,err := net.ResolveUDPAddr("udp",addr)
	if err != nil {
		panic(any(err))
	}
	return Option{localAddr:_addr}
}

func SetRemoteAddr(addr string) Option{
	_addr,err := net.ResolveUDPAddr("udp",addr)
	if err != nil {
		panic(any(err))
	}
	return Option{remoteAddr:_addr}
}

func SetResponsePackReceiveSize(n uint) Option{
	return Option{responsePackReceiveSize:n}
}

func SetSinglePackBufferSize(n uint) Option{
	return Option{byteBufferSize:n}
}


func margeOption(opt []Option,who uint8) *Option {
	var newOption = Option{
		localAddr: DefaultClientAddress,
		remoteAddr: DefaultServerAddress,
		responsePackReceiveSize: DefaultResponsePackReceiveSize,
		byteBufferSize: DefaultByteBufferSize,
	}
	if who == isServer {
		newOption[localAddr] = DefaultServerAddress
	}
	for _, Option := range opt {
		for k, v := range Option {
			newOption[k] = v
		}
	}

	return &newOption
}


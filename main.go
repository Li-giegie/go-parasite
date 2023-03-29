package main

import (
	"fmt"
	"log"
	"net"
)

func main()  {


}


func Read(conn *net.UDPConn)  {
	var data []byte
	defer conn.Close()

	for  {
		data = make([]byte,1024)
		n,uaddr,err := conn.ReadFromUDP(data)
		if err != nil {
			log.Fatalln("读取消息错误：",err)
		}
		//fmt.Printf("server addr:%v data:%v \n",uaddr.String(),string(data[:n]))

		_,err = conn.WriteToUDP([]byte("md5" + getMd5(data[:n])),uaddr)
		if err != nil {
			log.Fatalln(fmt.Sprintf("server panic addr:%v err :写入回复消息错误： %v \n",uaddr.String(),err))
		}
	}
}

func Read2(conn *net.UDPConn) ([]byte,*net.UDPAddr,error) {

	var data = make([]byte,1024)

	n,uaddr,err := conn.ReadFromUDP(data)
	if err != nil {
		return nil, nil, err
	}

	return data[:n],uaddr,nil
}



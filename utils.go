package main

import (
	"crypto/md5"
	"log"
	"net"
	"sync"
)

func _getMd5(src []byte) []byte {
	_md5 := md5.New()
	_, err := _md5.Write(src)
	if err != nil {
		return src[:16]
	}
	return _md5.Sum(nil)
}

func _sumMd5(b []byte) (u16 uint16) {
	for _, v := range b {
		u16 += uint16(v)
	}
	return
}

func _sum(i uint32) uint32 {
	if i == 1 {
		return 1
	}
	return i + _sum(i-1)
}

func _read(conn *net.UDPConn, BufferSize int) ([]byte, *net.UDPAddr, error) {
	var data = make([]byte, BufferSize)
	n, uaddr, err := conn.ReadFromUDP(data)
	if err != nil {
		return nil, nil, err
	}
	return data[:n], uaddr, nil
}

func _Receive(conn *net.UDPConn,byteBufferSize int,responseChan chan *Pack,taskCache *sync.Map) {
	defer conn.Close()
	for {
		buf, _, err := _read(conn, byteBufferSize)
		if err != nil {
			log.Fatalln(err)
		}
		pack,ok := UnmarshalPack(buf)
		if !ok {
			log.Println("验证失败的消息")
			continue
		}
		if pack.Length == 1 && pack.SN == 1 {
			responseChan <- pack
			log.Println("接收包1次完成")
			continue
		}

		storeTask,ok := LoadTaskSyncMayCache(taskCache,pack.ID)
		if !ok {
			log.Println("不存在消息任务新建一个")
			_Task :=NewTask(pack,responseChan)
			StoreTaskSyncMayCache(taskCache,_Task)
			go _Task.run()
			continue
		}
		if storeTask.State == TaskEnd {
			log.Println("向已经关闭的消息发送------",pack.ID)
			continue
		}
		storeTask.RequestChan <- pack
		log.Println("消息任务：向管道添加消息")
	}
}

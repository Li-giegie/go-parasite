package main

import (
	"fmt"
	"log"
	"testing"
)

type Info struct {
	data  int64
	data2 float64
	data3 []byte
}


func TestNewPack(t *testing.T) {

	pack := NewPack([]byte("1234554329"),103)
	subPack,err := pack.Marshal()
	if err != nil {
		log.Fatalln(err)
	}

	for _, s := range subPack {
		fmt.Println("msg: ",*s,string(s.Msg))
	}
}

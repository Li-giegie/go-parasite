package main

import (
	"fmt"
	"log"
	"testing"
	"time"
)

type Info struct {
	data  int64
	data2 float64
	data3 []byte
}

func TestNewPack(t *testing.T) {

	var data = []byte("123321")

	t1 := time.Now()
	Pack, err := Disassembly(data)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("ok ", time.Now().Sub(t1))

	for _, s := range Pack {
		fmt.Println(s.String(), s.Cheek())
	}
}

func TestSrv(t *testing.T) {
	Srv()
}
func TestCli(t *testing.T) {
	Cli()
}

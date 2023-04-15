package main

import (
	"fmt"
	"log"
	"sort"
	"testing"
	"time"
)

type Info struct {
	data  int64
	data2 float64
	data3 []byte
}

func TestNewPack(t *testing.T) {


	var tt = time.NewTicker(time.Second)
	var ttt = time.NewTicker(time.Second*3)
	for {
		select {
		case <-ttt.C:
			log.Println("ttt----")
		case <-tt.C:
			log.Println("tt------")
			ttt.Reset(time.Second*3)
		}
	}
	return
	var pa =[]*Pack{&Pack{
		ID:     2,
		SN:     2,
		Length: 0,
		Md5Sum: 0,
		Msg:    nil,
	},
	&Pack{
		ID:     1,
		SN:     1,
		Length: 0,
		Md5Sum: 0,
		Msg:    nil,
	},
		&Pack{
			ID:     3,
			SN:     3,
			Length: 0,
			Md5Sum: 0,
			Msg:    nil,
		},
	}
	var p = Packs(pa)
	sort.Sort(p)

	for _, pack := range p {
		fmt.Println(*pack)
	}

	return
	var is []*int
	var aaaa = 1
	is = append(is,&aaaa )
	fmt.Println(is)
	return
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

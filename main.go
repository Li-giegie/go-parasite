package main

import (
	"fmt"
	"log"
)

func main() {
	// go Srv()

	// // go Cli()
	// fmt.Println("ListenAndServer------------")
	// select {}
}

func Srv() {
	srv, err := NewServer(DefaultServerAddress)
	if err != nil {
		log.Fatalln(err)
	}

	if err = srv.ListenAndServer(); err != nil {
		log.Fatalln(err)
	}
}

func Cli() {
	cli, err := NewClient(DefaultServerAddress,SetLocalAddr("127.0.0.1:9999"))
	if err != nil {
		log.Fatalln(err)
	}

	if err = cli.Connect(); err != nil {
		defer cli.conn.Close()
		log.Fatalln(err)
	}

	fmt.Println(cli.Send([]byte("1")))
	fmt.Println(cli.Send([]byte("12")))
	fmt.Println(cli.Send([]byte("123")))
	fmt.Println(cli.Send([]byte("1234")))
	fmt.Println(cli.Send([]byte("12345")))
}

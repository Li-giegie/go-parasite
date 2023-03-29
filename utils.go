package main

import (
	"crypto/md5"
	"fmt"
	"log"
)

func getMd5(src []byte) string {
	_md5 := md5.New()
	_,err := _md5.Write(src)
	if err != nil {
		log.Println("get md5 err：",err)
		return ""
	}
	return fmt.Sprintf("%x",_md5.Sum(nil))
}

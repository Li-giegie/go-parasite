package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"strconv"
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
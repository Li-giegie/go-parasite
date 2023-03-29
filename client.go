package main

import (
	"net"
	"github.com/Li-giegie/errors"
)

type Client struct {
	conn *net.UDPConn
	localAddr *net.UDPAddr
	remoteAddr *net.UDPAddr
}

func NewClient(remoteAddr string,localAddr ...string) (*Client,error) {
	if localAddr == nil || len(localAddr) < 1 {
		localAddr = []string{DefaultClientAddress}
	}
	lAddr,lErr := net.ResolveUDPAddr(ProtocolUDP,localAddr[0])
	rAddr,rErr := net.ResolveUDPAddr(ProtocolUDP,remoteAddr)
	if lErr != nil || rErr != nil {
		return nil,errors.NewErrors(lErr,rErr)
	}
	return &Client{localAddr: lAddr,remoteAddr: rAddr},nil
}

func (c *Client) Connect() error {
	conn,err := net.DialUDP(ProtocolUDP,c.localAddr,c.remoteAddr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}
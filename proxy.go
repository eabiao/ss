package main

import (
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"log"
	"net"
)

var (
	ss = initSSProxy()
)

type SSProxy struct {
	server string
	cipher *shadowsocks.Cipher
}

func initSSProxy() *SSProxy {
	return &SSProxy{
		server: config.Server,
		cipher: initCipher(config.Method, config.Passwd),
	}
}

func initCipher(method, passwd string) *shadowsocks.Cipher {
	cipher, err := shadowsocks.NewCipher(method, passwd)
	if err != nil {
		log.Fatal(err)
	}
	return cipher
}

func (p *SSProxy) connect(addr string) (net.Conn, error) {
	return shadowsocks.Dial(addr, p.server, p.cipher.Copy())
}

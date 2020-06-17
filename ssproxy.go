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

// 初始化ss代理
func initSSProxy() *SSProxy {
	return &SSProxy{
		server: config.Server,
		cipher: initCipher(config.Method, config.Passwd),
	}
}

// 初始化加密器
func initCipher(method, passwd string) *shadowsocks.Cipher {
	cipher, err := shadowsocks.NewCipher(method, passwd)
	if err != nil {
		log.Fatal(err)
	}
	return cipher
}

// ss代理连接
func (sp *SSProxy) connect(addr string) (net.Conn, error) {
	return shadowsocks.Dial(addr, sp.server, sp.cipher.Copy())
}

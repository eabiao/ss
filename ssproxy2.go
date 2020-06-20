package main

import (
	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
	"log"
	"net"
)

var (
	ss2 = initSSProxy2()
)

type SSProxy2 struct {
	server string
	cipher core.Cipher
}

// 初始化ss代理
func initSSProxy2() *SSProxy2 {
	return &SSProxy2{
		server: config.Server,
		cipher: initCipher2(config.Method, config.Passwd),
	}
}

// 初始化加密器
func initCipher2(method, passwd string) core.Cipher {
	cipher, err := core.PickCipher(method, nil, passwd)
	if err != nil {
		log.Fatal(err)
	}
	return cipher
}

// ss代理连接
func (sp *SSProxy2) connect(addr string) (net.Conn, error) {
	rc, err := net.Dial("tcp", sp.server)
	if err != nil {
		log.Println("failed to connect to server", sp.server, err)
		return nil, err
	}

	rc = sp.cipher.StreamConn(rc)

	addrData := socks.ParseAddr(addr)
	if _, err = rc.Write(addrData); err != nil {
		log.Println("failed to send target address", addr, err)
		return nil, err
	}

	return rc, nil
}

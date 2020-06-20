package main

import (
	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
	"log"
	"net"
)

var (
	ss = initSSProxy()
)

type SSProxy struct {
	server string
	cipher core.Cipher
}

// 初始化ss代理
func initSSProxy() *SSProxy {
	return &SSProxy{
		server: config.Server,
		cipher: initCipher(config.Method, config.Passwd),
	}
}

// 初始化加密器
func initCipher(method, passwd string) core.Cipher {
	cipher, err := core.PickCipher(method, nil, passwd)
	if err != nil {
		log.Fatal(err)
	}
	return cipher
}

// ss代理连接
func (sp *SSProxy) connect(addr string) (net.Conn, error) {
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

//type SSProxy struct {
//	server string
//	cipher *shadowsocks.Cipher
//}
//
//// 初始化ss代理
//func initSSProxy() *SSProxy {
//	return &SSProxy{
//		server: config.Server,
//		cipher: initCipher(config.Method, config.Passwd),
//	}
//}
//
//// 初始化加密器
//func initCipher(method, passwd string) *shadowsocks.Cipher {
//	cipher, err := shadowsocks.NewCipher(method, passwd)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return cipher
//}
//
//// ss代理连接
//func (sp *SSProxy) connect(addr string) (net.Conn, error) {
//	return shadowsocks.Dial(addr, sp.server, sp.cipher.Copy())
//}

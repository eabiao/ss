package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Request struct {
	conn    net.Conn
	addr    string
	isHttps bool
	data    []byte
}

// 处理请求
func handleConnect(conn net.Conn) {
	defer conn.Close()

	req, err := parseRequest(conn)
	if err != nil {
		return
	}

	if req.isHttps {
		fmt.Fprint(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
	}

	doProxyConnect(req)
}

// 代理连接
func doProxyConnect(req *Request) {
	log.Println(req.addr)

	remote, err := ss.connect(req.addr)
	if err != nil {
		log.Println(req.addr, err)
		return
	}

	defer remote.Close()

	if !req.isHttps {
		remote.Write(req.data)
	}

	go copyStream(req.conn, remote)
	copyStream(remote, req.conn)
}

// 流复制
func copyStream(src, dst net.Conn) {
	var buff = connBuff.Get()
	defer connBuff.Put(buff)

	for {
		readN, err := src.Read(buff[:])
		if err != nil {
			return
		}

		_, err = dst.Write(buff[0:readN])
		if err != nil {
			return
		}
	}
}

// 解析请求信息
func parseRequest(client net.Conn) (*Request, error) {

	var buff = httpBuff.Get()
	defer httpBuff.Put(buff)

	readN, err := client.Read(buff[:])
	if err != nil {
		return nil, err
	}
	data := buff[:readN]

	var addr string
	var isHttps bool

	for _, line := range strings.Split(string(data), "\r\n") {
		if strings.HasPrefix(line, "CONNECT") {
			isHttps = true
			continue
		}
		if strings.HasPrefix(line, "Host:") {
			addr = strings.Fields(line)[1]
			break
		}
	}

	if !strings.Contains(addr, ":") {
		if isHttps {
			addr = addr + ":443"
		} else {
			addr = addr + ":80"
		}
	}

	request := &Request{
		conn:    client,
		addr:    addr,
		isHttps: isHttps,
		data:    data,
	}
	return request, nil
}
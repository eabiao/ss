package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// http请求
type HttpRequest struct {
	conn    net.Conn
	addr    string
	isHttps bool
	data    []byte
	host    string
	port    int
}

// 处理请求
func handleConnect(conn net.Conn) {
	defer conn.Close()

	req, err := parseRequest(conn)
	if err != nil {
		return
	}
	log.Println(req.addr)

	if req.isHttps {
		fmt.Fprint(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
	}

	remote, err := connectRemote(req)
	if err != nil {
		return
	}

	defer remote.Close()

	if !req.isHttps {
		remote.Write(req.data)
	}

	go copyStream(req.conn, remote)
	copyStream(remote, req.conn)
}

// 连接远端
func connectRemote(req *HttpRequest) (net.Conn, error) {
	// 如果域名不合法或为IP地址，走直连
	if !strings.Contains(req.host, ".") || net.ParseIP(req.host) != nil {
		return net.DialTimeout("tcp", req.addr, 2*time.Second)
	}
	return ss.connect(req.addr)
}

// 流复制
func copyStream(src, dst net.Conn) {
	var buff = connBuff.Get()
	defer func() {
		connBuff.Put(buff)
		src.SetDeadline(time.Now())
		dst.SetDeadline(time.Now())
	}()

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
func parseRequest(client net.Conn) (*HttpRequest, error) {

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

	addrParts := strings.SplitN(addr, ":", 2)
	host := addrParts[0]
	port, _ := strconv.Atoi(addrParts[1])

	request := &HttpRequest{
		conn:    client,
		addr:    addr,
		isHttps: isHttps,
		data:    data,
		host:    host,
		port:    port,
	}
	return request, nil
}

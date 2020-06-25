package main

import (
	"fmt"
	"io"
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

	if req.isHttps {
		fmt.Fprint(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
	}

	target, err := net.DialTimeout("tcp", req.addr, 100*time.Millisecond)
	if err == nil && connectTarget(req, target) == nil {
		direct.addDirect(req.host)
		return
	}

	if !direct.isDirect(req.host) {
		target, err = ss.connect(req.addr)
		if err == nil && connectTarget(req, target) == nil {
			return
		}
	}

	if err != nil {
		log.Println("all connect fail, bad network", req.host)
	}
}

// 连接目标地址
func connectTarget(req *HttpRequest, target net.Conn) error {
	defer target.Close()

	if !req.isHttps {
		_, err := target.Write(req.data)
		if err != nil {
			return err
		}
	}

	_, respN := relay(req.conn, target)
	if respN == 0 {
		return &Error{"read response data error"}
	}

	return nil
}

// 数据传输
func relay(left, right net.Conn) (int64, int64) {
	ch := make(chan int64)

	go func() {
		reqN, _ := io.Copy(right, left)
		right.SetDeadline(time.Now())
		left.SetDeadline(time.Now())
		ch <- reqN
	}()

	respN, _ := io.Copy(left, right)
	right.SetDeadline(time.Now())
	left.SetDeadline(time.Now())
	reqN := <-ch

	return reqN, respN
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

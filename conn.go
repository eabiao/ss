package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Request struct {
	conn    net.Conn
	addr    string
	host    string
	isHttps bool
	data    []byte
}

func handleConnect(conn net.Conn) {
	defer conn.Close()

	req, err := parseRequest(conn)
	if err != nil {
		return
	}

	if req.isHttps {
		fmt.Fprint(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
	}

	if block.contains(req.addr) {
		doProxyConnect(req)
	} else {
		doDirectConnect(req)
	}
}

// 直连
func doDirectConnect(req *Request) {
	log.Println("dr", req.addr)

	remote, err := net.DialTimeout("tcp", req.addr, 2*time.Second)
	if err != nil {
		block.put(req.addr)
		log.Println("dr", req.addr, err)
		return
	}

	if !req.isHttps {
		remote.Write(req.data)
	}

	go copyClientToRemote(req.addr, req.conn, remote, true)
	copyRemoteToClient(req.addr, remote, req.conn, true)
}

// 代理
func doProxyConnect(req *Request) {
	log.Println("ss", req.addr)

	remote, err := ss.connect(req.addr)
	if err != nil {
		log.Println("ss", req.addr, err)
		return
	}

	if !req.isHttps {
		remote.Write(req.data)
	}

	go copyClientToRemote(req.addr, req.conn, remote, false)
	copyRemoteToClient(req.addr, remote, req.conn, false)
}

func copyClientToRemote(addr string, client, remote net.Conn, isDirect bool) {
	var buff [2048]byte

	for {
		readN, err := client.Read(buff[:])
		if err != nil {
			return
		}

		_, err = remote.Write(buff[0:readN])
		if err != nil {
			if isDirect {
				block.put(addr)
				log.Println("dr write remote", addr, err)
			} else {
				log.Println("ss write remote", addr, err)
			}
			return
		}
	}
}

func copyRemoteToClient(addr string, remote, client net.Conn, isDirect bool) {
	var buff [2048]byte

	for {
		readN, err := remote.Read(buff[:])
		if err != nil {
			opErr, ok := err.(*net.OpError)
			if ok {
				if isDirect {
					block.put(addr)
					log.Println("dr read remote", addr, opErr.Err)
				} else {
					log.Println("ss read remote", addr, opErr.Err)
				}
			}
			return
		}

		_, err = client.Write(buff[0:readN])
		if err != nil {
			return
		}
	}
}

// 解析请求信息
func parseRequest(client net.Conn) (*Request, error) {

	var buff [1024]byte
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
		host:    "",
		isHttps: isHttps,
		data:    data,
	}
	return request, nil
}

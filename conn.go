package main

import (
	"fmt"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"log"
	"net"
	"strings"
	"time"
)

type Request struct {
	conn    net.Conn
	addr    string
	domain  string
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

	if block.contains(req.domain) {
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
		block.put(req.domain)
		log.Println("dr", req.addr, err)
		return
	}

	defer remote.Close()

	if !req.isHttps {
		remote.Write(req.data)
	}

	go copyClientToRemote(req, remote, true)
	copyRemoteToClient(req, remote, true)
}

// 代理
func doProxyConnect(req *Request) {
	log.Println("ss", req.addr)

	remote, err := ss.connect(req.addr)
	if err != nil {
		log.Println("ss", req.addr, err)
		return
	}

	defer remote.Close()

	if !req.isHttps {
		remote.Write(req.data)
	}

	go copyClientToRemote(req, remote, false)
	copyRemoteToClient(req, remote, false)
}

func copyClientToRemote(req *Request, remote net.Conn, isDirect bool) {
	var buff [2048]byte

	for {
		readN, err := req.conn.Read(buff[:])
		if err != nil {
			return
		}

		_, err = remote.Write(buff[0:readN])
		if err != nil {
			if isDirect {
				block.put(req.domain)
				log.Println("dr write remote", req.addr, err)
			} else {
				log.Println("ss write remote", req.addr, err)
			}
			return
		}
	}
}

func copyRemoteToClient(req *Request, remote net.Conn, isDirect bool) {
	var buff [2048]byte

	for {
		readN, err := remote.Read(buff[:])
		if err != nil {
			opErr, ok := err.(*net.OpError)
			if ok {
				if isDirect {
					block.put(req.domain)
					log.Println("dr read remote", req.addr, opErr.Err)
				} else {
					log.Println("ss read remote", req.addr, opErr.Err)
				}
			}
			return
		}

		_, err = req.conn.Write(buff[0:readN])
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

	domain, _ := publicsuffix.Domain(addr)

	request := &Request{
		conn:    client,
		addr:    addr,
		domain:  domain,
		isHttps: isHttps,
		data:    data,
	}
	return request, nil
}

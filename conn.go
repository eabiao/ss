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

func handleConnect(client net.Conn) {
	defer client.Close()

	var buff [1024]byte
	readN, err := client.Read(buff[:])
	if err != nil {
		return
	}
	data := buff[:readN]

	addr, isHttps := parseRequest(string(data))
	if isHttps {
		fmt.Fprint(client, "HTTP/1.1 200 Connection Established\r\n\r\n")
	}

	if block.contains(addr) {
		doProxyConnect(client, addr, isHttps, data)
	} else {
		doDirectConnect(client, addr, isHttps, data)
	}
}

// 直连
func doDirectConnect(client net.Conn, addr string, isHttps bool, data []byte) {
	log.Println("dr", addr)

	remote, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		block.put(addr)
		log.Println("dr", addr, err)
		return
	}

	if !isHttps {
		remote.Write(data)
	}

	go copyClientToRemote(addr, client, remote, true)
	copyRemoteToClient(addr, remote, client, true)
}

// 代理
func doProxyConnect(client net.Conn, addr string, isHttps bool, data []byte) {
	log.Println("ss", addr)

	remote, err := ss.connect(addr)
	if err != nil {
		log.Println("ss", addr, err)
		return
	}

	if !isHttps {
		remote.Write(data)
	}

	go copyClientToRemote(addr, client, remote, false)
	copyRemoteToClient(addr, remote, client, false)
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
func parseRequest(text string) (string, bool) {

	var addr string
	var isHttps bool

	for _, line := range strings.Split(text, "\r\n") {
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

	return addr, isHttps
}

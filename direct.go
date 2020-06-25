package main

import (
	"bufio"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	direct      = initDirect()
	directMutex = sync.Mutex{}
)

// 直连
type Direct struct {
	hostMap map[string]bool
}

// 初始化
func initDirect() *Direct {
	d := &Direct{
		hostMap: make(map[string]bool),
	}
	d.loadDirect()
	return d
}

// 加载文件
func (d *Direct) loadDirect() {

	resp, err := http.Get("http://lvzhanbiao.cn:3000/sidu/autoproxy/getAll")
	if err != nil {
		return
	}

	defer resp.Body.Close()

	br := bufio.NewReader(resp.Body)
	for {
		lineData, _, err := br.ReadLine()
		if err != nil {
			break
		}

		host := strings.TrimSpace(string(lineData))
		if host == "" {
			continue
		}

		d.hostMap[host] = true
	}
}

// 判断是否为直连
func (d *Direct) isDirect(host string) bool {
	for direct := range d.hostMap {
		if strings.HasSuffix(host, direct) {
			return true
		}
	}
	return false
}

// 增加记录
func (d *Direct) addDirect(host string) {
	directMutex.Lock()
	defer directMutex.Unlock()

	if d.isDirect(host) {
		return
	}

	domain := getTopDomain(host)
	if domain == "" {
		return
	}

	d.hostMap[domain] = true
	d.putDirect(domain)
}

// 保存
func (d *Direct) putDirect(domain string) {
	log.Println("add", domain)
	http.Get("http://lvzhanbiao.cn:3000/sidu/autoproxy/putHost?host=" + domain)
}

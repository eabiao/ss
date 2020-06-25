package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

var (
	direct      = initDirect()
	directMutex = sync.Mutex{}
)

// 直连
type Direct struct {
	defaultMap map[string]bool
	recordMap  map[string]bool
}

// 初始化
func initDirect() *Direct {
	d := &Direct{
		defaultMap: make(map[string]bool),
		recordMap:  make(map[string]bool),
	}
	d.loadDirect()
	d.saveDirect()
	return d
}

// 加载文件
func (d *Direct) loadDirect() {
	file, err := os.Open("./direct.txt")
	if err != nil {
		return
	}
	defer file.Close()

	br := bufio.NewReader(file)
	for {
		lineData, _, err := br.ReadLine()
		if err != nil {
			break
		}

		host := strings.TrimSpace(string(lineData))
		if host == "" {
			continue
		}

		if strings.HasPrefix(host, "||") {
			host = strings.TrimPrefix(host, "||")
			d.defaultMap[host] = true
		} else if !d.isDefaultDirect(host) {
			d.recordMap[host] = true
		}
	}
}

// 判断是否为直连
func (d *Direct) isDirect(host string) bool {
	return d.isDefaultDirect(host) || d.recordMap[host]
}

// 判断是否为默认直连
func (d *Direct) isDefaultDirect(host string) bool {
	for direct := range d.defaultMap {
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

	log.Println("add", host)
	domain := getTopDomain(host)
	if domain != "" {
		d.defaultMap[domain] = true
	} else {
		d.recordMap[host] = true
	}

	d.saveDirect()
}

// 保存
func (d *Direct) saveDirect() {
	var defaultList []string
	for k := range d.defaultMap {
		defaultList = append(defaultList, "||"+k)
	}
	sort.Strings(defaultList)

	var recordList []string
	for k := range d.recordMap {
		recordList = append(recordList, k)
	}
	sort.Strings(recordList)

	directText := strings.Join(defaultList, "\n") + "\n" + strings.Join(recordList, "\n")
	ioutil.WriteFile("./direct.txt", []byte(directText), os.ModePerm)
}

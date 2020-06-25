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

type Direct struct {
	directMap map[string]bool
}

func initDirect() *Direct {
	d := &Direct{
		directMap: make(map[string]bool),
	}
	d.loadDirect()
	return d
}

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

		d.directMap[host] = true
	}
}

func (d *Direct) isDirect(host string) bool {
	for direct := range d.directMap {
		if strings.HasSuffix(host, direct) {
			return true
		}
	}
	return false
}

func (d *Direct) addDirect(host string) {
	directMutex.Lock()
	defer directMutex.Unlock()

	if d.directMap[host] {
		return
	}
	log.Println("add direct", host)
	d.directMap[host] = true

	var list []string
	for k := range d.directMap {
		list = append(list, k)
	}

	sort.Strings(list)
	ioutil.WriteFile("./direct.txt", []byte(strings.Join(list, "\n")), os.ModePerm)
}

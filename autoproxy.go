package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
)

var (
	autoProxy      = initAutoProxy()
	autoProxyMutex = sync.Mutex{}
)

type AutoProxy struct {
	head       string
	blockList  map[string]bool
	directList map[string]bool
}

func initAutoProxy() *AutoProxy {
	ap := &AutoProxy{}
	ap.loadAutoProxy()
	return ap
}

func (ap *AutoProxy) loadAutoProxy() {

	file, err := os.Open("./autoproxy.txt")
	if err != nil {
		return
	}
	defer file.Close()

	br := bufio.NewReader(file)

	headSb := strings.Builder{}
	headEnd := false
	for {
		lineData, _, err := br.ReadLine()
		if err != nil {
			break
		}

		lineText := string(lineData)

		if !headEnd {
			headSb.WriteString(lineText + "\n")
		}

		if strings.HasPrefix(lineText, "!---------------") {
			headEnd = true
			continue
		}

		if strings.HasPrefix(lineText, "@@||") {
			host := strings.TrimPrefix(lineText, "@@||")
			ap.blockList[host] = true
			continue
		}

		if strings.HasPrefix(lineText, "||") {
			host := strings.TrimPrefix(lineText, "||")
			ap.directList[host] = true
			continue
		}
	}
	ap.head = headSb.String()
}

func (ap *AutoProxy) isBlock(host string) bool {
	return ap.blockList[host]
}

func (ap *AutoProxy) isDirect(host string) bool {
	return ap.directList[host]
}

func (ap *AutoProxy) addDirect(host string) {
	autoProxyMutex.Lock()
	defer autoProxyMutex.Unlock()

	if ap.directList[host] {
		return
	}

	ap.directList[host] = true

	var list []string
	for k := range ap.directList {
		list = append(list, k)
	}

	sort.Strings(list)

	sb := strings.Builder{}
	sb.WriteString(ap.head)
	sb.WriteString(strings.Join(list, "\n"))

	ioutil.WriteFile("./autoproxy.txt", []byte(sb.String()), os.ModePerm)
}

package main

import "sync"

var (
	block      = initBlock()
	blockMutex = sync.Mutex{}
)

type Block struct {
	m map[string]bool
}

func initBlock() *Block {
	return &Block{make(map[string]bool)}
}

func (b *Block) put(domain string) {
	blockMutex.Lock()
	defer blockMutex.Unlock()

	b.m[domain] = true
}

func (b *Block) contains(domain string) bool {
	return b.m[domain]
}

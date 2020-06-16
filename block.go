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

func (b *Block) put(addr string) {
	blockMutex.Lock()
	defer blockMutex.Unlock()

	b.m[addr] = true
}

func (b *Block) contains(addr string) bool {
	return b.m[addr]
}

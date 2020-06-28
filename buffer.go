package main

var (
	httpBuff = newBuffPool(512, 1024)
)

// 缓冲区池
type BuffPool struct {
	size int
	list chan []byte
}

// 初始化缓冲区池
func newBuffPool(n, size int) *BuffPool {
	return &BuffPool{
		size: size,
		list: make(chan []byte, n),
	}
}

// 获取缓冲区
func (b *BuffPool) Get() (buff []byte) {
	select {
	case buff = <-b.list:
	default:
		buff = make([]byte, b.size)
	}
	return
}

// 归还缓冲区
func (b *BuffPool) Put(buff []byte) {
	if len(buff) != b.size {
		panic("invalid buff size")
	}

	select {
	case b.list <- buff:
	default:
	}
	return
}

package main

var (
	httpBuff = newBuffer(512, 1024)
	connBuff = newBuffer(512, 4096)
)

// 缓冲区池
type Buffer struct {
	size int
	list chan []byte
}

func newBuffer(n, size int) *Buffer {
	return &Buffer{
		size: size,
		list: make(chan []byte, n),
	}
}

// 获取缓冲区
func (b *Buffer) Get() (buff []byte) {
	select {
	case buff = <-b.list:
	default:
		buff = make([]byte, b.size)
	}
	return
}

// 归还缓冲区
func (b *Buffer) Put(buff []byte) {
	if len(buff) != b.size {
		panic("invalid buff size")
	}

	select {
	case b.list <- buff:
	default:
	}
	return
}

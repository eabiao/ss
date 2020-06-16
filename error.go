package main

// 运行时异常
type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

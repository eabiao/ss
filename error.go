package main

type Error struct {
	msg string
}

func (err *Error) Error() string {
	return err.msg
}

package jtimer

import "errors"

var (
	ErrorUpdateTimer = errors.New("timer key is not exist")
	ErrorUpdateHeap  = errors.New("index out of range")
)

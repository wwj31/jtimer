package jtimer

import "errors"

var (
	ErrorAddTimer    = errors.New("add timer failed,cause by key is reapeated")
	ErrorUpdateTimer = errors.New("timer key is not exist")
	ErrorUpdateHeap  = errors.New("index out of range")
)

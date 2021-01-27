package jtimer

import (
	"fmt"
	"testing"
)

type array struct {
	val   int
	score int64
}

func (this *array) Priority() int64 {
	return this.score
}

func TestPropertyqueue(t *testing.T) {
	heap := NewQueue(nil, MIN_HEAP) // 创建一个小顶堆
	heap.Push(&array{val: 1, score: 154})
	fmt.Printf("1最低分:%v\n", heap.Peek().Priority())

	heap.Push(&array{val: 2, score: 410})
	fmt.Printf("2最低分:%v\n", heap.Peek().Priority())

	heap.Push(&array{val: 3, score: 360})
	fmt.Printf("3最低分:%v\n", heap.Peek().Priority())

	heap.Push(&array{val: 4, score: 90})
	fmt.Printf("4最低分:%v\n", heap.Peek().Priority())

	heap.Push(&array{val: 5, score: 103})
	fmt.Printf("5最低分:%v\n", heap.Peek().Priority())

	heap.Pop()
	fmt.Printf("pop最低分:%v\n", heap.Peek().Priority())

	heap.Change(MAX_HEAP) // 改变成大顶堆
	fmt.Printf("change最高分:%v\n", heap.Peek().Priority())

	for _, v := range heap.All() {
		fmt.Printf("[val:%v score:%v] ", v.(*array).val, v.(*array).score)
	}
	fmt.Printf("\n")

	// 用大顶堆排序是顺序
	for _, v := range heap.Sort() {
		fmt.Printf("[val:%v score:%v] ", v.(*array).val, v.(*array).score)
	}
	fmt.Printf("\n")
}

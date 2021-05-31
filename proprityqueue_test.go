package jtimer

import (
	"fmt"
	"math/rand"
	"testing"
)

type array struct {
	val   int
	score int64
}

func (s *array) Priority() int64 {
	return s.score
}

func TestPropertyqueue(t *testing.T) {
	heap := NewQueue(nil, MIN_HEAP, BINARY) // 创建一个小顶堆
	heap.Push(&array{val: 1, score: 154})
	fmt.Printf("1最低分:%v\n", heap.Peek().Priority())

	heap.Push(&array{val: 2, score: 410})
	fmt.Printf("2最低分:%v\n", heap.Peek().Priority())

	heap.Push(&array{val: 3, score: 360})
	fmt.Printf("3最低分:%v\n", heap.Peek().Priority())

	heap.Push(&array{val: 5, score: 103})
	fmt.Printf("5最低分:%v\n", heap.Peek().Priority())

	heap.Push(&array{val: 4, score: 90})
	fmt.Printf("4最低分:%v\n", heap.Peek().Priority())
	heap.Push(&array{val: 4, score: 80})
	fmt.Printf("4最低分:%v\n", heap.Peek().Priority())
	heap.Push(&array{val: 4, score: 60})
	fmt.Printf("4最低分:%v\n", heap.Peek().Priority())
	heap.Push(&array{val: 4, score: 100})
	fmt.Printf("4最低分:%v\n", heap.Peek().Priority())

	heap.Pop()
	fmt.Printf("pop最低分:%v\n", heap.Peek().Priority())
	heap.Pop()
	fmt.Printf("pop最低分:%v\n", heap.Peek().Priority())
	//heap.Pop()
	//fmt.Printf("pop最低分:%v\n", heap.Peek().Priority())
	//heap.Pop()
	//fmt.Printf("pop最低分:%v\n", heap.Peek().Priority())

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
	heap.Change(MIN_HEAP)
	for _, v := range heap.Sort() {
		fmt.Printf("[val:%v score:%v] ", v.(*array).val, v.(*array).score)
	}
	fmt.Printf("\n")
}

var arr = make([]int64, 0, 300000)

func init() {
	for i := 0; i < 300000; i++ {
		arr = append(arr, int64(rand.Intn(5000)+1))
	}
}

func BenchmarkHeap2(b *testing.B) {
	heap2 := NewQueue(nil, MIN_HEAP, BINARY)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range arr {
			heap2.Push(&array{val: 4, score: v})
		}
		//for heap2.Size() > 0 {
		//	heap2.Pop()
		//}
	}
}

func BenchmarkHeap4(b *testing.B) {
	heap2 := NewQueue(nil, MIN_HEAP, QUAD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range arr {
			heap2.Push(&array{val: 4, score: v})
		}
		//for heap2.Size() > 0 {
		//	heap2.Pop()
		//}
	}
}

package jtimer

type (
	IPriorityInterface interface {
		Priority() int64
	}
	Heap struct {
		arr []IPriorityInterface
		t   HEAPTYPE
		c   int
	}
)
type HEAPTYPE int

const (
	BINARY = 2	// 二叉堆
	QUAD   = 4	// 四叉堆
)

const (
	MIN_HEAP HEAPTYPE = 1 // 小顶堆
	MAX_HEAP HEAPTYPE = 2 // 大顶堆
)

func NewQueue(a []IPriorityInterface, t HEAPTYPE, c int) *Heap {
	h := &Heap{arr: a, t: t, c: c}
	if h.arr == nil {
		h.arr = make([]IPriorityInterface, 0, 10)
	} else {
		h.heapify()
	}
	return h
}

// 尾部插入一个元素
func (s *Heap) Push(val IPriorityInterface) {
	s.arr = append(s.arr, val)
	_up(s.arr, s.t, len(s.arr), s.c)
}

// 查看头部元素
func (s *Heap) Peek() IPriorityInterface {
	if len(s.arr) == 0 {
		return nil
	}
	return s.arr[0]
}

// 头部弹出一个元素
func (s *Heap) Pop() IPriorityInterface {
	n := len(s.arr)
	if n == 0 {
		return nil
	}
	s.arr[n-1], s.arr[0] = s.arr[0], s.arr[n-1]
	ret := s.arr[n-1]
	s.arr = s.arr[:n-1]
	_down(s.arr, s.t, 1, s.c)
	return ret
}

// 返回数量
func (s *Heap) Size() int {
	return len(s.arr)
}

// 获取数组
func (s *Heap) All() []IPriorityInterface {
	return s.arr
}

// 堆排序 小顶:大->小  大顶:小->大
func (s *Heap) Sort() []IPriorityInterface {
	for i := len(s.arr) - 1; i >= 0; i-- {
		s.arr[0], s.arr[i] = s.arr[i], s.arr[0]
		_down(s.arr[:i], s.t, 1, s.c)
	}
	return s.arr
}

//改变堆的性质
//t :1.小顶堆、2.大顶堆
func (s *Heap) Change(t HEAPTYPE) {
	if s.t != t {
		s.t = t
		s.heapify()
	}
}

func (s *Heap) heapify() {
	// 默认使用自底向上方式
	_buildHeap_bottom2top(s.arr, s.t, s.c)
	//_buildHeap_top2down(s.arr, s.t)
}

// 自底向上法 建堆
func _buildHeap_bottom2top(arr []IPriorityInterface, t HEAPTYPE, c int) {
	// 整体上自底向上， 每个节点自顶向下
	for i := len(arr) / 2; i > 0; i-- {
		_down(arr, t, i, c)
	}
}

// 自顶向下法 建堆
func _buildHeap_top2down(arr []IPriorityInterface, t HEAPTYPE, c int) {
	// 整体上自顶向下， 每个节点自底向上
	for i := 1; i <= len(arr); i++ {
		_up(arr, t, i, c)
	}
}

/*
向下渗透
arr: 使用的数组
t:   1.小顶堆  2.大顶堆
n:   当前数据的位置
*/
func _down(arr []IPriorityInterface, t HEAPTYPE, n int, x int) {
	_lr := func(n, x int) (int, int) {
		return (n * x) - (x - 2) - 1, (n * x)
	}

	l, r := _lr(n, x)
	for l < len(arr) {
		var _mum = l
		for i := l; i < r; i++ {
			if i+1 >= len(arr) {
				break
			}
			if t == MIN_HEAP {
				_mum = _min(arr, _mum, i+1)
			} else {
				_mum = _max(arr, _mum, i+1)
			}
		}
		if t == MIN_HEAP && arr[n-1].Priority() > arr[_mum].Priority() ||
			t == MAX_HEAP && arr[n-1].Priority() < arr[_mum].Priority() {
			arr[n-1], arr[_mum] = arr[_mum], arr[n-1]
			n = _mum + 1
			l, r = _lr(n, x)
		} else {
			return
		}
	}
}

/*
向上排查
arr: 使用的数组
t:   1.小顶堆  2.大顶堆
n:   当前数据的位置，1 表示根 n-1 表示数组下标
*/
func _up(arr []IPriorityInterface, t HEAPTYPE, n int, x int) {
	if n == 1 {
		return
	}

	var p int
	for n > 1 {
		p = (n - 2) / x
		if (t == MIN_HEAP && arr[p].Priority() > arr[n-1].Priority()) ||
			(t == MAX_HEAP && arr[p].Priority() < arr[n-1].Priority()) {
			arr[p], arr[n-1] = arr[n-1], arr[p]
			n = p + 1
		} else {
			return
		}
	}
}

func _min(arr []IPriorityInterface, i1, i2 int) int {
	if arr[i1].Priority() < arr[i2].Priority() {
		return i1
	}
	return i2
}

func _max(arr []IPriorityInterface, i1, i2 int) int {
	if arr[i1].Priority() > arr[i2].Priority() {
		return i1
	}
	return i2
}

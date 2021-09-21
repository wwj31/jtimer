package jtimer

const (
	BINARY = 2 // 二叉堆
	QUAD   = 4 // 四叉堆
	OCT    = 8 // 八叉堆
)
const (
	MAX_HEAP = 0 // 大顶堆
	MIN_HEAP = 1 // 小顶堆
)

type (
	IPriorityInterface interface {
		Priority() int64
		SetIndex(i int)
		GetIndex() int
	}
	Heap struct {
		arr []IPriorityInterface
		opt int
	}
)

func NewQueue(a []IPriorityInterface, htype, atype int) *Heap {
	h := &Heap{arr: a}
	if atype == 0 {
		atype = BINARY
	}
	h.stype(htype, atype)

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
	val.SetIndex(len(s.arr) - 1)

	heaptype, arytype := s.htype()
	_up(s.arr, len(s.arr), heaptype, arytype)
}

// 更新
func (s *Heap) Update(newVal IPriorityInterface) error {
	idx := newVal.GetIndex()
	if 0 > idx || idx >= len(s.arr) {
		return ErrorUpdateHeap
	}

	oldVal := s.arr[idx]
	if oldVal.Priority() == newVal.Priority() {
		return nil
	}
	s.arr[idx] = newVal
	heaptype, arytype := s.htype()
	_down(s.arr, idx+1, heaptype, arytype)
	_up(s.arr, idx+1, heaptype, arytype)
	return nil
}

// 重新维护堆顶，修改了堆顶元素后，需调用此函数
func (s *Heap) Topdown() {
	n := len(s.arr)
	if n == 0 {
		return
	}
	heaptype, arytype := s.htype()
	_down(s.arr, 1, heaptype, arytype)
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
	heaptype, arytype := s.htype()
	_down(s.arr, 1, heaptype, arytype)
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
	heaptype, arytype := s.htype()
	for i := len(s.arr) - 1; i >= 0; i-- {
		s.arr[0], s.arr[i] = s.arr[i], s.arr[0]
		_down(s.arr[:i], 1, heaptype, arytype)
	}
	return s.arr
}

//改变堆的性质
//t :1.小顶堆、2.大顶堆
func (s *Heap) Change(t int) {
	heaptype, arytype := s.htype()
	if heaptype != t {
		s.stype(t, arytype)
		s.heapify()
	}
}

func (s *Heap) heapify() {
	// 默认使用自底向上方式
	heaptype, arytype := s.htype()
	_buildHeap_bottom2top(s.arr, heaptype, arytype)
	//_buildHeap_top2down(s.arr, s.t)
}

// 设置堆顶类型，子节点数量
func (s *Heap) stype(h int, a int) {
	s.opt = (h << 4) | a
}

// 获取堆顶类型，子节点数量
func (s *Heap) htype() (t int, c int) {
	return s.opt >> 4, s.opt & 0xF
}

// 自底向上法 建堆
func _buildHeap_bottom2top(arr []IPriorityInterface, t int, c int) {
	// 整体上自底向上， 每个节点自顶向下
	for i := len(arr) / 2; i > 0; i-- {
		_down(arr, i, t, c)
	}
}

// 自顶向下法 建堆
func _buildHeap_top2down(arr []IPriorityInterface, t int, c int) {
	// 整体上自顶向下， 每个节点自底向上
	for i := 1; i <= len(arr); i++ {
		_up(arr, i, t, c)
	}
}

/*
向下渗透
arr: 使用的数组
n:   当前数据的位置
t:   1.小顶堆  2.大顶堆
c:   2\4\8叉堆
*/
func _down(arr []IPriorityInterface, n int, t int, c int) {
	_lr := func(n, x int) (int, int) {
		return (n * x) - (x - 2) - 1, (n * x)
	}

	l, r := _lr(n, c)
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

			arr[n-1].SetIndex(n - 1)
			arr[_mum].SetIndex(_mum)

			n = _mum + 1
			l, r = _lr(n, c)
		} else {
			return
		}
	}
}

/*
向上排查
arr: 使用的数组
n:   当前数据的位置，1 表示根 n-1 表示数组下标
t:   大\小顶堆
c:   2/4/8叉树
*/
func _up(arr []IPriorityInterface, n int, t int, c int) {
	if n == 1 {
		return
	}

	var p int
	for n > 1 {
		p = (n - 2) / c
		if (t == MIN_HEAP && arr[p].Priority() > arr[n-1].Priority()) ||
			(t == MAX_HEAP && arr[p].Priority() < arr[n-1].Priority()) {
			arr[p], arr[n-1] = arr[n-1], arr[p]

			arr[p].SetIndex(p)
			arr[n-1].SetIndex(n - 1)

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

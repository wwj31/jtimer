package jtimer

/*
	leftchild = n*2
	rightchild = n*2+1
	parent = n/2
*/
type (
	IPriorityInterface interface {
		Priority() int64
	}
	Heap struct {
		arr []IPriorityInterface
		t   HEAPTYPE
	}
)
type HEAPTYPE int

const (
	MIN_HEAP HEAPTYPE = 1 // 小顶堆
	MAX_HEAP HEAPTYPE = 2 // 大顶堆
)

func NewQueue(a []IPriorityInterface, t HEAPTYPE) *Heap {
	h := &Heap{arr: a, t: t}
	if h.arr == nil {
		h.arr = make([]IPriorityInterface, 0, 10)
	} else {
		h.heapify()
	}
	return h
}

// 尾部插入一个元素
func (this *Heap) Push(val IPriorityInterface) {
	this.arr = append(this.arr, val)
	_up(this.arr, this.t, len(this.arr))
}

// 查看头部元素
func (this *Heap) Peek() IPriorityInterface {
	if len(this.arr) == 0 {
		return nil
	}
	return this.arr[0]
}

// 头部弹出一个元素
func (this *Heap) Pop() IPriorityInterface {
	n := len(this.arr)
	if n == 0 {
		return nil
	}
	this.arr[n-1], this.arr[0] = this.arr[0], this.arr[n-1]
	ret := this.arr[n-1]
	this.arr = this.arr[:n-1]
	_down(this.arr, this.t, 1)
	return ret
}

// 返回数量
func (this *Heap) Size() int {
	return len(this.arr)
}

// 获取数组
func (this *Heap) All() []IPriorityInterface {
	return this.arr
}

// 堆排序 小顶:大->小  大顶:小->大
func (this *Heap) Sort() []IPriorityInterface {
	for i := len(this.arr) - 1; i >= 0; i-- {
		this.arr[0], this.arr[i] = this.arr[i], this.arr[0]
		_down(this.arr[:i], this.t, 1)
	}
	return this.arr
}

//改变堆的性质
//t :1.小顶堆、2.大顶堆
func (this *Heap) Change(t HEAPTYPE) {
	if this.t != t {
		this.t = t
		this.heapify()
	}
}

func (this *Heap) heapify() {
	// 默认使用自底向上方式
	_buildHeap_bottom2top(this.arr, this.t)
	//_buildHeap_top2down(this.arr, this.t)
}

// 自底向上法 建堆
func _buildHeap_bottom2top(arr []IPriorityInterface, t HEAPTYPE) {
	// 整体上自底向上， 每个节点自顶向下
	for i := len(arr) / 2; i > 0; i-- {
		_down(arr, t, i)
	}
}

// 自顶向下法 建堆
func _buildHeap_top2down(arr []IPriorityInterface, t HEAPTYPE) {
	// 整体上自顶向下， 每个节点自底向上
	for i := 1; i <= len(arr); i++ {
		_up(arr, t, i)
	}
}

/*
向下渗透
arr: 使用的数组
t:   1.小顶堆  2.大顶堆
n:   当前数据的位置，1 表示根 n-1 表示数组下标
*/
func _down(arr []IPriorityInterface, t HEAPTYPE, n int) {
	l := n * 2
	r := l + 1
	var s int

	lenarr := len(arr)
	if l > lenarr {
		return
	}

	if r > len(arr) {
		if (t == MIN_HEAP && arr[n-1].Priority() > arr[l-1].Priority()) ||
			(t == MAX_HEAP && arr[n-1].Priority() < arr[l-1].Priority()) {
			arr[n-1], arr[l-1] = arr[l-1], arr[n-1]
			_down(arr, t, l)
		}
	} else {
		if t == MIN_HEAP {
			s = _min(arr, l-1, r-1) + 1
			if arr[n-1].Priority() > arr[s-1].Priority() {
				arr[n-1], arr[s-1] = arr[s-1], arr[n-1]
				_down(arr, t, s)
			}
		} else if t == MAX_HEAP {
			s = _max(arr, l-1, r-1) + 1
			if arr[n-1].Priority() < arr[s-1].Priority() {
				arr[n-1], arr[s-1] = arr[s-1], arr[n-1]
				_down(arr, t, s)
			}
		}
	}
}

/*
向上排查
arr: 使用的数组
t:   1.小顶堆  2.大顶堆
n:   当前数据的位置，1 表示根 n-1 表示数组下标
*/
func _up(arr []IPriorityInterface, t HEAPTYPE, n int) {
	if n == 1 {
		return
	}
	p := n / 2

	if (t == MIN_HEAP && arr[p-1].Priority() > arr[n-1].Priority()) ||
		(t == MAX_HEAP && arr[p-1].Priority() < arr[n-1].Priority()) {
		arr[p-1], arr[n-1] = arr[n-1], arr[p-1]
		_up(arr, t, p)
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

package bp

import (
	"bytes"
	"sync"
)

var _pool = NewPool(256, 256*1024)

// GetPool Getting the default memory pool
func GetPool() *Pool { return _pool }

type Pool struct {
	begin      int
	pools      []*sync.Pool
	size2index map[int]int
}

// NewPool Creating a memory pool
// Left, right indicate the interval range of the memory pool, they will be transformed into pow(2,n)。
// Below left, Get method will return at least left bytes; above right, Put method will not reclaim the buffer.
func NewPool(left, right uint32) *Pool {
	var begin, end = binaryCeil(int(left)), binaryCeil(int(right))
	var p = &Pool{begin: begin, size2index: map[int]int{}}
	for i, j := begin, 0; i <= end; i *= 2 {
		capacity := i
		pool := &sync.Pool{New: func() any { return bytes.NewBuffer(make([]byte, 0, capacity)) }}
		p.pools = append(p.pools, pool)
		p.size2index[i] = j
		j++
	}
	return p
}

// Put Return buffer to memory pool
func (p *Pool) Put(b *bytes.Buffer) {
	if b != nil {
		if index, ok := p.size2index[b.Cap()]; ok {
			p.pools[index].Put(b)
		}
	}
}

// Get Fetch a buffer from the memory pool, of at least n bytes
func (p *Pool) Get(n int) *bytes.Buffer {
	var size = maxInt(binaryCeil(n), p.begin)
	var index, ok = p.size2index[size]
	if !ok {
		return bytes.NewBuffer(make([]byte, 0, n))
	}
	b := p.pools[index].Get().(*bytes.Buffer)
	if b.Cap() < size {
		b.Grow(size)
	}
	b.Reset()
	return b
}

func binaryCeil(x int) int {
	v := uint32(x)
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return int(v)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

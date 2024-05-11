package bp

import (
	"bytes"
	"sync"
)

var __pool = NewPool(256, 256*1024)

// Get fetch a buffer from the default pool, of at least n bytes
func Get(n int) *bytes.Buffer { return __pool.Get(n) }

// Put adds b to the default pool.
func Put(b *bytes.Buffer) { __pool.Put(b) }

type Pool struct {
	begin, end int
	shards     map[int]*sync.Pool
}

// NewPool Creating a buffer pool
// Left, right indicate the interval range of the buffer pool, they will be transformed into pow(2,n)。
// Below left, Get method will return at least left bytes; above right, Put method will not reclaim the buffer.
func NewPool(left, right uint32) *Pool {
	var begin, end = int(binaryCeil(left)), int(binaryCeil(right))
	var p = &Pool{
		begin: begin, end: end,
		shards: map[int]*sync.Pool{},
	}
	for i := begin; i <= end; i *= 2 {
		capacity := i
		p.shards[i] = &sync.Pool{
			New: func() any { return bytes.NewBuffer(make([]byte, 0, capacity)) },
		}
	}
	return p
}

// Put adds b to the pool.
func (p *Pool) Put(b *bytes.Buffer) {
	if b != nil {
		if pool, ok := p.shards[b.Cap()]; ok {
			pool.Put(b)
		}
	}
}

// Get fetch a buffer from the buffer pool, of at least n bytes
func (p *Pool) Get(n int) *bytes.Buffer {
	var size = maxInt(int(binaryCeil(uint32(n))), p.begin)
	if pool, ok := p.shards[size]; ok {
		b := pool.Get().(*bytes.Buffer)
		if b.Cap() < size {
			b.Grow(size)
		}
		b.Reset()
		return b
	}
	return bytes.NewBuffer(make([]byte, 0, n))
}

func binaryCeil(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

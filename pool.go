package bp

import (
	"bytes"
	"sync"
)

var _pool = NewPool(256, 256*1024)

func GetPool() *Pool { return _pool }

type Pool struct {
	pools      []*sync.Pool
	limits     []int
	size2index map[uint32]uint32
}

func NewPool(left, right uint32) *Pool {
	var p = &Pool{size2index: map[uint32]uint32{}}
	var begin, end = binaryCeil(left), binaryCeil(right)
	var index = uint32(0)
	for i := begin; i <= end; i *= 2 {
		capacity := i
		pool := &sync.Pool{New: func() any { return bytes.NewBuffer(make([]byte, 0, capacity)) }}
		p.pools = append(p.pools, pool)
		p.limits = append(p.limits, int(i))
		p.size2index[i] = index
		index++
	}
	return p
}

func (p *Pool) Put(b *bytes.Buffer) {
	if b == nil || b.Cap() == 0 {
		return
	}
	size := binaryCeil(uint32(b.Cap()))
	if index, ok := p.size2index[size]; ok {
		p.pools[index].Put(b)
	}
}

func (p *Pool) Get(n int) *bytes.Buffer {
	var size = max(int(binaryCeil(uint32(n))), p.limits[0])
	var index, ok = p.size2index[uint32(size)]
	if !ok {
		return bytes.NewBuffer(make([]byte, 0, n))
	}

	b := p.pools[index].Get().(*bytes.Buffer)
	b.Reset()
	if b.Cap() < size {
		b.Grow(size)
	}
	b.Reset()
	return b
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

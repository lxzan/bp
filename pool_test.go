package bp

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPut(t *testing.T) {
	Put(bytes.NewBuffer(make([]byte, 100)))
	assert.Equal(t, Get(24).Cap(), 256)
}

func TestBinaryCeil(t *testing.T) {
	assert.Equal(t, int(binaryCeil(3)), 4)
	assert.Equal(t, int(binaryCeil(125)), 128)
	assert.Equal(t, int(binaryCeil(1000)), 1024)
	assert.Equal(t, int(binaryCeil(65000)), 65536)
}

func TestPool_Get(t *testing.T) {
	var pool = NewPool(256, 1024*128)
	assert.Equal(t, pool.Get(100).Cap(), 256)
	assert.Equal(t, pool.Get(500).Cap(), 512)
	assert.Equal(t, pool.Get(1000).Cap(), 1024)
	assert.Equal(t, pool.Get(8000).Cap(), 8192)
	assert.Equal(t, pool.Get(362144).Cap(), 362144)
}

func TestPool_Put(t *testing.T) {
	var pool = NewPool(128, 1024*128)
	var b1 = pool.Get(200)
	var b2 = pool.Get(500)
	var b3 = pool.Get(600)
	pool.Put(b1)
	pool.Put(b2)
	pool.Put(b3)
	pool.Put(nil)
	pool.Put(new(bytes.Buffer))
	assert.Equal(t, pool.Get(501).Cap(), 512)

	pool.Put(bytes.NewBuffer(make([]byte, 800)))
	pool.Put(bytes.NewBuffer(make([]byte, 1000)))
	pool.Put(bytes.NewBuffer(make([]byte, 1024)))
	assert.Equal(t, pool.Get(900).Cap(), 1024)
	assert.Equal(t, pool.Get(800).Cap(), 1024)
	assert.Equal(t, pool.Get(600).Cap(), 1024)

	pool.shards[1024].Put(bytes.NewBuffer(make([]byte, 800)))
	assert.GreaterOrEqual(t, pool.Get(900).Cap(), 1024)
}

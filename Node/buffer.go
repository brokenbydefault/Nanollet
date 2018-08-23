package Node

import (
	"sync"
	"github.com/brokenbydefault/Nanollet/Node/Packets"
)

var (
	BytesBuffPool = NewBytesBufferPool(Packets.MessageSize, 0)
)

type BytesBufferPool struct {
	Pool sync.Pool
}

func NewBytesBufferPool(len int, cap int) *BytesBufferPool {
	return &BytesBufferPool{
		Pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, len, cap)
			},
		},
	}
}

func (p BytesBufferPool) Get() []byte {
	return p.Pool.Get().([]byte)
}

func (p BytesBufferPool) Put(b []byte) {
	p.Pool.Put(b)
}

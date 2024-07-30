// Package buffer implements a buffer for serialization, consisting of a chain of []byte-s to
// reduce copying and to allow reuse of individual chunks.
package buffer

import (
	"fmt"
	"io"
	"strconv"
	"sync"
)

// PoolConfig contains configuration for the allocation and reuse strategy.
type PoolConfig struct {
	PooledSize int // Minimum chunk size that is reused, reusing chunks too small will result in overhead.
}

var config = PoolConfig{
	PooledSize: 1024,
}

var buffers *sync.Pool

func initBuffers() {
	buffers = &sync.Pool{
		New: func() any {
			pb := make([]byte, 0, config.PooledSize)
			return &pb
		},
	}
}

func init() {
	initBuffers()
}

// Init sets up a non-default pooling and allocation strategy. Should be run before serialization is done.
func Init(cfg PoolConfig) {
	config = cfg
	initBuffers()
}

func putBuf(buf *[]byte) {
	b := (*buf)[:0]
	buffers.Put(&b)
}

func getBuf() *[]byte {
	v := buffers.Get()
	return v.(*[]byte)
}

// Buffer is a buffer optimized for serialization without extra copying.
type Buffer struct {

	// Buf is the current chunk that can be used for serialization.
	buf *[]byte

	// Data to be dumped
	bufs []*[]byte
}

// EnsureSpace makes sure that the current chunk contains at least s free bytes,
// possibly creating a new chunk.
func (b *Buffer) EnsureSpace(s int) {
	if b.buf == nil {
		if s > config.PooledSize {
			panic(fmt.Sprintf("require space %d large than %d", s, config.PooledSize))
		}
		b.buf = getBuf()
		// 只能存储8个byte slice
		// TODO FIXME
		b.bufs = make([]*[]byte, 0, 8)
	} else if cap(*b.buf)-len(*b.buf) < s {
		// 放到bufs
		b.bufs = append(b.bufs, b.buf)
		// get free buf
		b.buf = getBuf()
	}
}

// AppendByte appends a single byte to buffer.
func (b *Buffer) AppendByte(data byte) {
	b.EnsureSpace(1)
	t := append(*b.buf, data)
	b.buf = &t
}

// AppendBytes appends a byte slice to buffer.
func (b *Buffer) AppendBytes(data []byte) {
	b.EnsureSpace(1)
	if len(data) <= cap(*b.buf)-len(*b.buf) {
		t := append(*b.buf, data...)
		b.buf = &t
	} else {
		b.appendBytesSlow(data)
	}
}

func (b *Buffer) appendBytesSlow(data []byte) {
	for len(data) > 0 {
		b.EnsureSpace(1)

		sz := cap(*b.buf) - len(*b.buf)
		if sz > len(data) {
			sz = len(data)
		}

		t := append(*b.buf, data[:sz]...)
		b.buf = &t
		data = data[sz:]
	}
}

func (b *Buffer) AppendString(data string) {
	b.EnsureSpace(1)
	if len(data) <= cap(*b.buf)-len(*b.buf) {
		t := append(*b.buf, data...)
		b.buf = &t
	} else {
		b.appendStringSlow(data)
	}
}

func (b *Buffer) appendStringSlow(data string) {
	for len(data) > 0 {
		b.EnsureSpace(1)

		sz := cap(*b.buf) - len(*b.buf)
		if sz > len(data) {
			sz = len(data)
		}

		t := append(*b.buf, data[:sz]...)
		b.buf = &t
		data = data[sz:]
	}
}

func (b *Buffer) AppendInt(n int64, dataSize int) {
	if dataSize > cap(*b.buf)-len(*b.buf) {
		b.EnsureSpace(dataSize)
	}
	ret := strconv.AppendInt(*b.buf, n, 10)
	b.buf = &ret
}

func (b *Buffer) AppendUint(n uint64, dataSize int) {
	if dataSize > cap(*b.buf)-len(*b.buf) {
		b.EnsureSpace(dataSize)
	}
	ret := strconv.AppendUint(*b.buf, n, 10)
	b.buf = &ret
}

func (b *Buffer) AppendFloat32(f float32, dataSize int) {
	if dataSize > cap(*b.buf)-len(*b.buf) {
		b.EnsureSpace(dataSize)
	}
	ret := strconv.AppendFloat(*b.buf, float64(f), 'g', -1, 32)
	b.buf = &ret
}

func (b *Buffer) AppendFloat64(f float64, dataSize int) {
	if dataSize > cap(*b.buf)-len(*b.buf) {
		b.EnsureSpace(dataSize)
	}
	ret := strconv.AppendFloat(*b.buf, float64(f), 'g', -1, 64)
	b.buf = &ret
}

// Size computes the size of a buffer by adding sizes of every chunk.
func (b *Buffer) Size() int {
	// TODO: 在append的时候计数
	size := len(*b.buf)
	for _, buf := range b.bufs {
		size += len(*buf)
	}
	return size
}

// DumpTo outputs the contents of a buffer to a writer and resets the buffer.
func (b *Buffer) DumpTo(w io.Writer) (int, error) {
	written := 0
	for _, buf := range b.bufs {
		if n, err := w.Write(*buf); err != nil {
			return n, err
		} else {
			written += n
		}
	}
	if cap(*b.buf) > 0 {
		if n, err := w.Write(*b.buf); err != nil {
			return n, err
		} else {
			written += n
		}
	}

	// free bufs
	for _, buf := range b.bufs {
		putBuf(buf)
	}

	putBuf(b.buf)
	b.bufs = nil
	b.buf = nil

	return written, nil
}

func (b *Buffer) BuildBytes() []byte {
	if len(b.bufs) == 0 {
		cpy := make([]byte, len(*b.buf))
		copy(cpy, *b.buf)
		putBuf(b.buf)
		return cpy
	}

	var ret []byte
	size := b.Size()

	ret = make([]byte, size)

	written := 0
	for _, buf := range b.bufs {
		written += copy(ret[written:], *buf)
		putBuf(buf)
	}

	copy(ret[written:], *b.buf)
	putBuf(b.buf)

	b.bufs = nil
	b.buf = nil

	return ret
}

func (b *Buffer) Close() error {
	if b.bufs != nil {
		// Release all remaining buffers.
		for _, buf := range b.bufs {
			putBuf(buf)
		}
		// In case Close gets called multiple times.
		b.bufs = nil
	}

	if b.buf != nil {
		putBuf(b.buf)
		b.buf = nil
	}

	return nil
}

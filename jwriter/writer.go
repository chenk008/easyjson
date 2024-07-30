// Package jwriter contains a JSON writer.
package jwriter

import (
	"bytes"
	"github.com/chenk008/easyjson/buffer"
)

// Flags describe various encoding options. The behavior may be actually implemented in the encoder, but
// Flags field in Writer is used to set and pass them around.
type Flags int

const (
	NilMapAsEmpty   Flags = 1 << iota // Encode nil map as '{}' rather than 'null'.
	NilSliceAsEmpty                   // Encode nil slice as '[]' rather than 'null'.
)

type Writer interface {
	Flags() Flags

	MaybeFlush() (int, error)

	Flush() (int, error)

	// Close resets the buffer.
	Close() error

	// RawByte appends raw binary data to the buffer.
	RawByte(c byte) error

	// RawByte appends raw binary data to the buffer.
	RawBytes(data []byte) error

	RawBytesWithErr(data []byte, err error) error

	RawTextWithErr(data []byte, err error) error

	// RawByte appends raw binary data to the buffer.
	RawString(s string) error

	// Base64Bytes appends data to the buffer after base64 encoding it
	Base64Bytes(data []byte) error

	Uint8(n uint8) error

	Uint16(n uint16) error

	Uint32(n uint32) error

	Uint(n uint) error

	Uint64(n uint64) error

	Int8(n int8) error

	Int16(n int16) error

	Int32(n int32) error

	Int(n int) error

	Int64(n int64) error

	Uint8Str(n uint8) error

	Uint16Str(n uint16) error

	Uint32Str(n uint32) error

	UintStr(n uint) error

	Uint64Str(n uint64) error

	UintptrStr(n uintptr) error

	Int8Str(n int8) error

	Int16Str(n int16) error

	Int32Str(n int32) error

	IntStr(n int) error

	Int64Str(n int64) error

	Float32(n float32) error

	Float32Str(n float32) error

	Float64(n float64) error

	Float64Str(n float64) error

	Bool(v bool) error

	String(s string) error
}

const chars = "0123456789abcdef"

func getTable(falseValues ...int) [128]bool {
	table := [128]bool{}

	for i := 0; i < 128; i++ {
		table[i] = true
	}

	for _, v := range falseValues {
		table[v] = false
	}

	return table
}

var (
	htmlEscapeTable   = getTable(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, '"', '&', '<', '>', '\\')
	htmlNoEscapeTable = getTable(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, '"', '\\')
)

const encode = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
const padChar = '='

func base64(b *buffer.Buffer, in []byte) {

	if len(in) == 0 {
		return
	}

	b.EnsureSpace(((len(in)-1)/3 + 1) * 4)

	si := 0
	n := (len(in) / 3) * 3

	for si < n {
		// Convert 3x 8bit source bytes into 4 bytes
		val := uint(in[si+0])<<16 | uint(in[si+1])<<8 | uint(in[si+2])

		b.AppendByte(encode[val>>18&0x3F])
		b.AppendByte(encode[val>>12&0x3F])
		b.AppendByte(encode[val>>6&0x3F])
		b.AppendByte(encode[val&0x3F])

		si += 3
	}

	remain := len(in) - si
	if remain == 0 {
		return
	}

	// Add the remaining small block
	val := uint(in[si+0]) << 16
	if remain == 2 {
		val |= uint(in[si+1]) << 8
	}

	b.AppendByte(encode[val>>18&0x3F])
	b.AppendByte(encode[val>>12&0x3F])

	switch remain {
	case 2:
		b.AppendByte(encode[val>>6&0x3F])
		b.AppendByte(byte(padChar))
	case 1:
		b.AppendByte(byte(padChar))
		b.AppendByte(byte(padChar))
	}
}

// it is for test
type BufferWriter struct {
	tokenWriter
	buf         *bytes.Buffer
	currentByte []byte
}

func NewBufferWriter() *BufferWriter {
	var b bytes.Buffer
	streamWriter := tokenWriter{
		targetIOWriter:   &b,
		targetBufferSize: 1024,
	}
	streamWriter.buffer.EnsureSpace(1024)

	return &BufferWriter{
		tokenWriter: streamWriter,
		buf:         &b,
	}
}

func (b *BufferWriter) Flush() (int, error) {
	n, err := b.tokenWriter.Flush()
	if err == nil {
		b.currentByte = b.buf.Bytes()
		b.buf.Reset()
	}
	return n, err
}

func (b *BufferWriter) Size() int {
	return len(b.currentByte)
}

func (b *BufferWriter) BuildBytes() []byte {
	return b.currentByte
}

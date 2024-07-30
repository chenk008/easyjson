package jwriter

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBufferWriter(t *testing.T) {
	w := NewBufferWriter()

	var b bytes.Buffer
	for i := 0; i < 1002; i++ {
		b.WriteByte(1)
		b.WriteByte(2)
	}

	err := w.RawBytes(b.Bytes())
	b.Reset()
	if err != nil {
		t.Error(err)
	}
	written, err := w.MaybeFlush()
	if err != nil {
		t.Error(err)
	}

	if written != 2004 {
		t.Error(fmt.Errorf("got:%d, expect:%d", written, 2004))
	}

	// case2
	for i := 0; i < 500; i++ {
		b.WriteByte(1)
		b.WriteByte(2)
	}
	w = NewBufferWriter()
	err = w.RawBytes(b.Bytes())
	b.Reset()
	if err != nil {
		t.Error(err)
	}
	written, err = w.MaybeFlush()
	if err != nil {
		t.Error(err)
	}

	if written != 0 {
		t.Error(fmt.Errorf("got:%d, expect:%d", written, 0))
	}

	written, err = w.Flush()
	if err != nil {
		t.Error(err)
	}
	if written != 1000 {
		t.Error(fmt.Errorf("got:%d, expect:%d", written, 1000))
	}

	// case3
	for i := 0; i < 2022; i++ {
		b.WriteByte(1)
		b.WriteByte(2)
	}
	w = NewBufferWriter()
	err = w.RawBytes(b.Bytes())
	b.Reset()
	if err != nil {
		t.Error(err)
	}
	written, err = w.MaybeFlush()
	if err != nil {
		t.Error(err)
	}

	if written != 4044 {
		t.Error(fmt.Errorf("got:%d, expect:%d", written, 4044))
	}

	for i := 0; i < 1000; i++ {
		b.WriteByte(1)
		b.WriteByte(2)
	}
	err = w.RawBytes(b.Bytes())
	b.Reset()

	written, err = w.Flush()
	if err != nil {
		t.Error(err)
	}
	if written != 2000 {
		t.Error(fmt.Errorf("got:%d, expect:%d", written, 2000))
	}
}

//go:build use_easyjson
// +build use_easyjson

package benchmark

import (
	"testing"

	"github.com/chenk008/easyjson"
	"github.com/chenk008/easyjson/jwriter"
)

func BenchmarkEJ_Unmarshal_M(b *testing.B) {
	b.SetBytes(int64(len(largeStructText)))
	for i := 0; i < b.N; i++ {
		var s LargeStruct
		err := s.UnmarshalJSON(largeStructText)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkEJ_Unmarshal_S(b *testing.B) {
	b.SetBytes(int64(len(smallStructText)))

	for i := 0; i < b.N; i++ {
		var s Entities
		err := s.UnmarshalJSON(smallStructText)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkEJ_Marshal_M(b *testing.B) {
	var l int64
	for i := 0; i < b.N; i++ {
		data, err := easyjson.Marshal(&largeStructData)
		if err != nil {
			b.Error(err)
		}
		l = int64(len(data))
	}
	b.SetBytes(l)
}

func BenchmarkEJ_Marshal_L(b *testing.B) {
	var l int64
	for i := 0; i < b.N; i++ {
		data, err := easyjson.Marshal(&xlStructData)
		if err != nil {
			b.Error(err)
		}
		l = int64(len(data))
	}
	b.SetBytes(l)
}

func BenchmarkEJ_Marshal_L_ToWriter(b *testing.B) {
	var l int64
	out := &DummyWriter{}
	for i := 0; i < b.N; i++ {
		w := jwriter.NewStreamingTokenWriter(out, 1024)
		err := xlStructData.MarshalEasyJSON(w)
		if err != nil {
			b.Error(err)
		}
		if written, err := w.Flush(); err != nil {
			b.Error(err)
		} else {
			l = int64(written)
		}
	}
	b.SetBytes(l)

}
func BenchmarkEJ_Marshal_M_Parallel(b *testing.B) {
	b.SetBytes(int64(len(largeStructText)))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := largeStructData.MarshalJSON()
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkEJ_Marshal_M_ToWriter(b *testing.B) {
	var l int64
	out := &DummyWriter{}
	for i := 0; i < b.N; i++ {
		w := jwriter.NewStreamingTokenWriter(out, 1024)
		err := largeStructData.MarshalEasyJSON(w)
		if err != nil {
			b.Error(err)
		}
		if written, err := w.Flush(); err != nil {
			b.Error(err)
		} else {
			l = int64(written)
		}

	}
	b.SetBytes(l)

}
func BenchmarkEJ_Marshal_M_ToWriter_Parallel(b *testing.B) {
	out := &DummyWriter{}

	b.RunParallel(func(pb *testing.PB) {
		var l int64
		for pb.Next() {
			w := jwriter.NewStreamingTokenWriter(out, 1024)
			err := largeStructData.MarshalEasyJSON(w)
			if err != nil {
				b.Error(err)
			}

			if written, err := w.Flush(); err != nil {
				b.Error(err)
			} else {
				l = int64(written)
			}
		}
		if l > 0 {
			b.SetBytes(l)
		}
	})

}

func BenchmarkEJ_Marshal_L_Parallel(b *testing.B) {
	var l int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data, err := xlStructData.MarshalJSON()
			if err != nil {
				b.Error(err)
			}
			l = int64(len(data))
		}
	})
	b.SetBytes(l)
}

func BenchmarkEJ_Marshal_L_ToWriter_Parallel(b *testing.B) {
	out := &DummyWriter{}
	b.RunParallel(func(pb *testing.PB) {
		var l int64
		for pb.Next() {
			w := jwriter.NewStreamingTokenWriter(out, 1024)

			err := xlStructData.MarshalEasyJSON(w)
			if err != nil {
				b.Error(err)
			}
			if written, err := w.Flush(); err != nil {
				b.Error(err)
			} else {
				l = int64(written)
			}
		}
		if l > 0 {
			b.SetBytes(l)
		}
	})
}

func BenchmarkEJ_Marshal_S(b *testing.B) {
	var l int64
	for i := 0; i < b.N; i++ {
		data, err := smallStructData.MarshalJSON()
		if err != nil {
			b.Error(err)
		}
		l = int64(len(data))
	}
	b.SetBytes(l)
}

func BenchmarkEJ_Marshal_S_Parallel(b *testing.B) {
	var l int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data, err := smallStructData.MarshalJSON()
			if err != nil {
				b.Error(err)
			}
			l = int64(len(data))
		}
	})
	b.SetBytes(l)
}

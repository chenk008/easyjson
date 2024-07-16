package tests

import (
	"github.com/chenk008/easyjson"
	"github.com/chenk008/easyjson/jlexer"
	"github.com/chenk008/easyjson/jwriter"
)

//easyjson:json
type NestedMarshaler struct {
	Value  easyjson.MarshalerUnmarshaler
	Value2 int
}

type StructWithMarshaler struct {
	Value int
}

func (s *StructWithMarshaler) UnmarshalEasyJSON(w *jlexer.Lexer) {
	s.Value = w.Int()
}

func (s *StructWithMarshaler) MarshalEasyJSON(w jwriter.Writer) error {
	return w.Int(s.Value)
}

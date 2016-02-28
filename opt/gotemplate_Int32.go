package opt

import (
	"fmt"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// template type Optional(A)

// A 'gotemplate'-based type for providing optional semantics without using pointers.
type Int32 struct {
	V       int32
	Defined bool
}

// Creates an optional type with a given value.
func OInt32(v int32) Int32 {
	return Int32{V: v, Defined: true}
}

// Get returns the value or given default in the case the value is undefined.
func (v Int32) Get(deflt int32) int32 {
	if !v.Defined {
		return deflt
	}
	return v.V
}

// MarshalEasyJSON does JSON marshaling using easyjson interface.
func (v Int32) MarshalEasyJSON(w *jwriter.Writer) {
	if v.Defined {
		w.Int32(v.V)
	} else {
		w.RawString("null")
	}
}

// UnmarshalEasyJSON does JSON unmarshaling using easyjson interface.
func (v *Int32) UnmarshalEasyJSON(l *jlexer.Lexer) {
	if l.IsNull() {
		l.Skip()
		*v = Int32{}
	} else {
		v.V = l.Int32()
		v.Defined = true
	}
}

// MarshalJSON implements a standard json marshaler interface.
func (v *Int32) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	v.MarshalEasyJSON(&w)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalJSON implements a standard json marshaler interface.
func (v *Int32) UnmarshalJSON(data []byte) error {
	l := jlexer.Lexer{}
	v.UnmarshalEasyJSON(&l)
	return l.Error()
}

// IsDefined returns whether the value is defined, a function is required so that it can
// be used in an interface.
func (v Int32) IsDefined() bool {
	return v.Defined
}

// String implements a stringer interface using fmt.Sprint for the value.
func (v Int32) String() string {
	if !v.Defined {
		return "<undefined>"
	}
	return fmt.Sprint(v.V)
}
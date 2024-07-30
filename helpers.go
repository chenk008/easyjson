// Package easyjson contains marshaler/unmarshaler interfaces and helper functions.
package easyjson

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"unsafe"

	"github.com/chenk008/easyjson/jlexer"
	"github.com/chenk008/easyjson/jwriter"
)

// Marshaler is an easyjson-compatible marshaler interface.
type Marshaler interface {
	MarshalEasyJSON(w jwriter.Writer) error
}

// Marshaler is an easyjson-compatible unmarshaler interface.
type Unmarshaler interface {
	UnmarshalEasyJSON(w *jlexer.Lexer)
}

// MarshalerUnmarshaler is an easyjson-compatible marshaler/unmarshaler interface.
type MarshalerUnmarshaler interface {
	Marshaler
	Unmarshaler
}

// Optional defines an undefined-test method for a type to integrate with 'omitempty' logic.
type Optional interface {
	IsDefined() bool
}

// UnknownsUnmarshaler provides a method to unmarshal unknown struct fileds and save them as you want
type UnknownsUnmarshaler interface {
	UnmarshalUnknown(in *jlexer.Lexer, key string)
}

// UnknownsMarshaler provides a method to write additional struct fields
type UnknownsMarshaler interface {
	MarshalUnknowns(w jwriter.Writer, first bool) error
}

func isNilInterface(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}

// Marshal returns data as a single byte slice. Method is suboptimal as the data is likely to be copied
// from a chain of smaller chunks.
func Marshal(v Marshaler) ([]byte, error) {
	if isNilInterface(v) {
		return nullBytes, nil
	}

	bufferWriter := jwriter.NewBufferWriter()
	if err := v.MarshalEasyJSON(bufferWriter); err != nil {
		return nil, err
	}
	if _, err := bufferWriter.Flush(); err != nil {
		return nil, err
	}
	return bufferWriter.BuildBytes(), nil
}

// MarshalToWriter marshals the data to an io.Writer.
func MarshalToWriter(v Marshaler, w io.Writer) (written int, err error) {
	if isNilInterface(v) {
		return w.Write(nullBytes)
	}

	bufferWriter := jwriter.NewBufferWriter()
	if err := v.MarshalEasyJSON(bufferWriter); err != nil {
		return 0, err
	}
	if _, err := bufferWriter.Flush(); err != nil {
		return 0, err
	}
	return w.Write(bufferWriter.BuildBytes())
}

// MarshalToHTTPResponseWriter sets Content-Length and Content-Type headers for the
// http.ResponseWriter, and send the data to the writer. started will be equal to
// false if an error occurred before any http.ResponseWriter methods were actually
// invoked (in this case a 500 reply is possible).
func MarshalToHTTPResponseWriter(v Marshaler, w http.ResponseWriter) (started bool, written int, err error) {
	if isNilInterface(v) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(nullBytes)))
		written, err = w.Write(nullBytes)
		return true, written, err
	}

	bufferWriter := jwriter.NewBufferWriter()
	if err := v.MarshalEasyJSON(bufferWriter); err != nil {
		return false, 0, err
	}
	if _, err := bufferWriter.Flush(); err != nil {
		return false, 0, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(bufferWriter.Size()))

	started = true
	written, err = w.Write(bufferWriter.BuildBytes())
	return
}

// Unmarshal decodes the JSON in data into the object.
func Unmarshal(data []byte, v Unmarshaler) error {
	l := jlexer.Lexer{Data: data}
	v.UnmarshalEasyJSON(&l)
	return l.Error()
}

// UnmarshalFromReader reads all the data in the reader and decodes as JSON into the object.
func UnmarshalFromReader(r io.Reader, v Unmarshaler) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	l := jlexer.Lexer{Data: data}
	v.UnmarshalEasyJSON(&l)
	return l.Error()
}

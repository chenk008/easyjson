package easyjson

import (
	jlexer "github.com/chenk008/easyjson/jlexer"
	"github.com/chenk008/easyjson/jwriter"
)

// UnknownFieldsProxy implemets UnknownsUnmarshaler and UnknownsMarshaler
// use it as embedded field in your structure to parse and then serialize unknown struct fields
type UnknownFieldsProxy struct {
	unknownFields map[string][]byte
}

func (s *UnknownFieldsProxy) UnmarshalUnknown(in *jlexer.Lexer, key string) {
	if s.unknownFields == nil {
		s.unknownFields = make(map[string][]byte, 1)
	}
	s.unknownFields[key] = in.Raw()
}

func (s UnknownFieldsProxy) MarshalUnknowns(out jwriter.Writer, first bool) error {
	for key, val := range s.unknownFields {
		if first {
			first = false
		} else {
			if err := out.RawByte(','); err != nil {
				return err
			}
		}
		if err := out.String(string(key)); err != nil {
			return err
		}
		if err := out.RawByte(':'); err != nil {
			return err
		}
		return out.RawBytes(val)
	}
	return nil
}

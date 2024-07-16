package tests

import (
	"github.com/chenk008/easyjson"
	"github.com/chenk008/easyjson/jwriter"
)

//easyjson:json
type NestedInterfaces struct {
	Value interface{}
	Slice []interface{}
	Map   map[string]interface{}
}

type NestedEasyMarshaler struct {
	EasilyMarshaled bool
}

var _ easyjson.Marshaler = &NestedEasyMarshaler{}

func (i *NestedEasyMarshaler) MarshalEasyJSON(w jwriter.Writer) error {
	// We use this method only to indicate that easyjson.Marshaler
	// interface was really used while encoding.
	i.EasilyMarshaled = true
	return nil
}

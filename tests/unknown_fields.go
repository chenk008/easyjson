package tests

import "github.com/chenk008/easyjson"

//easyjson:json
type StructWithUnknownsProxy struct {
	easyjson.UnknownFieldsProxy

	Field1 string
}

//easyjson:json
type StructWithUnknownsProxyWithOmitempty struct {
	easyjson.UnknownFieldsProxy

	Field1 string `json:",omitempty"`
}

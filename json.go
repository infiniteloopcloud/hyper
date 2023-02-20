package hyper

import (
	"encoding/json"
	"io"
)

var jsonEncoder JSONEncoder = defaultJSONParser{}

type JSONEncoder interface {
	Decode(r io.Reader, v any) error
	Encode(w io.Writer, v any) error
}

func SetJSONEncoder(e JSONEncoder) {
	jsonEncoder = e
}

type defaultJSONParser struct{}

func (defaultJSONParser) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

func (defaultJSONParser) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

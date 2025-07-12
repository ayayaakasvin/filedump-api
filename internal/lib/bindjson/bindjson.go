package bindjson

import (
	"encoding/json"
	"io"
)

func BindJson(r io.Reader, obj any) error {
	return json.NewDecoder(r).Decode(obj)
}
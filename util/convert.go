package util

import (
	"bytes"
	"encoding/json"
	"io"
)

// StructToReader takes any struct and calls json.Marshal on it before returning an io.Reader wrapped around the returned
// bytes.
func StructToReader(v any) (io.Reader, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

// StructToReaderOrDie takes any struct and calls json.Marshal on it before returning an io.Reader wrapped around the returned
// bytes. If the json.Marshal fails, this function panics.
func StructToReaderOrDie(v any) io.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return bytes.NewReader(b)
}

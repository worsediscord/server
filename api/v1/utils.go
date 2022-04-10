package v1

import (
	"fmt"
	"io"
	"strings"
)

func ReadHeader(k string, m map[string][]string) (string, error) {
	if v, ok := m[k]; ok {
		if len(v) >= 1 {
			return strings.Join(v, ","), nil
		} else {
			return "", fmt.Errorf("key is empty")
		}
	}

	return "", fmt.Errorf("key not found")
}

func WriteState(r io.Reader, w io.Writer) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	_, err = w.Write(b)

	return err
}

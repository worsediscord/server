package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Flusher interface {
	Listen(ctx context.Context) (chan<- interface{}, <-chan error)
	Flatten(pk string, v interface{}) interface{}
}

type FileFlusher struct {
	BasePath string
}

type fileData struct {
	filepath string
	data     []byte
}

func NewFileFlusher(basePath string) (*FileFlusher, error) {
	err := os.MkdirAll(basePath, 0700)
	if err != nil {
		return nil, fmt.Errorf("could not open path %s: %w", basePath, err)
	}

	return &FileFlusher{BasePath: basePath}, nil
}

func (f *FileFlusher) Listen(ctx context.Context) (chan<- interface{}, <-chan error) {
	vChan := make(chan interface{})
	errChan := make(chan error)

	go func() {
		for {
			select {
			case v := <-vChan:
				if fd, ok := v.(fileData); ok {
					path := filepath.Join(f.BasePath, fd.filepath)
					errChan <- os.WriteFile(path, fd.data, 0700)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return vChan, errChan
}

func (f *FileFlusher) Flatten(pk string, v interface{}) interface{} {
	b, _ := json.Marshal(v)

	return fileData{
		filepath: pk,
		data:     b,
	}
}

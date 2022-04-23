package memcache

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

const (
	MethodSet    = "SET"
	MethodGet    = "GET"
	MethodList   = "LIST"
	MethodDelete = "DELETE"
	MethodSave   = "SAVE"
	MethodFlush  = "FLUSH"
)

type Item struct {
	Key   string
	Value []byte
}

type Document struct {
	Name  string
	state map[string][]byte
	lock  sync.RWMutex
}

type Datastore struct {
	documents map[string]*Document
	lock      sync.RWMutex
}

type Request struct {
	Method string
	Item
}

type DocumentReader interface {
	Get(key string) Item
	GetAll() []Item
}

type DocumentWriter interface {
	Set(key string, value interface{}) error
	Delete(key string)
}

type DocumentReadWriter interface {
	DocumentReader
	DocumentWriter
}

func (i Item) Decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &ErrInvalidDecode{reflect.TypeOf(v)}
	}

	if len(i.Value) == 0 {
		return ErrEmptyItem
	}

	return gob.NewDecoder(bytes.NewBuffer(i.Value)).Decode(v)
}

func newDocument(name string) *Document {
	return &Document{
		Name:  name,
		state: make(map[string][]byte),
	}
}

func (d *Document) Get(key string) Item {
	d.lock.RLock()
	defer d.lock.RUnlock()

	if key == "" {
		return Item{}
	}

	if _, ok := d.state[key]; !ok {
		return Item{}
	}

	return Item{Key: key, Value: d.state[key]}
}

func (d *Document) GetAll() []Item {
	d.lock.RLock()
	defer d.lock.RUnlock()

	var items []Item
	for k, v := range d.state {
		items = append(items, Item{Key: k, Value: v})
	}

	return items
}

func (d *Document) Set(key string, value interface{}) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if key == "" {
		// TODO probs better error
		return ErrKeyNotFound
	}

	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return fmt.Errorf("could not store value in document")
	}

	d.state[key] = b.Bytes()

	return nil
}

func (d *Document) RawSet(key string, b []byte) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if key == "" {
		return ErrKeyNotFound
	}

	d.state[key] = b

	return nil
}

func (d *Document) Delete(key string) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if key == "" {
		return
	}

	if _, ok := d.state[key]; !ok {
		return
	}

	delete(d.state, key)

	return
}

func NewDatastore() *Datastore {
	return &Datastore{
		documents: make(map[string]*Document),
	}
}

func Open(p string) (*Datastore, error) {
	datastore := NewDatastore()
	baseDepth := strings.Count(p, string(os.PathSeparator))

	stat, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return datastore, os.MkdirAll(p, 0700)
		} else {
			return datastore, err
		}
	}

	if !stat.IsDir() {
		return datastore, ErrInvalidPath
	}

	err = filepath.WalkDir(p, func(walkPath string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if de.IsDir() {
			if strings.Count(walkPath, string(os.PathSeparator))-baseDepth > 2 {
				return fmt.Errorf("failed open datastore")
			}
			datastore.Document(de.Name())
			return nil
		}

		dirSplit := strings.Split(filepath.Dir(walkPath), string(os.PathSeparator))
		documentName := dirSplit[len(dirSplit)-1]

		b, err := os.ReadFile(walkPath)
		if err != nil {
			return err
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(string(b))
		if err != nil {
			return err
		}

		if err = datastore.Document(documentName).RawSet(filepath.Base(walkPath), decodedBytes); err != nil {
			return err
		}

		return nil
	})

	return datastore, err
}

func (d *Datastore) Document(name string) *Document {
	d.lock.RLock()

	if _, ok := d.documents[name]; !ok {
		d.lock.RUnlock()
		d.lock.Lock()
		d.documents[name] = newDocument(name)
		d.lock.Unlock()
	} else {
		d.lock.RUnlock()
	}

	return d.documents[name]
}

func (d *Datastore) Close(p string) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	for dir, document := range d.documents {
		if err := os.MkdirAll(path.Join(p, dir), 0700); err != nil {
			return err
		}

		for file, data := range document.state {
			if err := os.WriteFile(path.Join(p, dir, file), []byte(base64.StdEncoding.EncodeToString(data)), 0600); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *Datastore) Listen(ctx context.Context) (chan<- Request, <-chan error) {
	reqChan := make(chan Request, 10)
	errChan := make(chan error, 10)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}()

	return reqChan, errChan
}

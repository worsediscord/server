package memcache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"sync"
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

func NewDocument(name string) *Document {
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

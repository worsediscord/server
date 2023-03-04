package memcache

import (
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Collection struct {
	name      string
	documents map[string]*Document
	lock      sync.RWMutex
}

type Datastore struct {
	collections map[string]*Collection
	lock        sync.RWMutex
}

type CollectionReader interface {
	Get(key string) *Document
	GetAll() []*Document
}

type CollectionWriter interface {
	Set(key string, document *Document) error
	Delete(key string)
}

type CollectionReadWriter interface {
	CollectionReader
	CollectionWriter
}

func newCollection(name string) *Collection {
	return &Collection{
		name:      name,
		documents: make(map[string]*Document),
	}
}

func (c *Collection) Get(key string) *Document {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if _, ok := c.documents[key]; ok {
		return c.documents[key]
	}

	return nil
}

func (c *Collection) GetAll() []*Document {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var documents []*Document
	for _, document := range c.documents {
		documents = append(documents, document)
	}

	return documents
}

func (c *Collection) Set(key string, document *Document) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if key == "" {
		// TODO probs better error
		return ErrKeyNotFound
	}

	if document == nil {
		return fmt.Errorf("cannot insert nil document into collection")
	}

	c.documents[key] = document

	return nil
}

func (c *Collection) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.documents[key]; ok {
		delete(c.documents, key)
	}
}

// Document is a helper function that returns an existing document if it exists, and creates it if it doesn't.
func (c *Collection) Document(name string) *Document {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.documents[name]; !ok {
		c.documents[name] = NewDocument(name)
	}

	return c.documents[name]
}

func NewDatastore() *Datastore {
	return &Datastore{
		collections: make(map[string]*Collection),
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

	var currentCollection string
	err = filepath.WalkDir(p, func(walkPath string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if de.IsDir() {
			switch strings.Count(walkPath, string(os.PathSeparator)) - baseDepth {
			case 1:
				datastore.Collection(de.Name())
				currentCollection = de.Name()
			case 2:
				datastore.Collection(currentCollection).Document(de.Name())
			case 3:
				return fmt.Errorf("failed open datastore")
			}

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

		if err = datastore.Collection(currentCollection).Document(documentName).RawSet(filepath.Base(walkPath), decodedBytes); err != nil {
			return err
		}

		return nil
	})

	return datastore, err
}

func (d *Datastore) Collection(name string) *Collection {
	d.lock.Lock()
	defer d.lock.Unlock()

	if _, ok := d.collections[name]; !ok {
		d.collections[name] = newCollection(name)
	}

	return d.collections[name]
}

func (d *Datastore) Close(outPath string) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	for collectionName, collection := range d.collections {
		for documentName, document := range collection.documents {
			documentPath := path.Join(outPath, collectionName, documentName)
			if err := os.MkdirAll(documentPath, 0700); err != nil {
				return err
			}
			for file, data := range document.state {
				if err := os.WriteFile(path.Join(documentPath, file), []byte(base64.StdEncoding.EncodeToString(data)), 0600); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

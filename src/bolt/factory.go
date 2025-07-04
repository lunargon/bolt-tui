package bolt

import (
	"fmt"
	"sync"
)

type Factory struct {
	lck       sync.RWMutex
	databases map[string]*DB
}

func NewFactory(name, defaultPath string) (*Factory, error) {
	databases := make(map[string]*DB)
	databases[name] = &DB{Path: defaultPath}
	if err := databases[name].Open(); err != nil {
		return nil, fmt.Errorf("could not open database %s: %v", name, err)
	}
	return &Factory{databases: databases}, nil
}

func (f *Factory) GetDatabases() ([]string, error) {
	f.lck.RLock()
	defer f.lck.RUnlock()

	databases := make([]string, 0, len(f.databases))
	for name := range f.databases {
		databases = append(databases, name)
	}
	return databases, nil
}

func (f *Factory) Open(name, path string) (*DB, error) {
	f.lck.Lock()
	defer f.lck.Unlock()

	db := &DB{Path: path}
	if err := db.Open(); err != nil {
		return nil, err
	}

	f.databases[name] = db
	return db, nil
}

func (f *Factory) Close(name string) error {
	f.lck.Lock()
	defer f.lck.Unlock()

	db, ok := f.databases[name]
	if !ok {
		return fmt.Errorf("database %s not found", name)
	}

	if err := db.Close(); err != nil {
		return err
	}

	delete(f.databases, name)
	return nil
}

func (f *Factory) Get(name string) (*DB, error) {
	f.lck.RLock()
	defer f.lck.RUnlock()

	db, ok := f.databases[name]
	if !ok {
		return nil, fmt.Errorf("database %s not found", name)
	}

	return db, nil
}

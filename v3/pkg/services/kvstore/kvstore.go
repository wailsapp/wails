package kvstore

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type KeyValueStore struct {
	config   *Config
	filename string
	data     map[string]any
	unsaved  bool
	lock     sync.RWMutex
}

type Config struct {
	Filename string
	AutoSave bool
}

type Service struct{}

func New(config *Config) *KeyValueStore {
	return &KeyValueStore{
		config: config,
		data:   make(map[string]any),
	}
}

// ServiceShutdown will save the store to disk if there are unsaved changes.
func (kvs *KeyValueStore) ServiceShutdown() error {
	if kvs.unsaved {
		err := kvs.Save()
		if err != nil {
			return errors.Wrap(err, "Error saving store")
		}
	}
	return nil
}

// ServiceName returns the name of the plugin.
func (kvs *KeyValueStore) ServiceName() string {
	return "github.com/wailsapp/wails/v3/plugins/kvstore"
}

// ServiceStartup is called when the plugin is loaded. This is where you should do any setup.
func (kvs *KeyValueStore) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	err := kvs.open(kvs.config.Filename)
	if err != nil {
		return err
	}

	return nil
}

func (kvs *KeyValueStore) open(filename string) (err error) {
	kvs.filename = filename
	kvs.data = make(map[string]any)

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer func() {
		err2 := file.Close()
		if err2 != nil {
			application.Get().Logger.Error("Key/Value Store Plugin Error:", "error", err.Error())
			if err == nil {
				err = err2
			}
		}
	}()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(bytes) > 0 {
		if err := json.Unmarshal(bytes, &kvs.data); err != nil {
			return err
		}
	}

	return nil
}

// Save saves the store to disk
func (kvs *KeyValueStore) Save() error {
	kvs.lock.Lock()
	defer kvs.lock.Unlock()

	bytes, err := json.Marshal(kvs.data)
	if err != nil {
		return err
	}

	err = os.WriteFile(kvs.filename, bytes, 0644)
	if err != nil {
		return err
	}

	kvs.unsaved = false

	return nil
}

// Get returns the value for the given key. If key is empty, the entire store is returned.
func (kvs *KeyValueStore) Get(key string) any {
	kvs.lock.RLock()
	defer kvs.lock.RUnlock()

	if key == "" {
		return kvs.data
	}

	return kvs.data[key]
}

// Set sets the value for the given key. If AutoSave is true, the store is saved to disk.
func (kvs *KeyValueStore) Set(key string, value any) error {
	kvs.lock.Lock()
	kvs.data[key] = value
	kvs.lock.Unlock()
	if kvs.config.AutoSave {
		err := kvs.Save()
		if err != nil {
			return err
		}
		kvs.unsaved = false
	} else {
		kvs.unsaved = true
	}
	return nil
}

// Delete deletes the key from the store. If AutoSave is true, the store is saved to disk.
func (kvs *KeyValueStore) Delete(key string) error {
	kvs.lock.Lock()
	delete(kvs.data, key)
	kvs.lock.Unlock()
	if kvs.config.AutoSave {
		err := kvs.Save()
		if err != nil {
			return err
		}
		kvs.unsaved = false
	} else {
		kvs.unsaved = true
	}
	return nil
}

package kvstore

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Config struct {
	// Filename specifies the path of the on-disk file associated to the key-value store.
	Filename string

	// AutoSave specifies whether the store
	// must be written to disk automatically after every modification.
	// When AutoSave is false, stores are only saved to disk upon shutdown
	// or when the [Service.Save] method is called manually.
	AutoSave bool
}

type KVStoreService struct {
	lock sync.RWMutex

	config *Config

	data    map[string]any
	unsaved bool
}

// New initialises an in-memory key-value store. See [NewWithConfig] for details.
func New() *KVStoreService {
	return NewWithConfig(nil)
}

// NewWithConfig initialises a key-value store with the given configuration:
//   - if config is nil, the new store is in-memory, i.e. not associated with a file;
//   - if config is non-nil, the associated file is not loaded until [Service.Load] is called.
//
// If the store is registered with the application as a service,
// [Service.Load] will be called automatically at startup.
func NewWithConfig(config *Config) *KVStoreService {
	result := &KVStoreService{data: make(map[string]any)}
	result.Configure(config)
	return result
}

// ServiceName returns the name of the plugin.
func (kvs *KVStoreService) ServiceName() string {
	return "github.com/wailsapp/wails/v3/plugins/kvstore"
}

// ServiceStartup loads the store from disk if it is associated with a file.
// It returns a non-nil error in case of failure.
func (kvs *KVStoreService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return errors.Wrap(kvs.Load(), "error loading store")
}

// ServiceShutdown saves the store to disk if it is associated with a file.
// It returns a non-nil error in case of failure.
func (kvs *KVStoreService) ServiceShutdown() error {
	return errors.Wrap(kvs.Save(), "error saving store")
}

// Configure changes the store's configuration.
// The contents of the store at call time are preserved and marked unsaved.
// Consumers will need to call [Service.Load] manually after Configure
// in order to load a new file.
//
// If the store is unsaved upon calling Configure, no attempt is made at saving it.
// Consumers will need to call [Service.Save] manually beforehand.
//
// See [NewWithConfig] for details on configuration.
//
//wails:ignore
func (kvs *KVStoreService) Configure(config *Config) {
	if config != nil {
		// Clone to prevent changes from the outside.
		clone := new(Config)
		*clone = *config
		config = clone
	}

	kvs.lock.Lock()
	defer kvs.lock.Unlock()

	kvs.config = config
	kvs.unsaved = true
}

// Load loads the store from disk.
// If the store is in-memory, i.e. not associated with a file, Load has no effect.
// If the operation fails, a non-nil error is returned
// and the store's content and state at call time are preserved.
func (kvs *KVStoreService) Load() error {
	kvs.lock.Lock()
	defer kvs.lock.Unlock()

	if kvs.config == nil {
		return nil
	}

	bytes, err := os.ReadFile(kvs.config.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	// Init new map because [json.Unmarshal] does not clear the previous one.
	data := make(map[string]any)

	if len(bytes) > 0 {
		if err := json.Unmarshal(bytes, &data); err != nil {
			return err
		}
	}

	kvs.data = data
	kvs.unsaved = false
	return nil
}

// Save saves the store to disk.
// If the store is in-memory, i.e. not associated with a file, Save has no effect.
func (kvs *KVStoreService) Save() error {
	kvs.lock.Lock()
	defer kvs.lock.Unlock()

	if kvs.config == nil {
		return nil
	}

	bytes, err := json.Marshal(kvs.data)
	if err != nil {
		return err
	}

	err = os.WriteFile(kvs.config.Filename, bytes, 0644)
	if err != nil {
		return err
	}

	kvs.unsaved = false
	return nil
}

// Get returns the value for the given key. If key is empty, the entire store is returned.
func (kvs *KVStoreService) Get(key string) any {
	kvs.lock.RLock()
	defer kvs.lock.RUnlock()

	if key == "" {
		return kvs.data
	}

	return kvs.data[key]
}

// Set sets the value for the given key. If AutoSave is true, the store is saved to disk.
func (kvs *KVStoreService) Set(key string, value any) error {
	var autosave bool
	func() {
		kvs.lock.Lock()
		defer kvs.lock.Unlock()

		kvs.data[key] = value
		kvs.unsaved = true

		if kvs.config != nil {
			autosave = kvs.config.AutoSave
		}
	}()

	if autosave {
		return kvs.Save()
	} else {
		return nil
	}
}

// Delete deletes the given key from the store. If AutoSave is true, the store is saved to disk.
func (kvs *KVStoreService) Delete(key string) error {
	var autosave bool
	func() {
		kvs.lock.Lock()
		defer kvs.lock.Unlock()

		delete(kvs.data, key)
		kvs.unsaved = true

		if kvs.config != nil {
			autosave = kvs.config.AutoSave
		}
	}()

	if autosave {
		return kvs.Save()
	} else {
		return nil
	}
}

// Clear deletes all keys from the store. If AutoSave is true, the store is saved to disk.
func (kvs *KVStoreService) Clear() error {
	var autosave bool
	func() {
		kvs.lock.Lock()
		defer kvs.lock.Unlock()

		kvs.data = make(map[string]any)
		kvs.unsaved = true

		if kvs.config != nil {
			autosave = kvs.config.AutoSave
		}
	}()

	if autosave {
		return kvs.Save()
	} else {
		return nil
	}
}

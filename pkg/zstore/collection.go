package zstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/zarlcorp/core/pkg/zcrypto"
)

// Collection is a typed, encrypted key-value collection within a Store.
// Each collection has its own HKDF-derived sub-key so compromising one
// collection does not expose another.
type Collection[V any] struct {
	store *Store
	name  string
	key   []byte
}

// NewCollection returns a typed collection backed by the given store.
// It creates a subdirectory for the collection and derives a sub-key
// via HKDF using the collection name as the info parameter.
func NewCollection[V any](store *Store, name string) (*Collection[V], error) {
	if err := store.fs.MkdirAll(name, 0o700); err != nil {
		return nil, fmt.Errorf("create collection directory: %w", err)
	}

	key, err := zcrypto.ExpandKey(store.masterKey, store.salt, []byte(name))
	if err != nil {
		return nil, fmt.Errorf("derive collection key: %w", err)
	}

	store.subKeys = append(store.subKeys, key)

	return &Collection[V]{store: store, name: name, key: key}, nil
}

// Put encrypts and stores a value under the given id.
func (c *Collection[V]) Put(id string, value V) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value: %w", err)
	}

	ct, err := zcrypto.Encrypt(c.key, data)
	if err != nil {
		return fmt.Errorf("encrypt value: %w", err)
	}

	path := filepath.Join(c.name, id+".enc")
	if err := c.store.fs.WriteFile(path, ct, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	return nil
}

// Get reads, decrypts, and unmarshals a value by id.
// Returns ErrNotFound if the id does not exist.
func (c *Collection[V]) Get(id string) (V, error) {
	var zero V

	path := filepath.Join(c.name, id+".enc")
	ct, err := c.store.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return zero, ErrNotFound
		}
		return zero, fmt.Errorf("read %s: %w", path, err)
	}

	plain, err := zcrypto.Decrypt(c.key, ct)
	if err != nil {
		return zero, fmt.Errorf("decrypt %s: %w", path, err)
	}

	var v V
	if err := json.Unmarshal(plain, &v); err != nil {
		return zero, fmt.Errorf("unmarshal %s: %w", path, err)
	}

	return v, nil
}

// Delete removes a value by id. Returns ErrNotFound if the id does not exist.
func (c *Collection[V]) Delete(id string) error {
	path := filepath.Join(c.name, id+".enc")
	if err := c.store.fs.Remove(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return ErrNotFound
		}
		return fmt.Errorf("remove %s: %w", path, err)
	}
	return nil
}

// List decrypts and returns all values in the collection.
// The order is not guaranteed.
func (c *Collection[V]) List() ([]V, error) {
	var values []V

	err := c.store.fs.WalkDir(c.name, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".enc") {
			return nil
		}

		ct, err := c.store.fs.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		plain, err := zcrypto.Decrypt(c.key, ct)
		if err != nil {
			return fmt.Errorf("decrypt %s: %w", path, err)
		}

		var v V
		if err := json.Unmarshal(plain, &v); err != nil {
			return fmt.Errorf("unmarshal %s: %w", path, err)
		}

		values = append(values, v)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

// Len returns the number of encrypted entries without decrypting them.
func (c *Collection[V]) Len() (int, error) {
	count := 0

	err := c.store.fs.WalkDir(c.name, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".enc") {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

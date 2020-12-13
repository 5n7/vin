package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Cache represents a file cache.
type Cache struct {
	dir string
	m   sync.RWMutex
}

func New(dir string) *Cache {
	return &Cache{dir: dir}
}

func (c *Cache) path(key string) string {
	return filepath.Join(c.dir, key)
}

func (c *Cache) Get(key string) []byte {
	c.m.RLock()
	defer c.m.RUnlock()

	b, err := ioutil.ReadFile(c.path(key))
	if err != nil {
		return nil
	}
	return b
}

func (c *Cache) GetString(key string) string {
	return string(c.Get(key))
}

func (c *Cache) Set(key string, value []byte) error {
	if len(value) == 0 {
		return fmt.Errorf("value is empty")
	}

	c.m.Lock()
	defer c.m.Unlock()

	if err := ioutil.WriteFile(c.path(key), value, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (c *Cache) SetString(key, value string) error {
	return c.Set(key, []byte(value))
}

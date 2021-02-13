package cache

import (
	"fmt"
	"sync"
	"time"
)

//Value struct
type Value struct {
	Object     interface{}
	Expiration int64
}

//Cache struct
type Cache struct {
	expiration time.Duration
	mu         sync.RWMutex
	items      map[string]Value
}

func New(expiration time.Duration) *Cache {
	return &Cache{items: map[string]Value{}, mu: sync.RWMutex{}, expiration: expiration}
}

//Set new value into cache.
//Return error if item already exist.
func (c Cache) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.exist(key) {
		return fmt.Errorf("Item %s already exists", key)
	}

	c.items[key] = Value{Object: value, Expiration: time.Now().Add(c.expiration).UnixNano()}
	return nil
}

//Replace value by key.
//Return error, if value doesn't exist.
func (c Cache) Replace(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.exist(key) {
		return fmt.Errorf("Item %s doesn't exists", key)
	}
	c.items[key] = Value{Object: value, Expiration: time.Now().Add(c.expiration).UnixNano()}
	return nil
}

//Get value by key.
//Return value and indication flag.
//Return nil and false if value expired.
func (c Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.items[key]
	if value.Expiration > 0 {
		if time.Now().UnixNano() > value.Expiration {
			return nil, false
		}
	}
	return value.Object, ok
}

//Delete value from cache.
//Return error if item doesn't exist.
func (c Cache) Del(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.exist(key) {
		return fmt.Errorf("Item %s doesn't exists", key)
	}

	delete(c.items, key)
	return nil
}

//Check is value exist in cache.
//Return isSuccess.
func (c Cache) Exist(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.exist(key)
}

func (c Cache) exist(key string) bool {
	_, ok := c.items[key]
	return ok
}

type Cacher interface {
	Set(key string, value interface{}) error
	Replace(key string, value interface{}) error
	Get(key string) (interface{}, bool)
	Del(key string) error
	Exist(key string) bool
	exist(key string) bool
}

package cache

import (
	"fmt"
	"sync"
)

//Value struct
type Value struct {
	Object interface{}
	//Expire
}

//Cache struct
type Cache struct {
	mu    sync.RWMutex
	items map[string]Value
}

func New() *Cache {
	return &Cache{items: map[string]Value{}, mu: sync.RWMutex{}}
}

//Set new value into cache.
//Return error if item already exist.
func (c Cache) Set(key string, value Value) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.exist(key) {
		return fmt.Errorf("Item %s already exists", key)
	}

	c.items[key] = value
	return nil
}

//Get value by key.
//Return value and isSuccess.
func (c Cache) Get(key string) (Value, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.items[key]
	return value, ok
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
	Set(key string, value Value) error
	Get(key string) (Value, bool)
	Del(key string) error
	Exist(key string) bool
	exist(key string) bool
}

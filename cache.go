package main

import (
	"sync"
	"time"
)

type CacheEntry[T any] struct {
	CreatedAt time.Time
	Val       *T
}

type CacheConfig struct {
	DeleteInterval time.Duration
}

type Cache[T any] struct {
	Ticker *time.Ticker
	Config CacheConfig
	Mutex  sync.Mutex
	Data   map[string]*CacheEntry[T]
}

func getNewCache[T any](config CacheConfig) *Cache[T] {
	c := &Cache[T]{
		Config: config,
		Ticker: time.NewTicker(config.DeleteInterval),
		Data:   make(map[string]*CacheEntry[T]),
	}

	go c.readLoop()

	return c
}

func (c *Cache[T]) Add(key string, value *T) {
	ce := &CacheEntry[T]{
		CreatedAt: time.Now(),
		Val:       value,
	}
	c.Mutex.Lock()
	c.Data[key] = ce
	c.Mutex.Unlock()
}

func (c *Cache[T]) Get(key string) *T {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	d := c.Data[key]
	if d == nil {
		return nil
	}
	return c.Data[key].Val
}

func (c *Cache[T]) readLoop() {
	for {
		<-c.Ticker.C
		c.Mutex.Lock()
		for key, data := range c.Data {
			if time.Since(data.CreatedAt) > c.Config.DeleteInterval {
				delete(c.Data, key)
			}
		}
		c.Mutex.Unlock()
	}
}

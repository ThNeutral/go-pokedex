package main

import (
	"sync"
	"time"
)

type CacheEntry struct {
	CreatedAt time.Time
	Val       *LocationSetType
}

type CacheConfig struct {
	DeleteInterval time.Duration
}

type Cache struct {
	Ticker *time.Ticker
	Config CacheConfig
	Mutex  sync.Mutex
	Data   map[int]*CacheEntry
}

func getNewCache(config CacheConfig) *Cache {
	c := &Cache{
		Config: config,
		Ticker: time.NewTicker(config.DeleteInterval),
		Data:   make(map[int]*CacheEntry),
	}

	go c.readLoop()

	return c
}

func (c *Cache) Add(key int, value *LocationSetType) {
	ce := &CacheEntry{
		CreatedAt: time.Now(),
		Val:       value,
	}
	c.Mutex.Lock()
	c.Data[key] = ce
	c.Mutex.Unlock()
}

func (c *Cache) Get(key int) *LocationSetType {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	d := c.Data[key]
	if d == nil {
		return nil
	}
	return c.Data[key].Val
}

func (c *Cache) readLoop() {
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

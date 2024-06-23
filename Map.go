package main

import (
	"sync"
	"time"
)

type KeyValueCache struct {
	cache      map[string]cacheEntry
	maxSize    int32
	mutex      sync.RWMutex
	expiration time.Duration
}

type cacheEntry struct {
	value      interface{}
	expireTime time.Time
}

func NewCache(MaxSize int32, expiration time.Duration) *KeyValueCache {
	cache := &KeyValueCache{
		cache:      make(map[string]cacheEntry),
		maxSize:    MaxSize,
		expiration: expiration,
	}

	go cache.startCleanupTask()
	return cache
}

func (cache *KeyValueCache) put(key string, value interface{}) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if int32(len(cache.cache)) >= cache.maxSize {
		cache.evict()
	}

	cache.cache[key] = cacheEntry{
		value:      value,
		expireTime: time.Now().Add(cache.expiration),
	}
}

func (cache *KeyValueCache) get(key string) (interface{}, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	entry, found := cache.cache[key]
	if !found || time.Now().After(entry.expireTime) {
		if found {
			delete(cache.cache, key)
		}
		return nil, false
	}
	return entry.value, true
}

func (cache *KeyValueCache) evict() {
	for key := range cache.cache {
		delete(cache.cache, key)
		break
	}
}

func (cache *KeyValueCache) remove(key string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	delete(cache.cache, key)
}

func (cache *KeyValueCache) isContain(key string) bool {
	_, found := cache.cache[key]
	return found
}

func (cache *KeyValueCache) update(key string, value interface{}) int8 {
	if cache.isContain(key) {
		cache.cache[key] = cacheEntry{
			value:      value,
			expireTime: time.Now().Add(cache.expiration),
		}
		return 1
	}
	return -1
}

func (cache *KeyValueCache) startCleanupTask() {
	ticker := time.NewTicker(cache.expiration)
	defer ticker.Stop()

	for {
		<-ticker.C
		cache.cleanup()
	}
}

func (cache *KeyValueCache) cleanup() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	now := time.Now()
	for key, entry := range cache.cache {
		if now.After(entry.expireTime) {
			delete(cache.cache, key)
		}
	}
}

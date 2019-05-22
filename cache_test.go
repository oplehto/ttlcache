package ttlcache

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	cache := &Cache{
		ttl:   time.Second,
		items: map[string]*Item{},
	}

	exists := cache.Get("hello")
	if exists {
		t.Errorf("Expected empty cache to not exist")
	}

	cache.Set("hello")
	exists = cache.Get("hello")
	if !exists {
		t.Errorf("Expected cache to return a positive response")
	}
}

func TestExpiration(t *testing.T) {
	cache := &Cache{
		ttl:   time.Second,
		items: map[string]*Item{},
	}

	cache.Set("x")
	cache.Set("y")
	cache.Set("z")
	cache.startCleanupTimer()

	count := cache.Count()
	if count != 3 {
		t.Errorf("Expected cache to contain 3 items")
	}

	<-time.After(500 * time.Millisecond)
	cache.mutex.Lock()
	cache.items["y"].touch(time.Second)
	item := cache.items["x"]
	cache.mutex.Unlock()
	if item == nil || item.expired() {
		t.Errorf("Expected `x` to not have expired after 200ms")
	}

	<-time.After(time.Second)
	cache.mutex.RLock()
	item = cache.items["x"]
	if item != nil {
		t.Errorf("Expected `x` to have expired")
	}
	item = cache.items["z"]
	if item != nil {
		t.Errorf("Expected `z` to have expired")
	}
	item = cache.items["y"]
	if item == nil {
		t.Errorf("Expected `y` to not have expired")
	}
	cache.mutex.RUnlock()

	count = cache.Count()
	if count != 1 {
		t.Errorf("Expected cache to contain 1 item")
	}

	<-time.After(600 * time.Millisecond)
	cache.mutex.RLock()
	item = cache.items["y"]
	if item != nil {
		t.Errorf("Expected `y` to have expired")
	}
	cache.mutex.RUnlock()

	count = cache.Count()
	if count != 0 {
		t.Errorf("Expected cache to be empty")
	}
}

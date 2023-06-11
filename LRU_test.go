package LRU

import (
	"fmt"
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {
	cache := NewLRUCache(3)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	value, ok := cache.Get("key1")
	if !ok || value != "value1" {
		t.Error("Expected value1, got", value)
	}


	cache.Add("key4", "value4")


	_, ok = cache.Get("key2")
	if ok {
		t.Error("Expected key2 to be evicted, but it was found")
	}


	cache.AddWithTTL("key5", "value5", time.Second*2)


	time.Sleep(time.Second * 3)

	_, ok = cache.Get("key5")
	if ok {
		t.Error("Expected key5 to be evicted, but it was found")
	}
}

func TestLRUCacheConcurrency(t *testing.T) {
	cache := NewLRUCache(10)

	for i := 0; i < 100; i++ {
		go func(index int) {
			key := fmt.Sprintf("key%d", index)
			value := fmt.Sprintf("value%d", index)

			cache.Add(key, value)

			v, ok := cache.Get(key)
			if !ok || v != value {
				t.Errorf("Expected %s, got %v", value, v)
			}
		}(i)
	}
}

func TestLRUCacheClear(t *testing.T) {
	cache := NewLRUCache(3)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	cache.Clear()

	if len(cache.cache) != 0 || cache.list.Len() != 0 {
		t.Error("Cache clear failed")
	}
}

func TestLRUCacheRemove(t *testing.T) {
	cache := NewLRUCache(3)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	cache.Remove("key2")

	_, ok := cache.Get("key2")
	if ok {
		t.Error("Expected key2 to be removed, but it was found")
	}
}

func TestLRUCacheCapacity(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	_, ok := cache.Get("key1")
	if ok {
		t.Error("Expected key1 to be evicted, but it was found")
	}

	cache.cap = 3

	cache.Add("key4", "value4")

	_, ok = cache.Get("key2")
	if !ok {
		t.Error("Expected key2 to be found, but it was evicted")
	}
}

func TestLRUCacheTTL(t *testing.T) {
	cache := NewLRUCache(5)


	cache.AddWithTTL("key1", "value1", time.Second)

	time.Sleep(time.Millisecond * 1500)

	_, ok := cache.Get("key1")
	if ok {
		t.Error("Expected key1 to be evicted, but it was found")
	}
}

func TestLRUCacheTTLUpdate(t *testing.T) {
	cache := NewLRUCache(5)

	cache.AddWithTTL("key1", "value1", time.Second)

	time.Sleep(time.Millisecond * 500)

	cache.AddWithTTL("key1", "value1-updated", time.Second*2)

	time.Sleep(time.Millisecond * 1500)


	value, ok := cache.Get("key1")
	if !ok || value != "value1-updated" {
		t.Error("Expected value1-updated, got", value)
	}
}

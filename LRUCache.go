package LRU

import (
	"container/list"
	"sync"
	"time"
)

type ICache interface {
	Cap() int
	Clear()
	Add(key, value interface{})
	AddWithTTL(key, value interface{}, ttl time.Duration)
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
}

type cacheElement struct {
	key   interface{}
	value interface{}
	exp   time.Time
	dur   time.Duration
}

// Используем List чтобы была очередь
// Мапа чтобы хранить по ключам

type LRUCache struct {
	cap   int
	cache map[interface{}]*list.Element
	list  *list.List
	mutex sync.Mutex
}

func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		cap:   cap,
		cache: make(map[interface{}]*list.Element),
		list:  list.New(),
	}
}

func (c *LRUCache) Cap() int {
	return c.cap
}

// Просто создаем новые пустые структуры

func (c *LRUCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = make(map[interface{}]*list.Element)
	c.list = list.New()
}

func (c *LRUCache) AddWithTTL(key, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Если значение уже в кеше, то обновляем значение, ттл и двигаем вперед в очереди

	if elem, ok := c.cache[key]; ok {
		entry := elem.Value.(*cacheElement)
		entry.value = value
		entry.dur = ttl
		entry.exp = time.Now().Add(ttl)
		c.list.MoveToFront(elem)
		return
	}

	entry := &cacheElement{
		key:   key,
		value: value,
		exp:   time.Now().Add(ttl),
		dur:   ttl,
	}
	el := c.list.PushFront(entry)
	c.cache[key] = el

	// Если превысили капасити, то убираем Least Recently Used
	if c.list.Len() > c.cap {
		last := c.list.Back()
		if last != nil {
			delete(c.cache, (last.Value.(*cacheElement)).key)
			c.list.Remove(last)
		}
	}
}

func (c *LRUCache) Add(key, value interface{}) {
	c.AddWithTTL(key, value, time.Duration(0))
}

func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	/*
	* Если нет в кеше, то возвращаем (nil, false)
	* Если если есть, но прошел срок годности, то удаляем из кеша и говорим что не нашли
	* Иначе передвигаем элемент к которому запросили доступ в начало очереди
	* Если у элемента нет срока годности, то считаем что он там навсегда, пока не привысят кап
	 */

	if elem, ok := c.cache[key]; ok {
		entry := elem.Value.(*cacheElement)
		if entry.dur != time.Duration(0) && entry.exp.Before(time.Now()) {
			delete(c.cache, entry.key)
			c.list.Remove(elem)
			return nil, false
		}
		c.list.MoveToFront(elem)
		return entry.value, true
	}
	return nil, false
}

// Удаляем из мапы и из очереди

func (c *LRUCache) Remove(key interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if elem, ok := c.cache[key]; ok {
		delete(c.cache, key)
		c.list.Remove(elem)
	}

}

package lru

import (
	"container/list"
	"time"
)

// Cache LRU,不支持并发安全
type Cache struct {
	maxBytes  int64
	nbytes    int64                         // 已使用的内存
	ll        *list.List                    // 双向链表
	cache     map[string]*list.Element      // 链表节点
	OnEvicted func(key string, value Value) // 记录被删除时的回调函数
}

type entry struct {
	key   string
	value Value
	// expireAt <= 0 代表不过期
	expireAt int64
}

func (e *entry) expired() bool {
	return e.expireAt > 0 && time.Now().UnixNano() >= e.expireAt
}

type Value interface {
	Len() int
}

// New Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look-ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		e := ele.Value.(*entry)
		if e.expired() {
			c.RemoveElement(ele)
			return nil, false
		}
		c.ll.MoveToFront(ele)
		return e.value, true
	}
	return nil, false
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value, ttl time.Duration) {
	var expireAt int64
	if ttl > 0 {
		expireAt = time.Now().Add(ttl).UnixNano()
	}
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		e := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(e.value.Len())
		e.expireAt = expireAt
		e.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value, expireAt})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}

func (c *Cache) RemoveElement(ele *list.Element) {
	c.ll.Remove(ele)
}

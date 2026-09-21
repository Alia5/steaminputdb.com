package memcache

import (
	"slices"
	"sync"
	"time"
)

type Cache struct {
	m           *sync.Map
	validity    time.Duration
	maxItems    int
	cleanupLock sync.Mutex
}

func New(validity time.Duration, maxItems int) *Cache {
	return &Cache{
		m:        &sync.Map{},
		validity: validity,
		maxItems: maxItems,
	}
}

func Get[T any](c *Cache, key string) (T, bool) {
	v, ok := c.Get(key)
	if !ok {
		var zero T
		return zero, false
	}
	return v.(T), true
}

func (c *Cache) Get(key string) (any, bool) {
	v, ok := c.m.Load(key)
	go c.cleanup()

	if !ok {
		return nil, false
	}
	item := v.(item)
	if time.Since(item.Added) > c.validity {
		c.m.Delete(key)
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) cleanup() {
	if !c.cleanupLock.TryLock() {
		return
	}
	defer c.cleanupLock.Unlock()

	var items []item
	c.m.Range(func(key, value any) bool {
		it := value.(item)
		if time.Since(it.Added) > c.validity {
			c.m.Delete(key)
		} else {
			items = append(items, it)
		}
		return true
	})

	if len(items) > c.maxItems {
		slices.SortFunc(items, func(a, b item) int {
			if a.Added.Before(b.Added) {
				return -1
			}
			if a.Added.After(b.Added) {
				return 1
			}
			return 0
		})
		for i := range len(items) - c.maxItems {
			c.m.Delete(items[i].Key)
		}
	}
}

func (c *Cache) Store(key string, value any) {
	c.m.Store(key, item{
		Key:   key,
		Value: value,
		Added: time.Now(),
	})
	go c.cleanup()
}

type item struct {
	Key   string
	Value any
	Added time.Time
}

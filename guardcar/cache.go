package guardcar

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// 简易带过期的缓存

var (
	// ErrKeyNotFound gets returned when a specific key couldn't be found
	ErrKeyNotFound = errors.New("Key not found in cache")
	// ErrKeyNotFoundOrLoadable gets returned when a specific key couldn't be
	// found and loading via the data-loader callback also failed
	ErrKeyNotFoundOrLoadable = errors.New("Key not found and could not be loaded into cache")
)

var (
	mutex sync.RWMutex
	cache *Cache
)

// 缓存结构
type Cache struct {
	sync.RWMutex
	// All cached items.
	items map[string]*CacheItem

	// Timer responsible for triggering cleanup.
	cleanupTimer *time.Timer
	// Current timer duration.
	cleanupInterval time.Duration

	// Callback method triggered before deleting an item from the cache.
	aboutToDeleteItem func(item *CacheItem)

	// The logger used for this table.
	logger *log.Logger
}

// 创建
func InitCache(fun func(item *CacheItem)) *Cache {
	mutex.RLock()
	if cache == nil {
		cache = &Cache{
			items:             make(map[string]*CacheItem),
			aboutToDeleteItem: fun,
		}
	}
	mutex.RUnlock()
	return cache
}

// 添加
func (cache *Cache) Add(key string, lifeSpan time.Duration, data interface{}) *CacheItem {
	item := NewCacheItem(key, lifeSpan, data)

	// Add item to cache.
	cache.Lock()
	cache.addInternal(item)

	return item
}

// 获取值
func (cache *Cache) Value(key string) (*CacheItem, error) {
	cache.RLock()
	r, ok := cache.items[key]
	cache.RUnlock()

	if ok {
		// Update access counter and timestamp.
		fmt.Println("更新时间key:", key)
		r.KeepAlive()
		return r, nil
	}

	return nil, ErrKeyNotFound
}

// 删除
func (cache *Cache) Delete(key string) (*CacheItem, error) {
	cache.Lock()
	defer cache.Unlock()

	return cache.deleteInternal(key)
}

// 是否存在
func (cache *Cache) Exists(key string) bool {
	cache.RLock()
	defer cache.RUnlock()
	_, ok := cache.items[key]

	return ok
}

// Count returns how many items are currently stored in the cache.
func (cache *Cache) Count() int {
	cache.RLock()
	defer cache.RUnlock()
	return len(cache.items)
}

func (cache *Cache) deleteInternal(key string) (*CacheItem, error) {
	r, ok := cache.items[key]
	if !ok {
		return nil, ErrKeyNotFound
	}

	// Cache value so we don't keep blocking the mutex.
	aboutToDeleteItem := cache.aboutToDeleteItem
	cache.Unlock()

	// Trigger callbacks before deleting an item from cache.
	if aboutToDeleteItem != nil {
		aboutToDeleteItem(r)
	}

	cache.Lock()
	cache.log("Deleting item with key", key, "created on", r.createdOn)
	delete(cache.items, key)

	return r, nil
}

func (cache *Cache) addInternal(item *CacheItem) {
	// Careful: do not run this method unless the table-mutex is locked!
	// It will unlock it for the caller before running the callbacks and checks
	cache.log("Adding item with key", item.key, "and lifespan of", item.lifeSpan)
	cache.items[item.key] = item

	// Cache values so we don't keep blocking the mutex.
	expDur := cache.cleanupInterval
	cache.Unlock()

	// If we haven't set up any expiration check timer or found a more imminent item.
	if item.lifeSpan > 0 && (expDur == 0 || item.lifeSpan < expDur) {
		cache.expirationCheck()
	}
}

// 过期检查
func (cache *Cache) expirationCheck() {
	cache.Lock()
	if cache.cleanupTimer != nil {
		cache.cleanupTimer.Stop()
	}
	if cache.cleanupInterval > 0 {
		cache.log("Expiration check triggered after", cache.cleanupInterval)
	} else {
		cache.log("Expiration check installed for table")
	}

	// To be more accurate with timers, we would need to update 'now' on every
	// loop iteration. Not sure it's really efficient though.
	now := time.Now()
	smallestDuration := 0 * time.Second
	for key, item := range cache.items {
		// Cache values so we don't keep blocking the mutex.
		item.RLock()
		lifeSpan := item.lifeSpan
		accessedOn := item.accessedOn
		item.RUnlock()

		if lifeSpan == 0 {
			continue
		}
		if now.Sub(accessedOn) >= lifeSpan {
			// Item has excessed its lifespan.
			cache.deleteInternal(key)
		} else {
			// Find the item chronologically closest to its end-of-lifespan.
			if smallestDuration == 0 || lifeSpan-now.Sub(accessedOn) < smallestDuration {
				smallestDuration = lifeSpan - now.Sub(accessedOn)
			}
		}
	}

	// Setup the interval for the next cleanup run.
	cache.cleanupInterval = smallestDuration
	if smallestDuration > 0 {
		cache.cleanupTimer = time.AfterFunc(smallestDuration, func() {
			go cache.expirationCheck()
		})
	}
	cache.Unlock()

}

// Internal logging method for convenience.
func (cache *Cache) log(v ...interface{}) {
	if cache.logger == nil {
		return
	}

	cache.logger.Println(v...)
}

// 缓存条目
type CacheItem struct {
	sync.RWMutex
	key        string
	data       interface{}
	lifeSpan   time.Duration
	createdOn  time.Time
	accessedOn time.Time
}

func NewCacheItem(key string, lifeSpan time.Duration, data interface{}) *CacheItem {
	t := time.Now()
	return &CacheItem{
		key:        key,
		lifeSpan:   lifeSpan,
		createdOn:  t,
		accessedOn: t,
		data:       data,
	}
}

// 保活
func (item *CacheItem) KeepAlive() {
	item.Lock()
	defer item.Unlock()
	item.accessedOn = time.Now()
}

// 返回值
func (item *CacheItem) Data() interface{} {
	// immutable
	return item.data
}

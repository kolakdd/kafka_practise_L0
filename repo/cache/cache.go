package cache

import "sync"

type CacheRepo interface {
	Get(k int, withMX bool) ([]byte, bool)
	Set(k int, v []byte)
	GetMaxLen() int
	Debug() string
}

type cacheRepo struct {
	mu     sync.Mutex
	maxLen int
	map_   map[int]*node
	ll     list
}

// LRU Cache
func NewCacheRepo(maxLen int) CacheRepo {
	if maxLen == 0 {
		maxLen = 100
	}
	return &cacheRepo{
		maxLen: maxLen,
		map_:   map[int]*node{},
		ll:     NewList(),
	}
}

func (c *cacheRepo) Get(k int, withMX bool) ([]byte, bool) {
	if withMX {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	node, ok := c.map_[k]
	if !ok {
		return nil, false
	}

	if err := node.Destructor(); err != nil {
		panic(err)
	}
	c.ll.len--

	c.ll.AppendToHead(node.val)
	return node.val, true
}

func (c *cacheRepo) Set(k int, v []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Get(k, false); ok {
		return
	}
	if c.ll.len >= c.maxLen {
		c.ll.RemoveLast()
	}
	newNode := c.ll.AppendToHead(v)
	c.map_[k] = newNode
}

func (c *cacheRepo) Debug() string {
	return c.ll.Represent()
}

func (c *cacheRepo) GetMaxLen() int {
	return c.maxLen
}

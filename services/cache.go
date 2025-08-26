package services

import (
	"encoding/json"
	"kafkapractisel0/repo"
	"kafkapractisel0/repo/cache"
	"log"
	"sync"
	"sync/atomic"
)

type CacheService interface {
	UpdateCacheNewest(count int)
}

type cacheService struct {
	c cache.CacheRepo
	o repo.OrderRepo
}

func NewCacheService(c cache.CacheRepo, o repo.OrderRepo) CacheService {
	return &cacheService{c, o}
}

// Обновить кеш новейшими записями
func (c *cacheService) UpdateCacheNewest(count int) {
	if count > c.c.GetMaxLen() {
		count = c.c.GetMaxLen()
	}
	var errCounter atomic.Int32
	var wg sync.WaitGroup
	wg.Add(count)
	for i := range count {
		go func(offset int) {
			defer wg.Done()
			order, err := c.o.SelectNewestWithOffset(offset)
			if err != nil {
				log.Println(err)
				errCounter.Add(1)
				return
			}
			jData, errM := json.Marshal(order)
			if errM != nil {
				log.Println(err)
				errCounter.Add(1)
				return
			}
			c.c.Set(order.UID, jData)
		}(i)
	}
	wg.Wait()
	log.Printf("[Cache Restored], total = %v , errors counter = %v", count, errCounter.Load())
}

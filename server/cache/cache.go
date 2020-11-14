/**
This caching system was made simply for fun and for
golang practice.
Here we are only interested in storing a certain result
in computer memory. We also implemented a simple eviction system .
*/
package cache

import (
	"time"
)

/**
entry is a cache element in the cache
this was built to give support to exp time
*/
type entry struct {
	obj interface{}
	exp int64 //Expiration time in linux epoch
}

/**
outEntry gives support to retrive elements from
the cache
*/
type outEntry struct {
	respchan chan *interface{}
	suc      bool
	reqKey   string
}

type inEntry struct {
	e   interface{}
	key string
}

/**
isExpired checks if a certain
entry has expired
*/
func (e entry) isExpired() bool {
	if e.exp == 0 {
		return false
	}

	return time.Now().Unix() > e.exp
}

type Cache interface {
	Set(k string, x interface{})
	Get(k string) (interface{}, bool)
}

type simpleCache struct {
	expTime time.Duration
	entries map[string]entry
	outchan chan outEntry
	inchan  chan inEntry
}

/**
Create a new simpleCache with a given evict time
*/
func New(expTime time.Duration) Cache {
	c := simpleCache{
		expTime: expTime,
		entries: make(map[string]entry),
		outchan: make(chan outEntry),
		inchan:  make(chan inEntry),
	}

	go c.cacheHandler()

	return &c
}

func returnExp(c *simpleCache) int64{
	if c.expTime == 0 {
		return 0
	}
	return  time.Now().Add(c.expTime).Unix()
}

func (c *simpleCache) cacheHandler() {
	wait := time.Millisecond * 100
	to := time.After(wait)
	for {
		select {
		case <-to:
			c.evictEntries()
			to = time.After(wait)
		case el := <-c.inchan:
			c.entries[el.key] = entry{
				obj: el.e,
				exp: returnExp(c),
			}
		case el := <-c.outchan:
			i, ok := c.entries[el.reqKey]
			if ok {
				el.respchan <- &i.obj
			} else {
				el.respchan <- nil
			}
		}
	}
}

func (c *simpleCache) evictEntries() {
	for k, v := range c.entries {
		if v.isExpired() {
			delete(c.entries, k)
		}
	}
}

/**
Get retrives a element from the Cache
*/
func (c *simpleCache) Get(k string) (interface{}, bool) {
	ch := make(chan *interface{})
	c.outchan <- outEntry{
		respchan: ch,
		suc:      false,
		reqKey:   k,
	}
	el := <-ch
	if el != nil {
		return *el, true
	}
	return struct{}{}, false
}

/**
Set puts a element in the Cache
*/
func (c *simpleCache) Set(k string, x interface{}) {
	c.inchan <- inEntry{
		e:   x,
		key: k,
	}
}

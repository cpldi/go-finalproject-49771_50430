package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSimpleCache(t *testing.T) {
	c := New(0)
	c.Set("3",3)

	res , ok := c.Get("3")
	if !ok {
		t.Fatalf("Element should be presented\n")
	}
	if res != 3 {
		t.Fatalf("Expecting %v  got %v d\n",3,res)
	}


	res , ok = c.Get("5")
	if ok {
		t.Fatalf("Element should not be presented\n")
	}

}
func TestSimpleCacheEvict(t *testing.T) {
	exptime := 100 * time.Millisecond
	c := New(exptime)
	c.Set("3",3)

	time.Sleep(exptime * 10)

	_ , ok := c.Get("3")
	if ok {
		t.Fatalf("Element should not be presented\n")
	}

}

func TestConcurrentInsertion(t *testing.T) {
	const N = 1000
	var wg sync.WaitGroup
	c := New(0)

	wg.Add(N)

	for i := 0 ; i < N ; i++{
		go func(i int) {
			c.Set(strconv.Itoa(i),i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	wg.Add(N)
	for i := 0 ; i < N ; i++{
		go func(i int) {
			_ , ok := c.Get(strconv.Itoa(i))
			if !ok {
				t.Fatalf("Element should be presented\n")
			}
			wg.Done()
		}(i)
	}

}
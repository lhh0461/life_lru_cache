package life_lru

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUseCase_Get(t *testing.T) {
	cache := NewLRUCache[int, int](3)

	var v int
	var ok bool

	cache.Set(1, 1, time.Hour)
	v, ok = cache.Get(1)
	assert.Equal(t, 1, v)
	assert.Equal(t, true, ok)

	v, ok = cache.Get(2)
	assert.Equal(t, 0, v)
	assert.Equal(t, false, ok)

	cache.Set(1, 2, time.Hour)
	v, ok = cache.Get(1)
	assert.Equal(t, 2, v)
	assert.Equal(t, true, ok)
}

func TestUseCase_LRU(t *testing.T) {
	cache := NewLRUCache[int, string](3)
	cache.Set(1, "a", time.Hour)
	cache.Set(2, "b", time.Hour)
	cache.Set(3, "c", time.Hour)
	cache.Set(4, "d", time.Hour)

	var v string
	var ok bool
	v, ok = cache.Get(1)
	assert.Equal(t, "", v)
	assert.Equal(t, false, ok)

	v, ok = cache.Get(2)
	assert.Equal(t, "b", v)
	assert.Equal(t, true, ok)

	v, ok = cache.Get(3)
	assert.Equal(t, "c", v)
	assert.Equal(t, true, ok)

	v, ok = cache.Get(4)
	assert.Equal(t, "d", v)
	assert.Equal(t, true, ok)

	v, ok = cache.Get(2)
	assert.Equal(t, "b", v)
	assert.Equal(t, true, ok)

	//del c
	cache.Set(5, "e", time.Hour)

	v, ok = cache.Get(3)
	assert.Equal(t, "", v)
	assert.Equal(t, false, ok)
}

func TestUseCase_Timeout(t *testing.T) {
	cache := NewLRUCache[int, string](3)
	cache.Set(2, "b", time.Hour)
	cache.Set(3, "c", time.Second)

	time.Sleep(time.Second * 2)

	cache.Set(4, "d", time.Hour)

	var v string
	var ok bool
	v, ok = cache.Get(1)
	assert.Equal(t, "a", v)
	assert.Equal(t, true, ok)

	v, ok = cache.Get(2)
	assert.Equal(t, "b", v)
	assert.Equal(t, true, ok)

	v, ok = cache.Get(3)
	assert.Equal(t, "", v)
	assert.Equal(t, false, ok)

	v, ok = cache.Get(4)
	assert.Equal(t, "d", v)
	assert.Equal(t, true, ok)
}

func TestUseCase_Struct(t *testing.T) {
	type TestType struct {
		age  uint32
		sex  bool
		name string
	}
	cache := NewLRUCache[int, TestType](3)
	cache.Set(1, TestType{1, false, "andy"}, time.Hour)
	cache.Set(2, TestType{2, false, "Tom"}, time.Hour)
	cache.Set(3, TestType{2, true, "July"}, time.Hour)
	cache.Set(4, TestType{3, true, "Lisa"}, time.Hour)

	var v TestType
	var ok bool
	v, ok = cache.Get(1)
	assert.Equal(t, TestType{}, v)
	assert.Equal(t, false, ok)

	v, ok = cache.Get(2)
	assert.Equal(t, TestType{2, false, "Tom"}, v)
	assert.Equal(t, true, ok)

	v, ok = cache.Get(3)
	assert.Equal(t, TestType{2, true, "July"}, v)
	assert.Equal(t, true, ok)

	v, ok = cache.Get(4)
	assert.Equal(t, TestType{3, true, "Lisa"}, v)
	assert.Equal(t, true, ok)
}

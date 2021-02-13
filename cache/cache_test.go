package cache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCacheSet_WithStringValue_ShouldSucceed(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)

	c.Set("test", "testValue")

	res := c.items["test"]
	assert.Equal(t, "testValue", res.Object.(string))
}

func TestCacheSet_WithIntValue_ShouldSucceed(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)
	c.Set("test", 42)
	res := c.items["test"]
	assert.Equal(t, 42, res.Object.(int))
}

func TestCacheSet_WithFloatValue_ShouldSucceed(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)
	c.Set("test", 42.42)
	res := c.items["test"]
	assert.Equal(t, 42.42, res.Object.(float64))
}

func TestCacheSet_WithBoolValue_ShouldSucceed(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)
	c.Set("test", true)
	res := c.items["test"]
	assert.Equal(t, true, res.Object.(bool))
}

func TestCacheSet_WithDoubleFloatValue_ShouldReturnError(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)
	err1 := c.Set("test", 42.42)
	err2 := c.Set("test", 41.42)

	assert.NoError(t, err1)
	assert.Error(t, err2)
}

func TestCacheSet_WithDoubleBoolValue_ShouldReturnError(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)
	err1 := c.Set("test", true)
	err2 := c.Set("test", false)

	assert.NoError(t, err1)
	assert.Error(t, err2)
}

func TestCacheSet_WithDoubleStringValue_ShouldReturnError(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)
	err1 := c.Set("test", "testValue")
	err2 := c.Set("test", "testValue2")

	assert.NoError(t, err1)
	assert.Error(t, err2)
}

func TestCacheSet_WithDoubleIntValue_ShouldReturnError(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)

	err1 := c.Set("test", 42)
	err2 := c.Set("test", 41)

	assert.NoError(t, err1)
	assert.Error(t, err2)
}

func TestCacheGet_WhenValueExpired_ShouldReturnNil(t *testing.T) {
	dur, err := time.ParseDuration("1s")
	if err != nil {
		return
	}
	c := New(dur)
	err = c.Set("test", 42)

	time.Sleep(2 * time.Second)

	var res int
	v, ok := c.Get("test")
	if ok {
		res = v.(int)
	}

	assert.NoError(t, err)
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, v)
	assert.Equal(t, 0, res)
}

func TestCacheReplace_WithValidData_ShouldSucceed(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)
	err1 := c.Set("test", 42)

	err2 := c.Replace("test", "Hello world!")

	var res string
	v, ok := c.Get("test")
	if ok {
		res = v.(string)
	}

	assert.NoError(t, err)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, true, ok)
	assert.Equal(t, "Hello world!", res)
}

func TestCacheReplace_WithoutInitialValue_ShouldReturnError(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)

	err1 := c.Replace("test", "Hello world!")

	assert.NoError(t, err)
	assert.Error(t, err1)
}

func TestCacheGet_WithConcurrent_ShouldSucceed(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)

	c.Set("test", 42)
	go c.Set("test1", 41)

	time.Sleep(1 * time.Second)

	ch := make(chan int)

	var res int
	if v, ok := c.Get("test"); ok {
		res = v.(int)
	}

	var res1 int
	go func() {
		if v, ok := c.Get("test1"); ok {
			ch <- v.(int)
		}
	}()
	res1 = <-ch

	assert.Equal(t, 42, res)
	assert.Equal(t, 41, res1)
}

func TestCacheExist_WithValidData_ShouldSucceed(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)

	err = c.Set("Test", "test")
	exist := c.Exist("Test")

	assert.NoError(t, err)
	assert.Equal(t, true, exist)
}

func TestCacheExist_WithoutValue_ShouldReturnFalse(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)

	exist := c.Exist("Test")

	assert.Equal(t, false, exist)
}

func TestCacheDelete_WithValidData_ShouldSucceed(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)

	err = c.Set("Test", "test")
	exist := c.Exist("Test")
	err1 := c.Del("Test")
	res := c.items["Test"]

	assert.NoError(t, err)
	assert.Equal(t, true, exist)
	assert.NoError(t, err1)
	assert.Nil(t, res.Object)
}

func TestCacheDelete_WithoutData_ShouldReturnError(t *testing.T) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)

	err1 := c.Del("Test")

	assert.Error(t, err1)
}

func BenchmarkCache_Set(b *testing.B) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := c.Set(fmt.Sprint(i), i); err != nil {
			return
		}
	}
}

func BenchmarkCache_Get(b *testing.B) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)
	for i := 0; i < b.N; i++ {
		if err := c.Set(fmt.Sprint(i), i); err != nil {
			return
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := c.Get(fmt.Sprint(i)); !ok {
			return
		}
	}
}

func BenchmarkCache_GetInParallels(b *testing.B) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)
	b.SetParallelism(30)
	for i := 0; i < b.N; i++ {
		if err := c.Set(fmt.Sprint(i), i); err != nil {
			return
		}
	}

	i := 0
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, ok := c.Get(fmt.Sprint(i)); !ok {
				return
			}
			i++
		}
	})
}

func BenchmarkCache_ExistInParallels(b *testing.B) {
	dur, err := time.ParseDuration("5m")
	if err != nil {
		return
	}
	c := New(dur)
	b.SetParallelism(30)
	for i := 0; i < b.N; i++ {
		if err := c.Set(fmt.Sprint(i), i); err != nil {
			return
		}
	}

	i := 0
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if ok := c.Exist(fmt.Sprint(i)); !ok {
				return
			}
			i++
		}
	})
}

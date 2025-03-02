package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
	c := &Cache{}

	// Test setting and getting a value
	c.Set("key1", "value1", 1*time.Minute)
	value, found := c.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value1", value)

	// Test getting a non-existent value
	value, found = c.Get("nonexistent")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestCache_Expiration(t *testing.T) {
	c := &Cache{}

	// Test setting a value with expiration
	c.Set("key2", "value2", 1*time.Second)
	time.Sleep(2 * time.Second)
	value, found := c.Get("key2")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestCache_Delete(t *testing.T) {
	c := &Cache{}

	// Test deleting a value
	c.Set("key3", "value3", 1*time.Minute)
	c.Delete("key3")
	value, found := c.Get("key3")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

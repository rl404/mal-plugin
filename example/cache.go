package example

import (
	"time"

	"github.com/rl404/mal-plugin/cache/redis"
)

// Cache is Cacher interface usage example.
func Cache() {
	// Init client.
	client, err := redis.New("localhost:6379", "", time.Hour)
	if err != nil {
		panic(err)
	}

	// Don't forget to close.
	defer client.Close()

	// Sample data. Can be any type.
	data := []string{"a", "b", "c"}

	// Save to cache. Data will be encoded first.
	err = client.Set("key", data)
	if err != nil {
		panic(err)
	}

	// Create a new or use existing variable.
	// Data type should be the same as when saving to cache.
	var newData []string

	// Get data from cache. Data will be decoded to inputted
	// variable. Don't forget to use pointer.
	err = client.Get("key", &newData)
	if err != nil {
		panic(err)
	}

	// Delete data from cache.
	err = client.Delete("key")
	if err != nil {
		panic(err)
	}
}

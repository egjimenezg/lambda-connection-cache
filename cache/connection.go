package cache

import (
  "sync"
)

var (
  once      sync.Once
  connCache *ConnectionCache
)

// Function that will be executed to get the connection pool
type ConnectionFunction func(connectionString string) (interface{}, error)

// Result of connection 
type connResult struct {
  value   interface{}
  err     error
}

// ConnectionCache saves the created connection pools 
type ConnectionCache struct {
  cache   map[string]connResult
}

func New() *ConnectionCache {
  once.Do(func(){
    connCache = &ConnectionCache{
      cache: make(map[string]connResult),
    }
  })

  return connCache
}


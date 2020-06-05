package cache

import (
  "sync"
)

var (
  once        sync.Once
  connCache   *ConnectionCache
)

// Function that will be executed to get the connection pool
type ConnectionFunction func(connectionString string) (interface{}, error)

// Entry with channel
type cacheEntry struct {
  entry   connectionPool
  ready   chan struct{}
}

// Result of connection 
type connectionPool struct {
  connection    interface{}
  err           error
}

// ConnectionCache saves the created connection pools 
type ConnectionCache struct {
  cache   map[string]*cacheEntry
  mutex   sync.Mutex
}

// New creates a connection cache only once
func New() *ConnectionCache {
  once.Do(func(){
    connCache = &ConnectionCache{
      cache: make(map[string]*cacheEntry),
    }
  })

  return connCache
}

// Get obtains or create a connection pool
func (connectionCache *ConnectionCache) Get(connectionString string, connectionFunction ConnectionFunction) (interface{}, error) {
  connectionCache.mutex.Lock()
  connEntry := connectionCache.cache[connectionString]

  if connEntry == nil {
    connEntry = &cacheEntry{ready: make(chan struct{})}

    connectionCache.cache[connectionString] = connEntry
    connectionCache.mutex.Unlock()

    connEntry.entry.connection, connEntry.entry.err = connectionFunction(connectionString)

    //Broadcast to the goroutines that the function was executed
    close(connEntry.ready)
  } else {
    connectionCache.mutex.Unlock()
    <-connEntry.ready
  }

  return connEntry.entry.connection, connEntry.entry.err
}

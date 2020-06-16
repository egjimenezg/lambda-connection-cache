package cache

import (
  "sync"
)

var (
  once        sync.Once
  connCache   *ConnectionCache
)

// Function that will be executed to get the connection handler
type ConnectionFunction func(connectionString string) (interface{}, error)

// Entry with channel
type cacheEntry struct {
  entry   connectionHandler
  ready   chan struct{}
}

// Result of connection 
type connectionHandler struct {
  connection    interface{}
  err           error
}

// ConnectionCache saves the created connection handlers
type ConnectionCache struct {
  cache   map[string]*cacheEntry
  mutex   sync.Mutex
}

// New creates a connection handler cache only once
func New() *ConnectionCache {
  once.Do(func(){
    connCache = &ConnectionCache{
      cache: make(map[string]*cacheEntry),
    }
  })

  return connCache
}

// Get obtains or create a connection handler
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
    //Block until function has been executed
    <-connEntry.ready
  }

  return connEntry.entry.connection, connEntry.entry.err
}

func (connectionCache *ConnectionCache) Size() int {
  return len(connectionCache.cache)
}

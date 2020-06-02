package cache_test

import (
  "github.com/egjimenezg/lambda-connection-cache/cache"
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestCreateNewConnectionCacheOnlyOnce(t *testing.T){
  connectionCache := cache.New()
  assert.Equal(t, connectionCache, cache.New())
}

package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestCreateNewConnectionCache(t *testing.T){
  connectionCache := New()
  assert.Equal(t, len(connectionCache), 0)
}

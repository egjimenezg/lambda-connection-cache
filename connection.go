package main

type ConnectionCache map[string]interface{}

func New() ConnectionCache {
  connectionCache := make(ConnectionCache)
  return connectionCache
}

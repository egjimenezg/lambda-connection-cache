package cache_test

import (
  "database/sql"
  cache "github.com/egjimenezg/lambda-connection-cache"
  "github.com/stretchr/testify/assert"
  "github.com/DATA-DOG/go-sqlmock"
  "reflect"
  "sync"
  "testing"
)

func TestCreateNewConnectionCacheOnlyOnce(t *testing.T){
  connectionCache := cache.New()
  assert.Equal(t, connectionCache, cache.New())
}

func TestCreateNewConnectionManagerAndSavingItInCacheMap(t *testing.T) {
  connectionCache := cache.New()
  connManagerA, err := connectionCache.Get("user:password@/dbname", createConnectionManager)

  if err != nil {
    t.Fatalf("An error '%s' has ocurred", err)
  }

  connManagerB, err := connectionCache.Get("user:password@/dbname", createConnectionManager)

  assert.True(t, reflect.DeepEqual(connManagerA.(*sql.DB), connManagerB.(*sql.DB)))
}

func TestCreateOnlyOneConnectionManager(t *testing.T){
  requestsNumber := 1000
  connectionCache := cache.New()

  var waitGroup sync.WaitGroup

  for i := 0; i<requestsNumber; i++ {
    waitGroup.Add(1)
    go func(){
      defer waitGroup.Done()
      connectionCache.Get("user:password@/dbname", createConnectionManager)
    }()
  }

  waitGroup.Wait()

  assert.Equal(t, connectionCache.Size(), 1)
}

func createConnectionManager(connectionString string) (interface{}, error) {
  db, _, err := sqlmock.New()

  if err != nil {
    return nil, err 
  }

  return db, nil
}

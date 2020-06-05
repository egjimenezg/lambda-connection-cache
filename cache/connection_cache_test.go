package cache_test

import (
  "database/sql"
  "github.com/egjimenezg/lambda-connection-cache/cache"
  "github.com/stretchr/testify/assert"
  "github.com/DATA-DOG/go-sqlmock"
  "reflect"
  "testing"
)

func TestCreateNewConnectionCacheOnlyOnce(t *testing.T){
  connectionCache := cache.New()
  assert.Equal(t, connectionCache, cache.New())
}

func TestCreateNewConnectionPoolAndSavingItInCacheMap(t *testing.T) {

  connectionFunction := func(connectionString string) (interface{}, error) {
    db, _, err := sqlmock.New()

    if err != nil {
      t.Fatalf("An error '%s' has ocurred opening the stub database connection", err)
    }

    return db, err
  }

  connectionCache := cache.New()

  connManagerA, err := connectionCache.Get("user:password@/dbname", connectionFunction)

  if err != nil {
    t.Fatalf("An error '%s' has ocurred", err)
  }

  connManagerB, err := connectionCache.Get("user:password@/dbname", connectionFunction)

  assert.True(t, reflect.DeepEqual(connManagerA.(*sql.DB), connManagerB.(*sql.DB)))
}

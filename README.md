# Lambda connection cache

Simple thread safe map used to store database handlers created by lambda functions.

[![Build Status](https://travis-ci.org/egjimenezg/lambda-connection-cache.svg?branch=master)](https://travis-ci.org/egjimenezg/lambda-connection-cache)

# Installation

```
go get github.com/egjimenezg/lambda-connection-cache
```

# How to use

Create a function that gets the handler for your database into your repository and use the **Get** function with a connection string to avoid the creation of multiple handlers on consecutive calls.

```
package main

import (
  "context"
  "github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
  Body string `json:body`
}

type Response struct {
  StatusCode int `json:statusCode`
}

func HandleRequest(ctx context.Context, request Request) (Response, error){
  repository, err := NewRepository()

  if err != nil {
    return Response{}, err
  }

  service := NewService(repository)
  response := service.DoSomething()

  return response, nil
}

func main(){
  lambda.Start(HandleRequest)
}
```

repository.go

```
package main

import (
  "database/sql"
  cache "github.com/egjimenezg/lambda-connection-cache"
  _ "github.com/lib/pq"
)

type Repository struct {
  databaseHandler *sql.DB
}

func NewRepository() (*Repository, error) {
  connCache := cache.New()
  databaseHandler, err := connCache.Get("host=host port=5432 user=user password=n0m3l0s3 dbname=dbname", getDatabaseManager)

  if err != nil {
    return nil, err
  }

  return &Repository{databaseHandler: databaseHandler.(*sql.DB)}, nil
}

func getDatabaseManager(connectionString string) (interface{}, error){
  db, err := sql.Open("postgres", connectionString)

  if err != nil {
    return nil, err
  }

  db.SetMaxOpenConns(100)

  return db, nil
}
```

Using the **connCache.Get** method will avoid getting the **Too many connections** error.

Don't forget to configure the db manager options for a better database performance.

* **MaxOpenConns**
* **MaxIdleConns**
* **ConnMaxLifetime**


[![CircleCI](https://circleci.com/gh/lab259/http-mysql-service.svg?style=shield)](https://circleci.com/gh/lab259/http-mysql-service)
[![codecov](https://codecov.io/gh/lab259/http-mysql-service/branch/master/graph/badge.svg)](https://codecov.io/gh/lab259/http-mysql-service)

# http-mysql-service

The http-mysql-service is the [lab259/http](//github.com/lab259/http) service
implementation for the [go-sql-driver/mysql](//github.com/go-sql-driver/mysql) library.

## Dependencies

It depends on the [lab259/http](//github.com/lab259/http) (and its dependencies,
of course) itself and the [go-sql-driver/mysql](//github.com/go-sql-driver/mysql) library.

## Installation

First, fetch the library to the repository.

```bash
go get github.com/lab259/http-mysql-service
```

## Usage

Applying configuration and starting service

```go
// Create MySQLService instance
var mySQLService MySQLService

// Applying configuration
err := mySQLService.ApplyConfiguration(MySQLServiceConfiguration{
    Host:        "host",
    User:        "user",
    Password:    "password",
    Database:    "database",
    Port:        3306,
    MaxPoolSize: 1,
})

if err != nil {
    panic(err)
}
        
// Starting service
err := mySQLService.Start()

if err != nil {
    panic(err)
}

// Create a custom string
var value string

// Executing something using a *sql.Conn
err := mySQLService.RunWithConn(func(conn *sql.Conn) error {
    // Retrieving a value from the MySQL
    ctx := context.Background()
	rows, err := conn.QueryContext(ctx, "SELECT value FROM some-table WHERE id=?", "my-custom-id")
        
    if err != nil {
        return err
    }
    
    defer rows.Close()

    if !rows.Next() {
        return errors.New("some-table is empty")
    }
		
    err = rows.Scan(&value)

    if err != nil {
        return err
    }
        
    return nil
})

if err != nil {
    panic(err)
}

value // "some-value-get-from-my-sql"
```

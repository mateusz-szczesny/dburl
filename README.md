
[![Build Status](https://travis-ci.com/mateusz-szczesny/dburl.svg?branch=master)](https://travis-ci.com/mateusz-szczesny/dburl)

# DB Url

Package dburl allows you to configure your database configuration based on DATABASE_URL environment variable.

The package provides Config(string) function to generate configuration based on environment variable passed as an argument.
Function returns a database config structure, populated with all the data specified in your environment variable.
Package helps to create config for docker based applications, which depends on external database.

If you'd rather not use an environment variable, you can pass a URL in directly instead to Parse(string) method.

## Supported Databases

Currently supports the following databases: **PostgreSQL**, **MSSQL** and **SQLite**. (I plan to expand this list in the future)

### Installation

Installation is simple
```bash
    $ go get github.com/mateusz-szczesny/dburl
```
### Usage

Configure your database with default environment variable ("DATABASE_URL")
```go
    import (
        ...
        "github.com/mateusz-szczesny/dburl"
    )

    db, err := dburl.Config(dburl.DefaultEnv)
```
Provide a custom environment variable name
```go
    db, err := dburl.Config("MY_CUSTOM_VARIABLE_NAME")
```

## URL schema

| Engine        | URL                                        |
| :------------ | :------------                              |
| PostgreSQL    | `postgres://USER:PASSWORD@HOST:PORT/NAME`  |
| MSSQL         | `mssql://USER:PASSWORD@HOST:PORT/NAME`     |
| SQLite        | `sqlite:///path/to/some/database.db`       |

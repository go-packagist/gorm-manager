# Gorm'er —— GORM Manager

[![Go Report Card](https://goreportcard.com/badge/github.com/go-packagist/gormer)](https://goreportcard.com/report/github.com/go-packagist/gormer)
[![tests](https://github.com/go-packagist/gormer/actions/workflows/go.yml/badge.svg)](https://github.com/go-packagist/gormer/actions/workflows/go.yml)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-packagist/gormer)](https://pkg.go.dev/github.com/go-packagist/gormer)
[![codecov](https://codecov.io/gh/go-packagist/gormer/branch/master/graph/badge.svg?token=5TWGQ9DIRU)](https://codecov.io/gh/go-packagist/gormer)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## Installation

```bash
go get github.com/go-packagist/gormer
```

## Usage

```go
package main

import (
	"github.com/go-packagist/gormer"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Username string
	Password string
}

func main() {
	g := gormer.New(&gormer.Config{
		Default: "db1",
		Connections: map[string]gormer.ConnectionFunc{
			"db1": func() *gorm.DB {
				return nil
			},
			"db2": func() *gorm.DB {
				dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

				db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
				if err != nil {
					panic(err)
				}

				return db
			},
		},
	})

	g.Connection().First(&User{}) // use default connection
	g.Connection("db1").First(&User{}) // use db1 connection
	g.Connection("db2").First(&User{})
}
```

## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
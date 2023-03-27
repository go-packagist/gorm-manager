# Gorm'er —— GORM Manager

[![Go Version](https://badgen.net/github/release/go-packagist/gormer/stable)](https://github.com/go-packagist/gormer/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-packagist/gormer/v2)](https://pkg.go.dev/github.com/go-packagist/gormer/v2)
[![codecov](https://codecov.io/gh/go-packagist/gormer/branch/master/graph/badge.svg?token=5TWGQ9DIRU)](https://codecov.io/gh/go-packagist/gormer)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-packagist/gormer)](https://goreportcard.com/report/github.com/go-packagist/gormer)
[![tests](https://github.com/go-packagist/gormer/actions/workflows/go.yml/badge.svg)](https://github.com/go-packagist/gormer/actions/workflows/go.yml)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## Installation

```bash
go get github.com/go-packagist/gormer/v2
```

## Usage

```go
package main

import (
	"github.com/go-packagist/gormer/v2"
	"gorm.io/gorm"
)

func main() {
	m := gormer.NewManager(&gormer.Config{
		Default: "db1",
		Connections: map[string]gormer.Connection{
			"db1": &gormer.MySQLConfig{
				DSN: gormer.DSN{
					Host: "127.0.0.1",
					Port: 3306,
					User: "root",
					Pass: "123456",
					Db:   "test",
					Options: map[string]string{
						"charset":   "utf8mb4",
						"parseTime": "True",
					},
				}.String(),
				GormConfig: &gorm.Config{
					SkipDefaultTransaction: true,
				},
			},
			"db2": &gormer.MySQLConfig{
				DSN: gormer.DSN{
					Host: "...",
				}.String(),
			},
		},
	}, gormer.WithInstance) // Instance mode

	var user = struct {
		Name string
	}{}
	m.Connect("db1").DB.First(&user)
	gormer.Factory().Connect("db1").DB.First(&user) // use factory mode
	gormer.Gormer().Connect("db2").DB.First(&user)  // use instance mode
	gormer.Connect("db1").DB.First(&user)           // use instance mode
	gormer.Connect().DB.First(&user)                // use default connection
	gormer.Connect().First(&user)                   // use default connection and ignore `DB`

	println(gormer.Connect("db2").Err) // if db2 is err
}
```

## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
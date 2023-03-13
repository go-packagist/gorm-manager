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
	g := gormer.NewManager(&gormer.Config{
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

	g.First(&User{})              // use default connection(The test product could be discarded at any time)
	g.Connection().First(&User{}) // use default connection(Recommended usage)
	g.Connection("db1").First(&User{})
	g.Connection("db2").First(&User{})
}

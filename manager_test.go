package gormer

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"testing"
	"time"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func setup(db *Manager) {
	for _, dbname := range []string{"write", "read"} {
		db.Connect("gormer").DB.Exec("DROP DATABASE IF EXISTS ", "db_"+dbname)
		db.Connect(dbname).DB.Exec("DROP TABLE IF EXISTS `users`")
		db.Connect(dbname).DB.Exec("CREATE TABLE `users` (`id` int(11) NOT NULL AUTO_INCREMENT,`name` varchar(255) DEFAULT NULL,`age` int(11) DEFAULT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
		db.Connect(dbname).DB.Exec("INSERT INTO `users` (`id`, `name`, `age`) VALUES (1, ?, 18)", dbname+":test1")
		db.Connect(dbname).DB.Exec("INSERT INTO `users` (`id`, `name`, `age`) VALUES (2, ?, 20)", dbname+":test2")
	}
}

func TestGormer(t *testing.T) {
	db := newManager()
	setup(db)

	assert.Nil(t, db.Connect("write").Err)
	assert.Nil(t, db.Connect("read").Err)
	assert.Nil(t, db.Connect("rw").Err)

	t.Run("test write", func(t *testing.T) {
		user := getUser(db, "write")
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "write:test1", user.Name)
		assert.Equal(t, 18, user.Age)
	})

	t.Run("test read", func(t *testing.T) {
		user := getUser(db, "read")

		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "read:test1", user.Name)
		assert.Equal(t, 18, user.Age)
	})

	t.Run("test write and read", func(t *testing.T) {
		user := getUser(db, "rw")

		// valid read
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "read:test1", user.Name)
		assert.Equal(t, 18, user.Age)

		// valid write
		user.Age = 25
		db.Connect("rw").DB.Save(user)
		userWrite, userRead, userRw := getUser(db, "write"), getUser(db, "read"), getUser(db, "rw")

		// valid get from write
		assert.Equal(t, 1, userWrite.ID)
		assert.Equal(t, "read:test1", userWrite.Name)
		assert.Equal(t, 25, userWrite.Age)

		// valid get from read
		assert.Equal(t, 1, userRead.ID)
		assert.Equal(t, "read:test1", userRead.Name)
		assert.Equal(t, 18, userRead.Age)

		// valid get from rw
		assert.Equal(t, 1, userRw.ID)
		assert.Equal(t, "read:test1", userRw.Name)
		assert.Equal(t, 18, userRw.Age)
	})
}

func getUser(db *Manager, dbname string) *User {
	var user User
	db.Connect(dbname).DB.Model(&User{}).First(&user)

	return &user
}

func newDSN(dbname string) DSN {
	return DSN{
		Host: "127.0.0.1",
		Port: 3306,
		User: "root",
		Pass: "123456",
		Db:   dbname,
	}
}

func newManager() *Manager {
	var (
		gormerDSN = newDSN("gormer")
		writeDSN  = newDSN("db_write")
		readDSN   = newDSN("db_read")
	)

	return NewManager(&Config{
		Default: "gormer",
		Connections: map[string]Connection{
			"gormer": &MySQLConfig{
				DSN: gormerDSN.String(),
				GormConfig: &gorm.Config{
					SkipDefaultTransaction: true, // 禁用默认事务
				},
			},
			"write": &MySQLConfig{
				DSN: writeDSN.String(),
				GormConfig: &gorm.Config{
					SkipDefaultTransaction: true, // 禁用默认事务
				},
			},
			"read": &MySQLConfig{
				DSN: readDSN.String(),
				GormConfig: &gorm.Config{
					SkipDefaultTransaction: true, // 禁用默认事务
				},
			},
			"rw": &MySQLConfig{
				DSN: writeDSN.String(),
				GormConfig: &gorm.Config{
					SkipDefaultTransaction: true, // 禁用默认事务
				},
				GormUsees: []gorm.Plugin{
					dbresolver.Register(
						dbresolver.Config{
							Sources:  []gorm.Dialector{mysql.Open(writeDSN.String())},
							Replicas: []gorm.Dialector{mysql.Open(readDSN.String())},
							Policy:   dbresolver.RandomPolicy{},
						}).
						SetConnMaxLifetime(time.Hour).
						SetConnMaxIdleTime(24 * time.Hour).
						SetMaxIdleConns(100).
						SetMaxOpenConns(200),
				},
			},
		},
	})
}

package repositories

import (
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

var (
	mysqlDB   *gorm.DB
	rdsClient redis.UniversalClient
)

func init() {
	mysqlDB = connectSqlite()
	rdsClient = connectRedis()
}

func connectSqlite() *gorm.DB {
	db, err := gorm.Open("sqlite3", "MzDB.db")
	if err != nil {
		panic(err)
	}

	// DEBUG
	db.LogMode(true)
	db.SetLogger(logrus.StandardLogger())

	db.AutoMigrate(&User{})
	return db
}

func connectRedis() redis.UniversalClient {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"127.0.0.1:6379"},
		Password: "xc315215241",
	})
	if err := client.Ping().Err(); err != nil {
		panic(err)
	}
	return client
}

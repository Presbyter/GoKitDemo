package repositories

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

var (
	mysqlDB *gorm.DB
)

func init() {
	db, err := gorm.Open("sqlite3", "MzDB.db")
	if err != nil {
		panic(err)
	}

	// DEBUG
	db.LogMode(true)
	db.SetLogger(logrus.StandardLogger())

	db.AutoMigrate(&User{})

	mysqlDB = db
}

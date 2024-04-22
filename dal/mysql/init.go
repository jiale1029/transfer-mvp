package mysql

import (
	"fmt"
	"os"

	"github.com/jiale1029/transaction/entity"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	user     = "root"
	password = ""
	addr     = "localhost"
	port     = "3306"
	dbName   = "take_home_test"
)

func InitMySQL() *gorm.DB {
	var (
		dsn string
		db  *gorm.DB
		err error
	)

	isUT := os.Getenv("IS_UNIT_TEST")
	if isUT == "1" {
		dsn = "file::memory:?cache=shared"
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	} else {
		if password != "" {
			dsn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", user, password, addr, port, dbName)
		} else {
			dsn = fmt.Sprintf("%v@tcp(%v:%v)/%v?parseTime=true", user, addr, port, dbName)
		}
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect to database")
		}
	}

	err = db.AutoMigrate(
		&entity.Account{},
		&entity.Transaction{},
	)
	if err != nil {
		panic(err)
	}
	return db
}

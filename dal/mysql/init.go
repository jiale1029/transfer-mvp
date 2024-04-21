package mysql

import (
	"github.com/jiale1029/transaction/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func init() {
	dsn := "root@tcp(localhost:3306)/take_home_test?parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	err = db.AutoMigrate(
		&entity.Account{},
		&entity.Transaction{},
	)
	if err != nil {
		panic(err)
	}
	Database = db
}

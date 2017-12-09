package model

import (
	"fmt"
	"../config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	DBConn *gorm.DB
)

func GormInit() error {
	conf := &config.DBConfig{}
	err := conf.Read()
	if err != nil {
		fmt.Println(err)
		return err
	}
	DBConn, err = gorm.Open("postgres",
		fmt.Sprintf("host=localhost user=%s dbname=%s sslmode=disable password=%s",
			conf.DBUser, conf.DBName, conf.DBPass))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GormClose() error {
	if DBConn != nil {
		return DBConn.Close()
	}
	return nil
}

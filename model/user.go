package model

import (
"errors"
)

type Usr struct {
	//gorm.Model
	ID         int64  `gorm:primary key;not_nil`
	Login      string `gorm:not_nil`
	Pass       string `gorm:not_nil`
	WorkNumber int32
}

var err_db_nul error = errors.New("DBConn is nil")

func (u *Usr) Get(login string, pass string) error {
	if DBConn != nil {
		return DBConn.Where("login =?",login).Where("pass =?",pass).First(u).Error
	} else {
		return err_db_nul
	}
}

func (u *Usr) Save(ID int64, newpass string) error {
	if DBConn != nil {
		err := DBConn.Where("id =?",ID).First(u).Error
		if err != nil {
			return err
		} 
		return DBConn.Model(&u).Update("pass", newpass).Error
	} else {
		return err_db_nul
	}
}
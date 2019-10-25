package dbHelper

import (
	"fmt"
	"github.com/funtoy/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Orm *Database

type Database struct {
	*gorm.DB
}

func NewOrm(ip, user, pass, name string) *Database {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, ip, name)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库不可用， 错误信息：", err.Error())
		return nil
	}

	// 初始化tables
	ret := &Database{db}

	Orm = ret

	return ret
}

func (d *Database) CheckTable(v interface{}) {
	if !d.HasTable(v) {
		if err := d.CreateTable(v).Error; err != nil {
			log.Fatal(err.Error())
		}

	} else {
		if err := d.AutoMigrate(v).Error; err != nil {
			log.Fatal(err.Error())
		}
	}
}

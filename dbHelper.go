package dbHelper

import (
	"fmt"
	"github.com/funtoy/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type DBType string

const (
	DBTypeMySql DBType = "mysql"
	DBTypePgSql DBType = "postgres"
)

var (
	writer *Database
	reader *Database
)

func Writer() *Database {
	return writer
}

func Reader() *Database {
	return reader
}

func InitWriter(t DBType, ip, port, user, pass, name string) {
	writer = NewOrm(t, ip, port, user, pass, name)
}

func InitReader(t DBType, ip, port, user, pass, name string) {
	reader = NewOrm(t, ip, port, user, pass, name)
}

type Database struct {
	*gorm.DB
}

func NewOrm(t DBType, ip, port, user, pass, name string) *Database {
	var dsn string
	if t == DBTypeMySql {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, ip, port, name)

	} else if t == DBTypePgSql {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", ip, port, user, name, pass)
	}
	db, err := gorm.Open(string(t), dsn)
	if err != nil {
		log.Fatal("数据库不可用， 错误信息：", err.Error())
		return nil
	}

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(10)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置连接的最大可复用时间。
	db.DB().SetConnMaxLifetime(time.Hour)

	// 初始化tables
	ret := &Database{db}

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

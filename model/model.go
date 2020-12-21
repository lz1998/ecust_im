package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func init() {
	dsn := "im:xxx@tcp(tmp.lz1998.xin:13306)/im?charset=utf8mb4&parseTime=True&loc=Local" // TODO gitignore 防止密码泄露
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	Db = db
}

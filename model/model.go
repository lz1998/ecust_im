package model

import (
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB       // MySQL
var RDb *redis.Client // Redis
// TODO leveldb

func init() {

	// MySQL
	pass := os.Getenv("MYSQL_PASSWORD")
	dsn := "im:" + pass + "@tcp(tmp.lz1998.xin:13306)/im?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	Db = db

	RDb = redis.NewClient(&redis.Options{
		Addr:     "tmp.lz1998:16379",
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})
}

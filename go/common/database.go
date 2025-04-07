// @Title  database
// @Description  该文件用于初始化mysql数据库，以及包装一个向外提供数据库的功能
package common

import (
	"fmt"
	"lianjiang/model"
	"net/url"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// @title    InitDB
// @description   从配置文件中读取数据库相关信息后，完成数据库初始化
// @param     void        void         没有入参
// @return    db        *gorm.DB         将返回一个初始化后的数据库指针
func InitDB() *gorm.DB {
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc),
	)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	// TODO  如果未能连接到数据库，终止程序并返回错误信息
	if err != nil {
		log.Fatalf("连接mysql数据库失败, 原因: " + err.Error())
	}
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.FileHistory{})
	db.AutoMigrate(&model.DataHistory{})
	db.AutoMigrate(&model.MapHistory{})
	db.AutoMigrate(&model.FileInfo{})
	db.AutoMigrate(&model.MapVersion{})
	db.AutoMigrate(&model.MapVersionDetail{})
	db.AutoMigrate(&model.DataTableInfo{})
	DB = db
	return db
}

// @title    GetDB
// @description   返回数据库的指针
// @param     void        void         没有入参
// @return    db        *gorm.DB         将返回一个初始化后的数据库指针
func GetDB() *gorm.DB {
	return DB
}


// 
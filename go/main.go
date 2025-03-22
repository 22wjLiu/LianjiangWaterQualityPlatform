// @Title  main
// @Description  程序的入口，读取配置，调用初始化函数以及运行路由
package main

import (
	"fmt"
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/util"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// @title    main
// @description   程序入口，完成一些初始化工作后将开始监听
// @param     void			没有入参
// @return    void			没有回参
func main() {
	InitConfig()
	common.InitDB()
	client0 := common.InitRedis(0)
	defer client0.Close()
	r := gin.Default()
	r = CollectRoute(r)
	port := viper.GetString("server.port")
	// TODO 定时备份映射表
	go func() {
		for {
			log.Println("备份映射表正在进行中...")
			// TODO 执行备份功能
			BackUp()
			log.Println("备份映射表完成...")

			now := time.Now()

			// TODO 计算下一个4:00
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 4, 00, 0, 0, next.Location())

                        sleepDuration := next.Sub(now)

			hours := sleepDuration / time.Hour
                        minutes := (sleepDuration % time.Hour) / time.Minute
			seconds := (sleepDuration % time.Minute) / time.Second

			log.Printf("下次备份时间: %s (等待 %02d时%02d分%02d秒)\n", next.Format("2006-01-02 15:04:05"), hours, minutes, seconds)
			time.Sleep(next.Sub(now))
		}
	}()
	if port != "" {
		log.Panic(r.Run(":" + port))
	}
	log.Panic(r.Run())
}

// @title    InitConfig
// @description   读取配置文件并完成初始化
// @param     void			没有入参
// @return    void			没有回参
func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Panicf("获取当前目录失败: %v", err)
	}

	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err = viper.ReadInConfig()
	// TODO 如果发生错误，终止程序
	if err != nil {
		log.Panicf("读取配置失败: %v", err)
	}
}

// @title    BackUp
// @description   备份映射文件
// @param     void			没有入参
// @return    void			没有回参
func BackUp() {
	db := common.GetDB()
	for id, v := range util.MapMap {
		for _, key := range v.Keys() {
			value, _ := v.Get(key)
			db.Create(&model.MapBackup{
				Table: id,
				Key:   key,
				Value: fmt.Sprint(value),
			})
		}
	}
}

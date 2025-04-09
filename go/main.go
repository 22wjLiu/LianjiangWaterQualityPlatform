// @Title  main
// @Description  程序的入口，读取配置，调用初始化函数以及运行路由
package main

import (
	"lianjiang/common"
	"lianjiang/util"
	"os"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// @title    main
// @description   程序入口，完成一些初始化工作后将开始监听
// @param     void			没有入参
// @return    void			没有回参
func main() {
	// 初始化后端配置
	InitConfig()
	
	// 初始化数据库
	common.InitDB()
	client0 := common.InitRedis(0)
	defer client0.Close()

	// 初始化路由
	r := gin.Default()
	r = CollectRoute(r)

	// 初始化系统用户
	if err := util.InitSysUser(); err != nil {
		log.Fatal("初始化系统用户失败：", err)
	}

	// 初始化映射表
	log.Println("开始初始化映射表...")
	if err := util.InitMapMap(); err != nil {
		log.Fatal("初始化映射表失败：", err)
	}
	log.Println("初始化映射表成功")

	// 定时备份映射表
	go func() {
		for {
			log.Println("开始备份映射表...")

			// 执行备份功能
			if err := util.BackUpSql(util.MapSqlPath, util.MapSqlFiles); err != nil {
				log.Println(err)
			} else {
				log.Println("备份映射表完成")
			}
			
			// 计算下次备份时间（每天凌晨4:00备份）
			now := time.Now()
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

	// 开启服务
	port := viper.GetString("server.port")
	if port != "" {
		log.Fatal(r.Run(":" + port))
	} else {
		log.Fatal(r.Run())
	}
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
	// viper.SetConfigName("2")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")

	// 读取配置文件，如果发生错误，终止程序
	err = viper.ReadInConfig()
	if err != nil {
		log.Panicf("读取配置失败: %v", err)
	}
}

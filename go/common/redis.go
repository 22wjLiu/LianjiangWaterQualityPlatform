// @Title  redis
// @Description  该文件用于初始化redis数据库，以及包装一个向外提供数据库的功能
package common

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var Client [16]*redis.Client

// @title    InitRedis
// @description   从配置文件中读取数据库相关信息后，完成数据库初始化
// @param     val int         表示初始化第几号redis库
// @return    *redis.Client         将返回一个初始化后的redis数据库指针
func InitRedis(val int) *redis.Client {
	host := viper.GetString("redis.host")
	addr := viper.GetString("redis.port")
	password := viper.GetString("redis.password")

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + addr,
		Password: password,
		DB:       val,
	})
	Client[val] = client
	return client
}

// @title    GetRedisClient
// @description   返回数据库的指针
// @param     void        void         没有入参
// @return    *redis.Client         将返回一个初始化后的数据库指针
func GetRedisClient(val int) *redis.Client {
	return Client[val]
}

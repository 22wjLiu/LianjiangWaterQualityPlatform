// @Title  RecoveryMiddleware
// @Description  该中间件用于“拦截”运行时恐慌的内建函数,防止程序崩溃
package middleware

import (
	"fmt"
	"lianjiang/response"

	"github.com/gin-gonic/gin"
)

// @title    RecoveryMiddleware
// @description   该中间件用于“拦截”运行时恐慌的内建函数,防止程序崩溃
// @param     void        void    		  无入参
// @return    HandlerFunc        gin.HandlerFunc            返回一个响应函数
func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			// TODO recover用于“拦截”运行时恐慌的内建函数,防止程序崩溃
			if err := recover(); err != nil {
				response.Fail(ctx, nil, fmt.Sprint(err))
			}
		}()
		ctx.Next()
	}
}

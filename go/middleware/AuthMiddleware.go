// @Title  AuthMiddleware
// @Description  中间件，用于解析token
package middleware

import (
	"fmt"
	"lianjiang/common"
	"lianjiang/model"
	"strings"

	"github.com/gin-gonic/gin"
)

// @title    AuthMiddleware
// @description   中间件，用于解析token
// @param    void
// @return   gin.HandlerFunc	将token解析完毕后传回上下文
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 authorization header
		tokenString := ctx.GetHeader("Authorization")

		fmt.Println("请求token", tokenString)

		if tokenString == "" {
			ctx.JSON(201, gin.H{
				"code": 201,
				"msg":  "请先登录",
			})
			ctx.Abort()
			return
		}

		// validate token formate
		if  !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(202, gin.H{
				"code": 202,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		// TODO 截取字符
		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)

		if err != nil || !token.Valid {
			ctx.JSON(202, gin.H{
				"code": 202,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		// TODO token通过验证, 获取claims中的UserID
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// TODO 验证用户是否存在
		if user.Id == 0 {
			ctx.JSON(203, gin.H{
				"code": 203,
				"msg":  "用户不存在，请先注册后登录",
			})
			ctx.Abort()
			return
		}

		// TODO 用户存在 将user信息写入上下文
		ctx.Set("user", user)

		ctx.Next()
	}
}

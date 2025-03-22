// @Title  PersonalController
// @Description  该文件用于提供操作个人界面的各种函数
package controller

import (
	"lianjiang/common"
	"lianjiang/dto"
	"lianjiang/model"
	"lianjiang/response"
	"lianjiang/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @title    PersonalPage
// @description   提供用户的个人信息
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalPage(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	response.Success(ctx, gin.H{"user": dto.ToUserDto(user)}, "请求成功")

}

// @title    Level
// @description   用户设置其它用户的权限等级
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Level(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在五级以下的用户不能设置其它用户的等级
	if user.Level < 5 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 获取path中的id
	userId := ctx.Params.ByName("id")

	// TODO 尝试在数据库中查找这个用户
	db := common.GetDB()

	if db.Where("id = ?", userId).First(&model.User{}).Error != nil {
		response.Fail(ctx, nil, "不存在该用户")
		return
	}

	// TODO 获取path中的level
	level, _ := strconv.Atoi(ctx.Params.ByName("level"))

	db.Where("id = ?", userId).Update("level", level)

	response.Success(ctx, nil, "更新成功")

}

// @title    Users
// @description   用户查看其它用户
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Users(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在四级以下的用户不能查看其它用户的等级
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 尝试在数据库中查找这所有用户
	db := common.GetDB()

	var users []dto.UserDto

	if db.Table("users").Find(users).Error != nil {
		response.Fail(ctx, nil, "用户表出错")
		return
	}

	response.Success(ctx, gin.H{"users": users}, "查找成功")

}

// @title    FindUser
// @description   用户通过search字段查看其它用户
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func FindUser(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在四级以下的用户不能查看其它用户
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 尝试在数据库中查找这所有用户
	db := common.GetDB()

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取path中的search
	search := ctx.Params.ByName("search")

	field, ok := util.UserMap[search]

	if !ok {
		response.Fail(ctx, nil, "字段"+search+"不存在")
		return
	}

	// TODO 通过field字段搜寻用户
	if db.Where(field+" = ?", id).First(user).Error != nil {
		response.Fail(ctx, nil, "该用户不存在")
		return
	}

	response.Success(ctx, gin.H{"user": user}, "查找成功")
}

// @Title  PersonalController
// @Description  该文件用于提供操作个人界面的各种函数
package controller

import (
	"lianjiang/common"
	"lianjiang/dto"
	"lianjiang/model"
	"lianjiang/response"
	"strconv"
	"time"
	"fmt"

	"gorm.io/gorm"

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


// @title    Users
// @description   用户查看其它用户
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Users(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在四级以下的用户不能查询用户信息
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 尝试在数据库中查找这所有用户
	db := common.GetDB().Table("users")

	// 时间格式定义
	layout := "2006-01-02T15:04:05"

	// 读取参数请求
	start := ctx.Params.ByName("start")

	if start != "" && start != "null" {
		start, err := time.Parse(layout, start)
		if err == nil{
			db = db.Where("created_at >= ?", start)
		} else {
			response.Fail(ctx, nil, "错误的用户创建开始时间")
			return
		}
	}

	end := ctx.Params.ByName("end")

	if end != "" && end != "null" {
		end, err := time.Parse(layout, end)
		if err == nil{
			db = db.Where("created_at <= ?", end)
		} else {
			response.Fail(ctx, nil, "错误的用户创建结束时间")
			return
		}
	}

	applyFilter := func(field, value string) {
		if value != "" {
			db = db.Where("`" + field + "` = ?", value)
		}
	}

	applyFilter_2 := func(field, value string) {
		if value != "" {
			db = db.Where("`" + field + "` like ?", "%" + value + "%")
		}
	}

	// 获取ID
	applyFilter("id", ctx.DefaultQuery("id", ""))
	// 获取操作方式
	applyFilter_2("user_name", ctx.DefaultQuery("userName", ""))
	// 获取用户ID
	applyFilter_2("email", ctx.DefaultQuery("email", ""))
	// 获取键
	applyFilter("level", ctx.DefaultQuery("level", ""))
	// 获取分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "25"))
	
	offset := (page - 1) * pageSize

	var users []dto.UserDto

	result := db.
		Limit(pageSize).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		response.Fail(ctx, nil, "参数有误")
		return
	}

	// 查询总数
	var total int64
	dbCount := db.Session(&gorm.Session{})
	dbCount.Count(&total)
	
	// 返回分页数据
	response.Success(ctx, gin.H{
		"users": users,
		"total": total,
	}, "查询成功")

}

// @title    DeleteUser
// @description   删除用户
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteUser(ctx *gin.Context){
	// 获取登录用户
	tuser, exists := ctx.Get("user")
	if !exists {
		response.Fail(ctx, nil, "获取用户信息失败")
		return
	}
	
	user, ok := tuser.(model.User)
	if !ok {
		response.Fail(ctx, nil, "用户信息解析失败")
		return
	}
	
	// 安全等级在四级以下的用户不能删除用户
	if user.Level < 5 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// 获取path中的id
	userIdInt, err := strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		response.Fail(ctx, nil, "用户ID格式错误")
		return
	}

	userId := uint(userIdInt)

	if userId == user.Id {
		response.Fail(ctx, nil, "不能删除自己")
		return
	}

	// 尝试在数据库中查找这个用户
	db := common.GetDB()

	var userToDelete model.User

	if db.Where("id = ?", userId).First(&userToDelete).Error != nil {
		response.Fail(ctx, nil, "不存在该用户")
		return
	}

	result := db.Delete(&userToDelete)

	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败: "+result.Error.Error())
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// @title    DeleteUsers
// @description   删除一组用户
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteUsers(ctx *gin.Context) {
	// 获取登录用户
	tuser, exists := ctx.Get("user")
	if !exists {
		response.Fail(ctx, nil, "获取用户信息失败")
		return
	}
	
	user, ok := tuser.(model.User)
	if !ok {
		response.Fail(ctx, nil, "用户信息解析失败")
		return
	}

	// 安全等级在四级以下的用户不能删除用户
	if user.Level < 5 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// 解析请求体
	var resq dto.BatchDeleteId
	if err := ctx.ShouldBindJSON(&resq); err != nil {
		response.Fail(ctx, nil, "请求体参数解析错误")
		return
	}

	// 过滤掉自己的 ID
	for _, id := range resq.Ids {
		if id == user.Id {
			response.Fail(ctx, nil, "不能删除自己")
			return
		}
	}
	

	// 开启事务
	db := common.GetDB()
	tx := db.Begin()

	for _, id := range resq.Ids {
		if err := tx.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, fmt.Sprintf("删除用户ID为 %d 的用户失败，已回滚", id))
			return
		}
	}

	tx.Commit()
	response.Success(ctx, nil, "删除成功")
}

// @title    UpdateUser
// @description   用户设置其它用户的权限等级
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func UpdateUser(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// 安全等级在五级以下的用户不能修改其它用户的信息
	if user.Level < 5 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// 获取path中的id
	userId := ctx.Params.ByName("id")

	// 尝试在数据库中查找这个用户
	db := common.GetDB()

	if db.Where("id = ?", userId).First(&model.User{}).Error != nil {
		response.Fail(ctx, nil, "不存在该用户")
		return
	}

	var userData dto.UserDto

	
	// 解析请求体参数
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		response.Fail(ctx, nil, "请求体参数解析错误")
		return
	}

	// 构建更新字段的 map
	updates := map[string]interface{}{}

	if userData.Level != 0 {
		updates["level"] = userData.Level
	}

	if userData.Name != "" {
		updates["name"] = userData.Name
	}

	// 执行更新
	result := db.
		Model(&model.User{}).
		Where("id = ?", userId).
		Updates(updates)
	
	if result.Error != nil {
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, nil, "更新成功")

}
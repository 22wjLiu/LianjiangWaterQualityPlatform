// @Title  HistoryController
// @Description  该文件用于提供获取历史操作记录的各种函数
package controller

import (
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/response"

	"github.com/gin-gonic/gin"
)

// @title    FileHistory
// @description   提供文件的操作记录
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func FileHistory(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在四级以下的用户不能查看历史操作记录
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	db := common.GetDB().Table("file_historys")

	var fileHistorys []model.FileHistory

	// TODO 读取参数请求
	start := ctx.Params.ByName("start")

	if start != "" {
		db = db.Where("created_at >= ", start)
	}

	end := ctx.Params.ByName("end")

	if end != "" {
		db = db.Where("created_at <= ", end)
	}

	// TODO 取出用户id
	userId := ctx.DefaultQuery("id", "")

	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}

	// TODO 操作方式
	option := ctx.DefaultQuery("option", "")

	if option != "" {
		db = db.Where("option = ?", option)
	}

	if db.Find(fileHistorys).Error != nil {
		response.Fail(ctx, nil, "参数有误")
		return
	}

	response.Success(ctx, gin.H{"fileHistorys": fileHistorys}, "请求成功")

}

// @title    DataHistory
// @description   提供数据的操作记录
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DataHistory(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在四级以下的用户不能查看历史操作记录
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	db := common.GetDB().Table("data_historys")

	var dataHistorys []model.DataHistory

	// TODO 读取参数请求
	start := ctx.Params.ByName("start")

	if start != "" {
		db = db.Where("created_at >= ", start)
	}

	end := ctx.Params.ByName("end")

	if end != "" {
		db = db.Where("created_at <= ", end)
	}

	// TODO 取出用户id
	userId := ctx.DefaultQuery("id", "")

	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}

	// TODO 操作方式
	option := ctx.DefaultQuery("option", "")

	if option != "" {
		db = db.Where("option = ?", option)
	}

	// TODO 站名
	station_name := ctx.DefaultQuery("station_name", "")

	if station_name != "" {
		db = db.Where("station_name = ?", station_name)
	}

	// TODO 制度
	system := ctx.DefaultQuery("system", "")

	if system != "" {
		db = db.Where("system = ?", system)
	}

	if db.Find(dataHistorys).Error != nil {
		response.Fail(ctx, nil, "参数有误")
		return
	}

	response.Success(ctx, gin.H{"dataHistorys": dataHistorys}, "请求成功")

}

// @title    MapHistory
// @description   提供映射的操作记录
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func MapHistory(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在四级以下的用户不能查看历史操作记录
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	db := common.GetDB().Table("map_historys")

	var mapHistorys []model.MapHistory

	// TODO 读取参数请求
	start := ctx.Params.ByName("id")

	if start != "" {
		db = db.Where("created_at >= ", start)
	}

	end := ctx.Params.ByName("end")

	if end != "" {
		db = db.Where("created_at <= ", end)
	}

	// TODO 取出表id
	Id := ctx.DefaultQuery("id", "")

	if Id != "" {
		db = db.Where("id = ?", Id)
	}

	// TODO 操作方式
	option := ctx.DefaultQuery("option", "")

	if option != "" {
		db = db.Where("option = ?", option)
	}

	// TODO 取出用户id
	userId := ctx.DefaultQuery("userId", "")

	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}

	if db.Find(mapHistorys).Error != nil {
		response.Fail(ctx, nil, "参数有误")
		return
	}

	response.Success(ctx, gin.H{"mapHistorys": mapHistorys}, "请求成功")

}

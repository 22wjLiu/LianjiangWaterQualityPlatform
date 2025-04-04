// @Title  HistoryController
// @Description  该文件用于提供获取历史操作记录的各种函数
package controller

import (
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/response"
	"time"
	"strconv"

    "gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// @title    FileHistory
// @description   提供文件的操作记录
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func FileHistory(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// 安全等级在四级以下的用户不能查看历史操作记录
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	db := common.GetDB().Table("file_histories")

	var fileHistories []model.FileHistory

	// 时间格式定义
	layout := "2006-01-02T15:04:05"

	// 读取参数请求
	start := ctx.Params.ByName("start")

	if start != "" && start != "null" {
		start, err := time.Parse(layout, start)
		if err == nil{
			db = db.Where("created_at >= ?", start)
		} else {
			response.Fail(ctx, nil, "错误的文件日志开始时间")
			return
		}
	}

	end := ctx.Params.ByName("end")

	if end != "" && end != "null" {
		end, err := time.Parse(layout, end)
		if err == nil{
			db = db.Where("created_at <= ?", end)
		} else {
			response.Fail(ctx, nil, "错误的数据日志结束时间")
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

	// 用户id
	applyFilter("user_id", ctx.DefaultQuery("id", ""))
	// 文件名
	applyFilter_2("file_name", ctx.DefaultQuery("fileName", ""))
	// 操作方式
	applyFilter("option", ctx.DefaultQuery("option", ""))
	// 获取分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "25"))
		
	offset := (page - 1) * pageSize
	
	result := db.
		Limit(pageSize).
		Offset(offset).
		Order("created_at desc").
		Find(&fileHistories)
	
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
		"fileHistories": fileHistories,
		"total": total,
	}, "查询成功")
}

// @title    DataHistory
// @description   提供数据的操作记录
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DataHistory(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// 安全等级在四级以下的用户不能查看历史操作记录
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	db := common.GetDB().Table("data_histories")

	var dataHistories []model.DataHistory

	// 时间格式定义
	layout := "2006-01-02T15:04:05"

	// 读取参数请求
	start := ctx.Params.ByName("start")

	if start != "" && start != "null" {
		start, err := time.Parse(layout, start)
		if err == nil{
			db = db.Where("created_at >= ?", start)
		} else {
			response.Fail(ctx, nil, "错误的数据日志开始时间")
			return
		}
	}

	end := ctx.Params.ByName("end")

	if end != "" && end != "null" {
		end, err := time.Parse(layout, end)
		if err == nil{
			db = db.Where("created_at <= ?", end)
		} else {
			response.Fail(ctx, nil, "错误的数据日志结束时间")
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

	// 获取取出用户id
	applyFilter("user_id", ctx.DefaultQuery("id", ""))
	// 获取操作方式
	applyFilter("option", ctx.DefaultQuery("option", ""))
	// 获取站名
	applyFilter_2("station_name", ctx.DefaultQuery("station_name", ""))
	// 获取制度
	applyFilter("system", ctx.DefaultQuery("system", ""))
	// 获取分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "25"))
	
	offset := (page - 1) * pageSize

	result := db.
		Limit(pageSize).
		Offset(offset).
		Order("created_at desc").
		Find(&dataHistories)

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
		"dataHistories": dataHistories,
		"total": total,
	}, "查询成功")
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

	db := common.GetDB().Table("map_histories")

	var mapHistories []model.MapHistory

	// 时间格式定义
	layout := "2006-01-02T15:04:05"

	// 读取参数请求
	start := ctx.Params.ByName("start")

	if start != "" && start != "null" {
		start, err := time.Parse(layout, start)
		if err == nil{
			db = db.Where("created_at >= ?", start)
		} else {
			response.Fail(ctx, nil, "错误的映射日志开始时间")
			return
		}
	}

	end := ctx.Params.ByName("end")

	if end != "" && end != "null" {
		end, err := time.Parse(layout, end)
		if err == nil{
			db = db.Where("created_at <= ?", end)
		} else {
			response.Fail(ctx, nil, "错误的映射日志结束时间")
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
	applyFilter("option", ctx.DefaultQuery("option", ""))
	// 获取用户ID
	applyFilter("user_id", ctx.DefaultQuery("userId", ""))
	// 获取键
	applyFilter_2("key", ctx.DefaultQuery("key", ""))
	// 获取分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "25"))
	
	offset := (page - 1) * pageSize

	result := db.
		Limit(pageSize).
		Offset(offset).
		Order("created_at desc").
		Find(&mapHistories)

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
		"mapHistories": mapHistories,
		"total": total,
	}, "查询成功")

}

// @title    DeleteFileHistory
// @description   提供文件的操作记录
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteFileHistory(ctx *gin.Context){
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查看历史操作记录
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	db := common.GetDB().Table("file_histories")

	// 时间格式定义
	layout := "2006-01-02T15:04:05"

	// 读取参数请求
	start := ctx.Query("start")

	if start != "" {
		start, err := time.Parse(layout, start)
		if err == nil{
			db = db.Where("created_at >= ?", start)
		} else {
			response.Fail(ctx, nil, "错误的数据日志开始时间")
			return
		}
	}

	end := ctx.Query("end")

	if end != "" {
		end, err := time.Parse(layout, end)
		if err == nil{
			db = db.Where("created_at <= ?", end)
		} else {
			response.Fail(ctx, nil, "错误的数据日志结束时间")
			return
		}
	}

	result := db.Delete(nil)

	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败: "+result.Error.Error())
		return
	}

	response.Success(ctx, gin.H{"num": result.RowsAffected}, "删除成功")
}

// @title    DeleteDataHistory
// @description   删除数据的操作记录
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteDataHistory(ctx *gin.Context){
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查看历史操作记录
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	db := common.GetDB().Table("data_histories")

	// 时间格式定义
	layout := "2006-01-02T15:04:05"

	// 读取参数请求
	start := ctx.Query("start")

	if start != "" {
		start, err := time.Parse(layout, start)
		if err == nil{
			db = db.Where("created_at >= ?", start)
		} else {
			response.Fail(ctx, nil, "错误的数据日志开始时间")
			return
		}
	}

	end := ctx.Query("end")

	if end != "" {
		end, err := time.Parse(layout, end)
		if err == nil{
			db = db.Where("created_at <= ?", end)
		} else {
			response.Fail(ctx, nil, "错误的数据日志结束时间")
			return
		}
	}

	result := db.Delete(nil)

	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败: "+result.Error.Error())
		return
	}

	response.Success(ctx, gin.H{"num": result.RowsAffected}, "删除成功")
}

// @title    DeleteMapHistory
// @description   删除映射的操作记录
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteMapHistory(ctx *gin.Context){
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查看历史操作记录
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	db := common.GetDB().Table("map_histories")

	// 时间格式定义
	layout := "2006-01-02T15:04:05"

	// 读取参数请求
	start := ctx.Query("start")

	if start != "" {
		start, err := time.Parse(layout, start)
		if err == nil{
			db = db.Where("created_at >= ?", start)
		} else {
			response.Fail(ctx, nil, "错误的数据日志开始时间")
			return
		}
	}

	end := ctx.Query("end")

	if end != "" {
		end, err := time.Parse(layout, end)
		if err == nil{
			db = db.Where("created_at <= ?", end)
		} else {
			response.Fail(ctx, nil, "错误的数据日志结束时间")
			return
		}
	}

	result := db.Delete(nil)

	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败: "+result.Error.Error())
		return
	}

	response.Success(ctx, gin.H{"num": result.RowsAffected}, "删除成功")
}
// @Title  HistoryController
// @Description  该文件用于提供获取历史操作记录的各种函数
package controller

import (
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/response"
	"lianjiang/util"
	"lianjiang/dto"
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

	// 读取参数请求
	start := ctx.Params.ByName("start")

	if start != "" && start != "null" {
		start, err := time.Parse(util.ReadableTimeFormat, start)
		if err == nil{
			db = db.Where("created_at >= ?", start)
		} else {
			response.Fail(ctx, nil, "错误的文件日志开始时间")
			return
		}
	}

	end := ctx.Params.ByName("end")

	if end != "" && end != "null" {
		end, err := time.Parse(util.ReadableTimeFormat, end)
		if err == nil{
			db = db.Where("created_at <= ?", end)
		} else {
			response.Fail(ctx, nil, "错误的数据日志结束时间")
			return
		}
	}

	cond1 := map[string]interface{}{
		"user_id": ctx.DefaultQuery("id", ""),
		"system":  ctx.DefaultQuery("system", ""),
		"file_type":  ctx.DefaultQuery("fileType", ""),
		"option":  ctx.DefaultQuery("option", ""),
	} 

	db = util.DbConditionsEqual(db ,cond1)

	cond2 := map[string]interface{}{
		"file_name":  ctx.DefaultQuery("fileName", ""),
	} 

	db = util.DbConditionsLike(db ,cond2)

	// 查询总数
	var total int64
	dbCount := db.Session(&gorm.Session{})
	dbCount.Count(&total)

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

	// 读取参数请求
	start := ctx.Params.ByName("start")

	if start != "" && start != "null" {
		start, err := time.Parse(util.ReadableTimeFormat, start)
		if err == nil{
			db = db.Where("created_at >= ?", start)
		} else {
			response.Fail(ctx, nil, "错误的数据日志开始时间")
			return
		}
	}

	end := ctx.Params.ByName("end")

	if end != "" && end != "null" {
		end, err := time.Parse(util.ReadableTimeFormat, end)
		if err == nil{
			db = db.Where("created_at <= ?", end)
		} else {
			response.Fail(ctx, nil, "错误的数据日志结束时间")
			return
		}
	}

	cond1 := map[string]interface{}{
		"user_id": ctx.DefaultQuery("id", ""),
		"option":  ctx.DefaultQuery("option", ""),
		"system":  ctx.DefaultQuery("system", ""),
	} 

	db = util.DbConditionsEqual(db ,cond1)

	cond2 := map[string]interface{}{
		"station_name":  ctx.DefaultQuery("station_name", ""),
	} 

	db = util.DbConditionsLike(db ,cond2)

	// 查询总数
	var total int64
	dbCount := db.Session(&gorm.Session{})
	dbCount.Count(&total)

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

	db := common.GetDB().
		Table("map_histories").
		Select(`
		map_histories.created_at,
		map_histories.id,
		map_histories.user_id,
		map_histories.key,
		map_histories.value,
		map_histories.option,
		map_versions.version_name AS version_name
	`).
	Joins("LEFT JOIN map_versions ON map_histories.ver_id = map_versions.id")
	

	var mapHistories []dto.MapHistoryWithVerName

	// 读取参数请求
	start := ctx.Params.ByName("start")

	if start != "" && start != "null" {
		start, err := time.Parse(util.ReadableTimeFormat, start)
		if err == nil{
			db = db.Where("created_at >= ?", start)
		} else {
			response.Fail(ctx, nil, "错误的映射日志开始时间")
			return
		}
	}

	end := ctx.Params.ByName("end")

	if end != "" && end != "null" {
		end, err := time.Parse(util.ReadableTimeFormat, end)
		if err == nil{
			db = db.Where("created_at <= ?", end)
		} else {
			response.Fail(ctx, nil, "错误的映射日志结束时间")
			return
		}
	}

	cond1 := map[string]interface{}{
		"user_id": ctx.DefaultQuery("userId", ""),
		"option":  ctx.DefaultQuery("option", ""),
		"system":  ctx.DefaultQuery("system", ""),
	} 

	db = util.DbConditionsEqual(db ,cond1)

	cond2 := map[string]interface{}{
		"version_name": ctx.DefaultQuery("version_name", ""),
		"key":  ctx.DefaultQuery("key", ""),
	} 

	db = util.DbConditionsLike(db ,cond2)

	// 查询总数
	var total int64
	dbCount := db.Session(&gorm.Session{})
	dbCount.Count(&total)

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
	
	// 返回分页数据
	response.Success(ctx, gin.H{
		"mapHistories": mapHistories,
		"total": total,
	}, "查询成功")

}

// @title    DeleteFileHistory
// @description   删除一组文件的操作记录
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

	// 解析请求体
	var resq dto.BatchDeleteId
	if err := ctx.ShouldBindJSON(&resq); err != nil {
		response.Fail(ctx, nil, "请求体参数解析错误")
		return
	}

	// 开启事务
	db := common.GetDB()
	tx := db.Begin()
	
	result := tx.Where("id IN ?", resq.Ids).Delete(&model.FileHistory{})

	if result.Error != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "删除失败，已撤销")
		return
	}
	
	tx.Commit()

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

	// 解析请求体
	var resq dto.BatchDeleteId
	if err := ctx.ShouldBindJSON(&resq); err != nil {
		response.Fail(ctx, nil, "请求体参数解析错误")
		return
	}

	// 开启事务
	db := common.GetDB()
	tx := db.Begin()
	
	result := tx.Where("id IN ?", resq.Ids).Delete(&model.DataHistory{})

	if result.Error != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "删除失败，已撤销")
		return
	}
	
	tx.Commit()

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

	// 解析请求体
	var resq dto.BatchDeleteId
	if err := ctx.ShouldBindJSON(&resq); err != nil {
		response.Fail(ctx, nil, "请求体参数解析错误")
		return
	}

	// 开启事务
	db := common.GetDB()
	tx := db.Begin()
	
	result := tx.Where("id IN ?", resq.Ids).Delete(&model.MapHistory{})

	if result.Error != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "删除失败，已撤销")
		return
	}
	
	tx.Commit()

	response.Success(ctx, gin.H{"num": result.RowsAffected}, "删除成功")
}
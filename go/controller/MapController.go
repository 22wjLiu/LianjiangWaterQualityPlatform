// @Title  MapController
// @Description  该文件用于提供操作映射表的各种函数
package controller

import (
	"fmt"
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/response"
	"lianjiang/util"
	"strconv"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// @title    ShowMapTables
// @description   查询映射类型
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowMapTables(ctx *gin.Context) {

	var tables []string
	for t := range util.MapMap {
    tables = append(tables, t)
	}
	
	if len(tables) == 0 {
    response.Fail(ctx, nil, "无映射类型")
	}
	// 返回映射类型
	response.Success(ctx, gin.H{"tables": tables}, "请求成功")
}

// @title    ShowMapKeys
// @description   用户查看映射表主键
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowMapKeys(ctx *gin.Context) {

	// TODO 获取path中的id
	Id := ctx.Params.ByName("id")

	M, ok := util.MapMap[Id]

	if !ok {
		response.Fail(ctx, nil, "不存在该映射表")
		return
	}
	// TODO 返回所有key
	response.Success(ctx, gin.H{"keys": M.Keys()}, "请求成功")
}

// @title    ShowMapValue
// @description   用户查看映射表的对应键的值
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowMapValue(ctx *gin.Context) {

	// TODO 获取path中的id
	Id := ctx.Params.ByName("id")

	key := ctx.Params.ByName("key")

	M, ok := util.MapMap[Id]

	if !ok {
		response.Fail(ctx, nil, "不存在该映射表")
		return
	}

	value, ok := M.Get(key)

	if !ok {
		response.Fail(ctx, nil, "该键值不存在")
		return
	}

	// TODO 返回所有value
	response.Success(ctx, gin.H{"value": value}, "请求成功")
}

// @title    ShowCurrentMaps
// @description   查询当前映射表
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowCurrentMaps(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查询用户信息
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	var mapVer model.MapVersion

	db := common.GetDB()

	if err := db.Where("active = ?", 1).First(&mapVer).Error; err != nil {
		response.Fail(ctx, nil, "当前映射表不存在")
		return
	}

	db = db.Table("map_version_details")

	cond1 := map[string]interface{}{
		"ver_id": mapVer.Id,
		"table": ctx.DefaultQuery("table", ""),
	} 

	db = util.DbConditionsEqual(db ,cond1)

	cond2 := map[string]interface{}{
		"key": ctx.DefaultQuery("key", ""),
		"value":  ctx.DefaultQuery("value", ""),
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

	var maps []model.MapVersionDetail

	result := db.
		Limit(pageSize).
		Offset(offset).
		Find(&maps)

	if result.Error != nil {
		response.Fail(ctx, nil, "参数有误")
		return
	}
	
	// 返回分页数据
	response.Success(ctx, gin.H{
		"maps": maps,
		"total": total,
	}, "查询成功")
}

// @title    CreateMapKey
// @description   通过同名键值创建映射
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func CreateMapKey(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在二级以下的用户不能操作映射表
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取path中的key
	key1 := ctx.Query("key1")

	key2 := ctx.Query("key2")

	M, ok := util.MapMap[id]

	if !ok {
		response.Fail(ctx, nil, "不存在该映射表")
		return
	}

	value, ok := M.Get(key1)

	if !ok {
		response.Fail(ctx, nil, "该键值不存在")
		return
	}

	// TODO 做历史记录
	common.GetDB().Create(&model.MapHistory{
		UserId: user.Id,
		Table:  id,
		Key:    key1,
		Value:  fmt.Sprint(value),
		Option: "创建",
	})

	// TODO 设置值
	M.Set(key2, value)

	// TODO 返回所有value
	response.Success(ctx, nil, "设置成功")
}

// @title    CreateMapValue
// @description   创建映射键值对
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func CreateMapValue(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在二级以下的用户不能操作映射表
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取path中的key
	key := ctx.Params.ByName("key")

	// TODO 取出value
	value := ctx.DefaultQuery("value", "")

	M, ok := util.MapMap[id]

	if !ok {
		response.Fail(ctx, nil, "不存在该映射表")
		return
	}

	// TODO 设置值
	M.Set(key, value)

	// TODO 做历史记录
	common.GetDB().Create(&model.MapHistory{
		UserId: user.Id,
		Table:  id,
		Key:    key,
		Value:  fmt.Sprint(value),
		Option: "创建",
	})

	// TODO 返回所有value
	response.Success(ctx, nil, "设置成功")
}

// @title    DeleteMapKey
// @description   删除映射
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteMapKey(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在二级以下的用户不能操作映射表
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取path中的key
	key := ctx.Params.ByName("key")

	M, ok := util.MapMap[id]

	if !ok {
		response.Fail(ctx, nil, "不存在该映射表")
		return
	}

	value, ok := M.Get(key)

	// TODO 检查键值是否存在
	if !ok {
		response.Fail(ctx, nil, "键值不存在")
		return
	}

	// TODO 做历史记录
	common.GetDB().Create(&model.MapHistory{
		UserId: user.Id,
		Table:  id,
		Key:    key,
		Value:  fmt.Sprint(value),
		Option: "删除",
	})

	// TODO 删除值
	M.Remove(key)

	response.Success(ctx, nil, "删除成功")
}

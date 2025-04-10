// @Title  MapController
// @Description  该文件用于提供操作映射表的各种函数
package controller

import (
	"errors"
	"fmt"
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/response"
	"lianjiang/util"
	"lianjiang/dto"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// @title    ShowMapTables
// @description   查询映射类型
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowMapTables(ctx *gin.Context) {

	var tables []dto.MapTables
	for t := range util.MapMap {
    tables = append(tables, dto.MapTables {
			Value: t,
		})
	}
	
	if len(tables) == 0 {
    response.Fail(ctx, nil, "无映射类型")
	}

	// 返回映射类型
	response.Success(ctx, gin.H{"tables": tables}, "请求成功")
}

// @title    ShowActiveMapInfoByStationName
// @description   根据站名和映射类型获取相关的使用中映射信息
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowActiveMapInfosByStationName(ctx *gin.Context) {
	// 获取站名和映射类型
	mapType := ctx.Params.ByName("mapType")
	stationName := ctx.Params.ByName("stationName")

	db := common.GetDB()

	var tableInfo model.DataTableInfo

	if err := db.
		Table("data_table_infos").
		Where("station_name = ? and active = 1", stationName).
		First(&tableInfo).Error; err != nil {
			response.Fail(ctx, nil, "查询数据表信息出错")
			return
	}
	
	var mapVer model.MapVersion

	if err := db.
		Table("map_versions").
		Where("id = ?", tableInfo.MapVerId).
		First(&mapVer).Error; err != nil {
			response.Fail(ctx, nil, "查询映射表信息出错")
			return
	}

	var mapInfos []dto.MapVersionInfos

	if err := db.
	Table("map_version_details").
	Select("`key`, `value`").
	Where("ver_id = ? and `table` = ?", mapVer.Id, mapType).
	Find(&mapInfos).Error; err != nil {
		response.Fail(ctx, nil, "查询映射表信息出错")
		return
	}

	// 返回所有value
	response.Success(ctx, gin.H{"mapInfos": mapInfos}, "请求成功")
}

// @title    ShowMapVersions
// @description   查询映射版本信息
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowMapVersions(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查询用户信息
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	db := common.GetDB()

	db = db.Table("map_versions")

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
		"active": ctx.DefaultQuery("active", ""),
	} 

	db = util.DbConditionsEqual(db ,cond1)

	cond2 := map[string]interface{}{
		"version_name":  ctx.DefaultQuery("version_name", ""),
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

	var mapVersions []model.MapVersion

	result := db.
		Limit(pageSize).
		Offset(offset).
		Find(&mapVersions)

	if result.Error != nil {
		response.Fail(ctx, nil, "参数有误")
		return
	}
	
	// 返回分页数据
	response.Success(ctx, gin.H{
		"mapVersions": mapVersions,
		"total": total,
	}, "查询成功")
}

// @title    ShowMapInfos
// @description   查询映射信息
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowMapInfos(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查询用户信息
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	mapId := ctx.Params.ByName("id")

	if mapId == "" || mapId == "null" {
		response.Fail(ctx, nil, "参数错误")
		return
	}
	
	temp, err := strconv.ParseUint(mapId, 10, 32)
	if err != nil {
		response.Fail(ctx, nil, "参数错误")
		return
	}
	id := uint(temp)

	var mapVer model.MapVersion

	// 获取数据库指针
	db := common.GetDB()

	err = db.Where("id = ?", id).First(&mapVer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(ctx, nil, "映射表信息不存在")
		return
	} else if err != nil {
		response.Fail(ctx, nil, "查询错误")
		return
	}

	db = db.Model(&model.MapVersionDetail{})

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

	var mapInfos []model.MapVersionDetail

	result := db.
		Limit(pageSize).
		Offset(offset).
		Find(&mapInfos)

	if result.Error != nil {
		response.Fail(ctx, nil, "参数有误")
		return
	}
	
	// 返回分页数据
	response.Success(ctx, gin.H{
		"mapInfos": mapInfos,
		"total": total,
	}, "查询成功")
}

// @title    CreateMapVersion
// @description   创建映射版本
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func CreateMapVersion(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查询用户信息
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	isCopy := ctx.DefaultQuery("isCopy", "true")

	var newMapVersion dto.MapVersionName

	if err := ctx.ShouldBindJSON(&newMapVersion); err != nil {
		response.Fail(ctx, nil, "请求数据格式错误")
		return
	}

	versionName := newMapVersion.VersionName
	if versionName == "" {
		response.Fail(ctx, nil, "版本名不能为空")
		return
	}

	// 获取数据库指针
	db := common.GetDB()

	var mapVer model.MapVersion

	err := db.
		Model(&model.MapVersion{}).
		Where("version_name = ?", versionName).
		First(&mapVer).
		Error

	if err == nil {
		response.Fail(ctx, nil, "版本名重复")
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(ctx, nil, "创建失败")
		return
	}

	var mapDetails []model.MapVersionDetail

	if isCopy == "true" {
		var activeVersion model.MapVersion

		err = db.
			Where("active = ?", 1).
			First(&activeVersion).
			Error
		
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(ctx, nil, "没有找到使用中的映射版本信息")
			return
		} else if err != nil {
			response.Fail(ctx, nil, "创建失败")
			return
		}

		err = db.
			Model(&model.MapVersionDetail{}).
			Select("`table`, `key`, `value`").
			Where("ver_id = ?", activeVersion.Id).
			Find(&mapDetails).
			Error
		
		if err != nil {
			response.Fail(ctx, nil, "创建失败")
			return
		}
	}

	mapVer.VersionName = versionName
	mapVer.Active = false;

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.
		Model(&model.MapVersion{}).
		Create(&mapVer).
		Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "创建失败")
		return
	}

	if len(mapDetails) > 0 {
		for i, _ := range mapDetails {
			mapDetails[i].VerId = mapVer.Id
		}
	
		if err := tx.
			Model(&model.MapVersionDetail{}).
			Create(&mapDetails).
			Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "创建失败")
			return
		}
	}

	if err := tx.
		Create(&model.MapHistory{
			UserId: user.Id,
			VerName: versionName,
			Option: "创建(版本)",
		}).
		Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "创建失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "创建失败")
		return
	}

	response.Success(ctx, nil, "创建成功")
}

// @title    DeleteMapVersion
// @description   删除映射版本
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteMapVersion(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查询用户信息
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

	db := common.GetDB()

	// 查找相应表信息
	var tableInfos []model.DataTableInfo

	if err := db.
		Model(&model.DataTableInfo{}).
		Where("map_ver_id in ?", resq.Ids).
		Find(&tableInfos).
		Error; err != nil {
		response.Fail(ctx, nil, "查找相关表信息出错")
		return
	}

	// 查找相应版本名
	var mapVerNames []dto.MapVersionName

	if err := db.
		Model(&model.MapVersion{}).
		Select("version_name").
		Where("id in ?", resq.Ids).
		Find(&mapVerNames).
		Error; err != nil {
		response.Fail(ctx, nil, "查找相关表信息出错")
		return
	}

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, info := range tableInfos {
		if err := tx.Exec("DROP TABLE IF EXISTS " + info.DataTableName).Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "删除失败")
			return
		}
	}
	
	if err := tx.
		Model(&model.MapVersion{}).
		Where("id in ?", resq.Ids).
		Delete(nil).
		Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "删除失败")
			return
	}

	for _, name := range mapVerNames {
		if err := tx.
		Create(&model.MapHistory{
			UserId: user.Id,
			VerName: name.VersionName,
			Option: "删除(版本)",
		}).
		Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "删除失败")
			return
		}
	} 

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// @title    ChangeMapVersion
// @description   切换映射版本
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ChangeMapVersion(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查询用户信息
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	mapId := ctx.DefaultQuery("id", "")

	if mapId == "" {
		response.Fail(ctx, nil, "参数错误")
		return
	}
	
	temp, err := strconv.ParseUint(mapId, 10, 32)
	if err != nil {
		response.Fail(ctx, nil, "参数错误")
		return
	}
	id := uint(temp)

	// 获取数据库指针
	db := common.GetDB()

	var curMapVer model.MapVersion

	err = db.
		Model(&model.MapVersion{}).
		Where("active = 1").
		First(&curMapVer).
		Error; 

	if errors.Is(err, gorm.ErrRecordNotFound){
		response.Fail(ctx, nil, "严重错误：使用中的映射版本不存在")
		return
	} else if err != nil {
		response.Fail(ctx, nil, "查询使用中映射版本失败")
		return
	}

	if curMapVer.Id == id {
		response.Fail(ctx, nil, "当前已在使用改映射版本")
		return
	}

	var destMapVer model.MapVersion

	err = db.
		Model(&model.MapVersion{}).
		Where("id = ?", id).
		First(&destMapVer).
		Error; 

	if errors.Is(err, gorm.ErrRecordNotFound){
		response.Fail(ctx, nil, "严重错误：要切换的映射版本不存在")
		return
	} else if err != nil {
		response.Fail(ctx, nil, "查询要切换映射版本失败")
		return
	}

	tx := db.Begin()

	if err := tx.
		Model(&model.MapVersion{}).
		Where("id = ?", curMapVer.Id).
		Update("active", 0).
		Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "切换失败")
			return
	}

	if err := tx.
		Model(&model.MapVersion{}).
		Where("id = ?", id).
		Update("active", 1).
		Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "切换失败")
			return
	}

	if err := tx.
	Create(&model.MapHistory{
		UserId: user.Id,
		VerName: destMapVer.VersionName,
		Option: "切换版本(后)",
	}).
	Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "删除失败")
		return
	}

	if err := tx.
		Create(&model.MapHistory{
			UserId: user.Id,
			VerName: curMapVer.VersionName,
			Option: "切换版本(前)",
		}).
		Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "删除失败")
			return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "切换失败")
		return
	}

	response.Success(ctx, nil, "切换成功")
}

// @title    CreateMap
// @description   创建映射
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func CreateMap(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查询用户信息
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	mapId := ctx.Params.ByName("id")

	if mapId == "" || mapId == "null" {
		response.Fail(ctx, nil, "参数错误")
		return
	}
	
	temp, err := strconv.ParseUint(mapId, 10, 32)
	if err != nil {
		response.Fail(ctx, nil, "参数错误")
		return
	}
	id := uint(temp)

	var data dto.CreateMapDetail
	if err := ctx.ShouldBindJSON(&data); err != nil {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	if data.Table == "" || data.Key == "" || data.Value == "" {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	isMutiLineMap := false

	if data.Table == "行字段一对多映射" {
		isMutiLineMap = true
	}

	if isMutiLineMap && data.Formula == "" {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	// 获取数据库指针
	db := common.GetDB()

	var curDetail model.MapVersionDetail

	err = db.
		Model(&model.MapVersionDetail{}).
		Where("ver_id = ? and `table` = ? and `key` = ? and `value` = ?", id, data.Table, data.Key, data.Value).
		First(&curDetail).
		Error;

	if err == nil {
		response.Fail(ctx, nil, "添加的映射已存在")
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(ctx, nil, "查询映射信息失败")
		return
	}

	var curMapVer model.MapVersion

	isActiveVersion := true

	err = db.
		Model(&model.MapVersion{}).
		Where("id = ? and active = 1", id).
		First(&curMapVer).
		Error; 

	if errors.Is(err, gorm.ErrRecordNotFound){
		isActiveVersion = false
	} else if err != nil {
		response.Fail(ctx, nil, "查询编辑的映射版本失败")
		return
	}

	if !isActiveVersion {
		err = db.
		Model(&model.MapVersion{}).
		Where("id = ?", id).
		First(&curMapVer).
		Error; 

		if errors.Is(err, gorm.ErrRecordNotFound){
			response.Fail(ctx, nil, "严重错误：编辑的映射版本不存在")
			return
		} else if err != nil {
			response.Fail(ctx, nil, "查询编辑的映射版本失败")
			return
		}
	}

	var tableInfos []model.DataTableInfo

	if isMutiLineMap || data.Table == "列字段映射" {
		if err := db.
		Select("data_table_name").
		Model(&model.DataTableInfo{}).
		Where("map_ver_id = ?", id).
		Find(&tableInfos).
		Error; err != nil {
			response.Fail(ctx, nil, "查询编辑的映射数据失败")
			return
		}
	}

	var mapDetails []model.MapVersionDetail

	if isMutiLineMap {
		if err := db.
		Select("value").
		Model(&model.MapVersionDetail{}).
		Where("ver_id = ?", id).
		Find(&mapDetails).
		Error; err != nil {
			response.Fail(ctx, nil, "查询编辑的映射数据失败")
			return
		}
	}

	tx := db.Begin()

	if err := tx.Create(&model.MapVersionDetail{
		Table: data.Table,
		Key: data.Key,
		Value: data.Value,
		VerId: id,
	}).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "创建失败")
		return
	}

	if isMutiLineMap {
		if err := tx.Create(&model.MapVersionDetail{
			Table: "行字段一对多公式映射",
			Key: data.Key,
			Value: data.Formula,
			VerId: id,
		}).Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "创建失败")
			return
		}
	}

	if len(tableInfos) > 0 {
		if isMutiLineMap {
			rowAllTime := []string{"created_at"}
			var rowAllStr []string

			if isActiveVersion {
				for _, key := range util.PointMap.Keys() {
					v,_ := util.PointMap.Get(key)
					rowAllStr = append(rowAllStr, v.(string))
				}
			} else {
				for _, detail := range mapDetails {
					rowAllStr = append(rowAllStr, detail.Value)
				}
			}
		
			rowAllStr = append(rowAllStr, "rowall_formula")

			sql := util.BuildCreateTableSQL_Str_T_FId(data.Value, rowAllTime, rowAllStr, "table_id", "data_table_infos")
			if err := tx.
				Exec(sql).
				Error; err != nil {
				tx.Rollback()
				response.Fail(ctx, nil, "创建失败")
				return
			}
		} else {
			for _, info := range tableInfos {
				sql := "ALTER TABLE `" + info.DataTableName + "` ADD " + data.Key + " VARCHAR(30) DEFAULT '' NOT NULL"
				if err := tx.
					Exec(sql).
					Error; err != nil {
						tx.Rollback()
						response.Fail(ctx, nil, "创建失败")
						return
				}
			}
		}
	}

	if err := tx.
	Create(&model.MapHistory{
		UserId: user.Id,
		VerName: curMapVer.VersionName,
		Table: data.Table,
		Key: data.Key,
		Value: data.Value,
		Option: "创建",
	}).
	Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "创建失败")
		return
	}

	if isMutiLineMap {
		if err := tx.
		Create(&model.MapHistory{
			UserId: user.Id,
			VerName: curMapVer.VersionName,
			Table: "行字段一对多公式映射",
			Key: data.Key,
			Value: data.Formula,
			Option: "创建",
		}).
		Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "创建失败")
			return
		}
	}

	if isActiveVersion {
		val, ok := util.MapMap[data.Table]
		if !ok {
			tx.Rollback()
			response.Fail(ctx, nil, "创建失败")
			return
		}

		val.Set(data.Key, data.Value)

		if isMutiLineMap {
			util.RowAllFormulaMap.Set(data.Key, data.Formula)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "创建失败")
		return
	}

	response.Success(ctx, nil, "创建成功")
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

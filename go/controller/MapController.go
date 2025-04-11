// @Title  MapController
// @Description  该文件用于提供操作映射表的各种函数
package controller

import (
	"errors"
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

	var curMapVer model.MapVersion

	err := db.
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

	for _, i := range resq.Ids {
		if curMapVer.Id == i {
			response.Fail(ctx, nil, "不能删除使用中的映射版本")
			return
		}
	}

	// 查找相应表信息
	var tableInfos []model.DataTableInfo

	if err := db.
		Model(&model.DataTableInfo{}).
		Select("DISTINCT data_table_name").
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
	var mapMutiLineDetails []model.MapVersionDetail

	if data.Table == "列字段映射" {
		if err := db.
		Select("DISTINCT data_table_name").
		Model(&model.DataTableInfo{}).
		Where("map_ver_id = ?", id).
		Find(&tableInfos).
		Error; err != nil {
			response.Fail(ctx, nil, "查询编辑的映射数据失败")
			return
		}

		if err := db.
		Model(&model.MapVersionDetail{}).
		Where("ver_id = ? and `table` = ?", id, "行字段一对多映射").
		Find(&mapMutiLineDetails).
		Error; err != nil {
			response.Fail(ctx, nil, "查询编辑的映射信息失败")
			return
		}
	}

	var mapDetails []model.MapVersionDetail

	if isMutiLineMap {
		if err := db.
		Select("value").
		Model(&model.MapVersionDetail{}).
		Where("ver_id = ? and `table` = ?", id, "列字段映射").
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
	}

	if len(tableInfos) > 0 {
		for _, info := range tableInfos {
			sql := "ALTER TABLE `" + info.DataTableName + "` ADD `" + data.Value + "` VARCHAR(30) DEFAULT '' NOT NULL"
			if err := tx.
				Exec(sql).
				Error; err != nil {
					tx.Rollback()
					response.Fail(ctx, nil, "创建失败")
					return
			}
		}

		for _, item := range mapMutiLineDetails {
			sql := "ALTER TABLE `" + item.Value + "_" + strconv.FormatUint(uint64(id), 10) + "` ADD `" + data.Value + "` VARCHAR(30) DEFAULT '' NOT NULL"
			if err := tx.
				Exec(sql).
				Error; err != nil {
					tx.Rollback()
					response.Fail(ctx, nil, "创建失败")
					return
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

	// 如果是使用中版本就修改内存中的数据
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

// @title    DeleteMap
// @description   删除映射
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteMap(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查看历史操作记录
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// 获取版本ID
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

	// 解析请求体
	var resq dto.BatchDeleteId
	if err := ctx.ShouldBindJSON(&resq); err != nil {
		response.Fail(ctx, nil, "请求体参数解析错误")
		return
	}

	if len(resq.Ids) <= 0 {
		response.Fail(ctx, nil, "删除成功")
		return
	}

	// 获取数据库指针
	db := common.GetDB()

	// 获取当前映射版本信息
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

	// 获取映射信息
	var mapDetails []model.MapVersionDetail

	if err := db.
		Model(&model.MapVersionDetail{}).
		Where("ver_id = ? and id in ?", id, resq.Ids).
		Find(&mapDetails).
		Error; err != nil {
			response.Fail(ctx, nil, "查询映射信息错误")
			return
	}

	if len(mapDetails) != len(resq.Ids) {
		response.Fail(ctx, nil, "查询部分映射信息错误")
		return
	}

	// 获取特殊映射数量
	var allSingleColumnMapCount int64

	if err := db.
		Model(&model.MapVersionDetail{}).
		Where("ver_id = ? and `table` = ?", id, "列字段映射").
		Count(&allSingleColumnMapCount).
		Error; err != nil {
			response.Fail(ctx, nil, "查询映射信息错误")
			return
	}

	// 筛选特殊映射类型
	var singleColumnMaps []model.MapVersionDetail
	var mutiLineMaps []model.MapVersionDetail
	var remainMutiLineMaps []model.MapVersionDetail
	var formulaMaps	[]model.MapVersionDetail

	for _, detail := range mapDetails {
		if detail.Table == "列字段映射" {
			singleColumnMaps = append(singleColumnMaps, detail)
		} else if detail.Table == "行字段一对多映射" {
			mutiLineMaps = append(mutiLineMaps, detail)
		} else if detail.Table == "行字段一对多公式映射" {
			formulaMaps = append(formulaMaps, detail)
		}
	}

	if allSingleColumnMapCount - int64(len(singleColumnMaps)) <= 0 {
		response.Fail(ctx, nil, "需要至少保留一个列字段映射")
		return
	}

	if len(mutiLineMaps) != len(formulaMaps) {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	formulaMapKeys := make(map[string]struct{})
	for _, j := range formulaMaps {
		formulaMapKeys[j.Key] = struct{}{}
	}
	
	for _, i := range mutiLineMaps {
		if _, ok := formulaMapKeys[i.Key]; !ok {
			response.Fail(ctx, nil, "参数错误")
			return
		}
	}

	// 获取当前版本所有一对多行字段映射
	var allMutiLineMaps []model.MapVersionDetail

	if len(mutiLineMaps) > 0 {
		if err := db.
		Model(&model.MapVersionDetail{}).
		Where("ver_id = ? and `table` = ?", id, "行字段一对多映射").
		Find(&allMutiLineMaps).
		Error; err != nil {
			response.Fail(ctx, nil, "查询映射信息错误")
			return
		}
	
		if len(allMutiLineMaps) - len(mutiLineMaps) < 0 {
			response.Fail(ctx, nil, "参数错误")
			return
		}
	
		mutiLineMapKeys := make(map[string]struct{})
		for _, j := range mutiLineMaps {
			mutiLineMapKeys[j.Key] = struct{}{}
		}
	
		for _, i := range allMutiLineMaps {
			if _, ok := mutiLineMapKeys[i.Key]; !ok {
				remainMutiLineMaps = append(remainMutiLineMaps, i)
			}
		}
	}

	// 查询数据表信息
	var tableInfos []model.DataTableInfo

	if err := db.
		Select("DISTINCT data_table_name").
		Model(&model.DataTableInfo{}).
		Where("map_ver_id = ?", id).
		Find(&tableInfos).
		Error; err != nil {
			response.Fail(ctx, nil, "查询数据表信息错误")
			return
	}

	// 开启事务
	tx := db.Begin()

	if len(tableInfos) > 0 && len(singleColumnMaps) > 0 {
		for _, info := range tableInfos {
			sql := "ALTER TABLE " + info.DataTableName

			for _, item := range singleColumnMaps {
				sql = sql + " DROP COLUMN `" + item.Value + "`,"
			}

			sql = sql[:len(sql) - 1]

			if err := tx.Exec(sql).Error; err != nil {
					tx.Rollback()
					response.Fail(ctx, nil, "删除失败")
					return
			}
		}

		if len(remainMutiLineMaps) > 0 {
			for _, i := range remainMutiLineMaps {
				sql := "ALTER TABLE " + i.Value + "_" + strconv.FormatUint(uint64(id), 10)

				for _, item := range singleColumnMaps {
					sql = sql + " DROP COLUMN `" + item.Value + "`,"
				}

				sql = sql[:len(sql) - 1]

				if err := tx.Exec(sql).Error; err != nil {
						tx.Rollback()
						response.Fail(ctx, nil, "删除失败")
						return
				}
			}
		}
	}

	if len(mutiLineMaps) > 0 {
		for _, item := range mutiLineMaps {
			sql := "DROP TABLE IF EXISTS " + item.Value + "_" + strconv.FormatUint(uint64(id), 10)
			if err := tx.Exec(sql).Error; err != nil {
				tx.Rollback()
				response.Fail(ctx, nil, "删除失败")
				return
			}
		}
	}

	// 删除所选映射详情
	if err := tx.Delete(&mapDetails).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "删除失败")
		return
	}

	// 创建删除记录
	for _, item := range mapDetails {
		if err := tx.
		Create(&model.MapHistory{
			UserId: user.Id,
			VerName: curMapVer.VersionName,
			Table: item.Table,
			Key: item.Key,
			Value: item.Value,
			Option: "删除",
		}).
		Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "创建失败")
			return
		}
	}

	// 如果是使用中版本就修改内存中的数据
	if isActiveVersion {
		for _, item := range mapDetails {
			val, ok := util.MapMap[item.Table]
			if !ok {
				tx.Rollback()
				response.Fail(ctx, nil, "删除失败")
				return
			}
	
			val.Remove(item.Key)

			if _, ok := val.Get(item.Key); ok {
				tx.Rollback()
				response.Fail(ctx, nil, "删除失败")
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// @title    UpdateMap
// @description   更新映射
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func UpdateMap(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查询用户信息
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	mapId := ctx.Params.ByName("id")
	cmapId := ctx.Params.ByName("curMapId")

	if mapId == "" || mapId == "null" || cmapId == "" || cmapId == "null" {
		response.Fail(ctx, nil, "参数错误")
		return
	}
	
	temp, err := strconv.ParseUint(mapId, 10, 32)
	if err != nil {
		response.Fail(ctx, nil, "参数错误")
		return
	}
	id := uint(temp)

	temp, err = strconv.ParseUint(cmapId, 10, 32)
	if err != nil {
		response.Fail(ctx, nil, "参数错误")
		return
	}
	curMapId := uint(temp)

	var data dto.UpdateMapDetail
	if err := ctx.ShouldBindJSON(&data); err != nil {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	if data.Table == "" || data.Key == "" || data.Value == "" {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	// 获取数据库指针
	db := common.GetDB()

	var curDetail model.MapVersionDetail

	err = db.
		Model(&model.MapVersionDetail{}).
		Where("id = ?", curMapId).
		First(&curDetail).
		Error;

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(ctx, nil, "编辑的映射不存在")
		return
	} else if err != nil {
		response.Fail(ctx, nil, "查询映射信息失败")
		return
	}

	isValueChange := false

	if curDetail.Value != data.Value {
		isValueChange = true
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
	var mapDetails []model.MapVersionDetail

	if isValueChange && data.Table == "列字段映射" {
		if err := db.
		Select("DISTINCT data_table_name").
		Model(&model.DataTableInfo{}).
		Where("map_ver_id = ?", id).
		Find(&tableInfos).
		Error; err != nil {
			response.Fail(ctx, nil, "查询编辑的映射数据失败")
			return
		}

		if err := db.
		Model(&model.MapVersionDetail{}).
		Where("ver_id = ? and `table` = ?", id, "行字段一对多映射").
		Find(&mapDetails).
		Error; err != nil {
			response.Fail(ctx, nil, "查询编辑的映射信息失败")
			return
		}
	}

	// 开启事务
	tx := db.Begin()

	if isValueChange {
		if data.Table == "列字段映射" {
			for _, info := range tableInfos {
				sql := "ALTER TABLE `" + info.DataTableName + "` CHANGE `" + curDetail.Value + "` `" +  data.Value
				sql = sql + "` VARCHAR(30) DEFAULT '' NOT NULL"
				if err := tx.
					Exec(sql).
					Error; err != nil {
						tx.Rollback()
						response.Fail(ctx, nil, "更新失败")
						return
				}
			}

			for _, detail := range mapDetails {
				sql := "ALTER TABLE " + detail.Value + "_" + strconv.FormatUint(uint64(id), 10) + " CHANGE `" + curDetail.Value + "` `" +  data.Value
				sql = sql + "` VARCHAR(30) DEFAULT '' NOT NULL"
				if err := tx.
					Exec(sql).
					Error; err != nil {
						tx.Rollback()
						response.Fail(ctx, nil, "更新失败")
						return
				}
			}
		} else if data.Table == "行字段一对多映射" {
			sql := "RENAME TABLE " + curDetail.Value + "_" + strconv.FormatUint(uint64(id), 10) + " TO " + data.Value + "_" + strconv.FormatUint(uint64(id), 10)
			if err := tx.
				Exec(sql).
				Error; err != nil {
					tx.Rollback()
					response.Fail(ctx, nil, "更新失败")
					return
			}
		}
	}

	if err := tx.
	Model(&model.MapVersionDetail{}).
	Where("id = ?", curDetail.Id).
	Updates(&data).
	Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "更新失败")
		return
	}

	if err := tx.
	Create(&model.MapHistory{
		UserId: user.Id,
		VerName: curMapVer.VersionName,
		Table: data.Table,
		Key: data.Key,
		Value: data.Value,
		Option: "更新(后)",
	}).
	Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "更新失败")
		return
	}

	if err := tx.
	Create(&model.MapHistory{
		UserId: user.Id,
		VerName: curMapVer.VersionName,
		Table: curDetail.Table,
		Key: curDetail.Key,
		Value: curDetail.Value,
		Option: "更新(前)",
	}).
	Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "更新失败")
		return
	}

	// 如果是使用中版本就修改内存中的数据
	if isActiveVersion {
		val, ok := util.MapMap[data.Table]
		if !ok {
			tx.Rollback()
			response.Fail(ctx, nil, "更新失败")
			return
		}
		if curDetail.Key != data.Key {
			val.Remove(curDetail.Key)
		}
		val.Set(data.Key, data.Value)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, nil, "更新成功")
}
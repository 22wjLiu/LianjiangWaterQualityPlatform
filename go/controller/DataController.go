// @Title  DataController
// @Description  该文件用于提供操作数据的各种函数
package controller

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os"
	"os/exec"
	"errors"
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/util"
	"lianjiang/dto"
	"log"
	"time"
	"strconv"

	"lianjiang/response"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// // @title    ShowStationsWhichHasData
// // @description   查询当前的有数据的站点
// // @param    ctx *gin.Context       接收一个上下文
// // @return   void
func ShowStationsWhichHasData(ctx *gin.Context) {
	db := common.GetDB()

	var stationNames []dto.StationNameData

	if err := db.
		Table("data_table_infos").
		Select("DISTINCT station_name, system").
		Where("active = 1").
		Find(&stationNames).Error; err != nil {
		response.Fail(ctx, nil, "获取站名失败")
		return
	}

	// 返回映射类型
	response.Success(ctx, gin.H{"stationNames": stationNames}, "请求成功")
}

// // @title    ShowDataTimeRange
// // @description   获取点集数据时间范围
// // @param    ctx *gin.Context       接收一个上下文
// // @return   void
func ShowDataTimeRange(ctx *gin.Context) {
	// 获取站名
	name := ctx.Params.ByName("name")

	_, ok := util.StationMap.Get(name)

	if !ok {
		response.Fail(ctx, nil, "不存在站名"+name)
		return
	}

	// 获取制度
	system := ctx.Params.ByName("system")

	_, ok = util.SysMap.Get(system)

	if !ok {
		response.Fail(ctx, nil, "不存在制度"+system)
		return
	}

	db := common.GetDB()

	var tableInfo model.DataTableInfo
	if err := db.Table("data_table_infos").Where("station_name = ? and system = ? and active = 1", name, system).First(&tableInfo).Error; err != nil {
		response.Fail(ctx, nil, "数据丢失")
		return
	}

	// 查看是否存在该表
	var exists bool
	err := db.Raw(`
	SELECT COUNT(*) > 0 FROM information_schema.tables 
	WHERE table_schema = DATABASE() AND table_name = ?
	`, tableInfo.DataTableName).Scan(&exists).Error

	if err != nil {
		response.Fail(ctx, nil, "数据丢失")
		return
	}

	if !exists {
		response.Fail(ctx, nil, "不存在对应表")
		return
	}

	var minTime, maxTime time.Time

	err = db.Table(tableInfo.DataTableName).
			Select("MIN(time) as min_time, MAX(time) as max_time").
			Row().
			Scan(&minTime, &maxTime)
	
	if err != nil {
		response.Fail(ctx, nil, "查询失败")
		return
	}

	response.Success(ctx, gin.H{
		"minTime": minTime.Format(util.ReadableTimeFormat),
		"maxTime": maxTime.Format(util.ReadableTimeFormat),
	},"查找成功")
}

// @title    ShowData
// @description   前台获取点集数据
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowData(ctx *gin.Context) {

	// 获取站名
	name := ctx.Params.ByName("name")

	if name == "" {
		response.Fail(ctx, nil, "不存在站名数据")
		return
	}

	_, ok := util.StationMap.Get(name)

	if !ok {
		response.Fail(ctx, nil, "不存在站名"+name)
		return
	}

	// 获取制度
	system := ctx.Params.ByName("system")

	_, ok = util.SysMap.Get(system)

	if !ok {
		response.Fail(ctx, nil, "不存在制度"+system)
		return
	}

	db := common.GetDB()

	var tableInfo model.DataTableInfo
	if err := db.Table("data_table_infos").Where("station_name = ? and system = ? and active = 1", name, system).First(&tableInfo).Error; err != nil {
		response.Fail(ctx, nil, "数据丢失")
		return
	}

	// 查看是否存在该表
	var exists bool
	err := db.Raw(`
	SELECT COUNT(*) > 0 FROM information_schema.tables 
	WHERE table_schema = DATABASE() AND table_name = ?
	`, tableInfo.DataTableName).Scan(&exists).Error

	if err != nil {
		response.Fail(ctx, nil, "数据丢失")
		return
	}

	if !exists {
		response.Fail(ctx, nil, "不存在对应表")
		return
	}

	queryDB := common.GetDB().Table(tableInfo.DataTableName)

	startNotNull := false
	endNotNull := false

	// 取出请求
	start := ctx.DefaultQuery("start", "")
	if start != "" && start != "null" {
		startNotNull = true
		s, err := time.Parse(util.ReadableTimeFormat, start)
		if err != nil{
			response.Fail(ctx, nil, "错误的数据开始时间")
			return
		}
		start = s.Format(util.ReadableTimeFormat)
		queryDB = queryDB.Where("time >= ?", start)
	}

	end := ctx.DefaultQuery("end", "")
	if end != "" && end != "null" {
		endNotNull = true
		e, err := time.Parse(util.ReadableTimeFormat, end)
		if err != nil{
			response.Fail(ctx, nil, "错误的数据结束时间")
			return
		}
		end = e.Format(util.ReadableTimeFormat)
		if !startNotNull {
			start = e.AddDate(0, -1, 0).Format(util.ReadableTimeFormat)
		}
		queryDB = queryDB.Where("time <= ?", end)
	}

	var endTime time.Time

	if !endNotNull {
		err = db.Table(tableInfo.DataTableName).Select("MAX(time) as max_time").Row().Scan(&endTime)

		if err != nil {
			response.Fail(ctx, nil, "查询时间失败")
			return
		}

		end = endTime.Format(util.ReadableTimeFormat)

		if !startNotNull {
			start = endTime.AddDate(0, -1, 0).Format(util.ReadableTimeFormat)
		}
	}

	// 查找对应数组
	resultArr := make([]map[string]interface{}, 0)

	queryDB.Table(tableInfo.DataTableName).
		Where("time >= ? and time <= ?", start, end).
		Scan(&resultArr)

	response.Success(ctx, gin.H{
		"resultArr": resultArr,
		"startTime": start,
		"endTime": end,
		}, "查找成功")
}

// @title    ShowFieldsData
// @description   前台获取所有单字段点集数据
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowFieldsData(ctx *gin.Context) {

	// 获取站名
	name := ctx.Params.ByName("name")

	if name == "" {
		response.Fail(ctx, nil, "不存在站名数据")
		return
	}

	_, ok := util.StationMap.Get(name)

	if !ok {
		response.Fail(ctx, nil, "不存在站名"+name)
		return
	}

	// 获取制度
	system := ctx.Params.ByName("system")

	_, ok = util.SysMap.Get(system)

	if !ok {
		response.Fail(ctx, nil, "不存在制度"+system)
		return
	}

	db := common.GetDB()

	var tableInfo model.DataTableInfo
	if err := db.Table("data_table_infos").Where("station_name = ? and system = ? and active = 1", name, system).First(&tableInfo).Error; err != nil {
		response.Fail(ctx, nil, "数据丢失")
		return
	}

	var mapDetails []model.MapVersionDetail
	if err := db.Model(&model.MapVersionDetail{}).Where("ver_id = ? and `table` = ?", tableInfo.MapVerId, "列字段映射").Find(&mapDetails).Error; err != nil {
		response.Fail(ctx, nil, "查找映射信息错误")
		return
	}

	// 查看是否存在该表
	var exists bool
	err := db.Raw(`
	SELECT COUNT(*) > 0 FROM information_schema.tables 
	WHERE table_schema = DATABASE() AND table_name = ?
	`, tableInfo.DataTableName).Scan(&exists).Error

	if err != nil {
		response.Fail(ctx, nil, "数据丢失")
		return
	}

	if !exists {
		response.Fail(ctx, nil, "不存在对应表")
		return
	}

	queryDB := common.GetDB().Table(tableInfo.DataTableName)

	startNotNull := false
	endNotNull := false

	// 取出请求
	start := ctx.DefaultQuery("start", "")
	if start != "" && start != "null" {
		startNotNull = true
		s, err := time.Parse(util.ReadableTimeFormat, start)
		if err != nil{
			response.Fail(ctx, nil, "错误的数据开始时间")
			return
		}
		start = s.Format(util.ReadableTimeFormat)
		queryDB = queryDB.Where("time >= ?", start)
	}

	end := ctx.DefaultQuery("end", "")
	if end != "" && end != "null" {
		endNotNull = true
		e, err := time.Parse(util.ReadableTimeFormat, end)
		if err != nil{
			response.Fail(ctx, nil, "错误的数据结束时间")
			return
		}
		end = e.Format(util.ReadableTimeFormat)
		if !startNotNull {
			start = e.AddDate(0, -1, 0).Format(util.ReadableTimeFormat)
		}
		queryDB = queryDB.Where("time <= ?", end)
	}

	var endTime time.Time

	if !endNotNull {
		err = db.Table(tableInfo.DataTableName).Select("MAX(time) as max_time").Row().Scan(&endTime)

		if err != nil {
			response.Fail(ctx, nil, "查询时间失败")
			return
		}

		end = endTime.Format(util.ReadableTimeFormat)

		if !startNotNull {
			start = endTime.AddDate(0, -1, 0).Format(util.ReadableTimeFormat)
		}
	}

	// 查找对应数组
	resultArr := make([][]map[string]interface{}, 0)
	var mapInfos []dto.MapVersionInfos

	queryDB = queryDB.Where("time >= ? and time <= ?", start, end)

	for _, item := range mapDetails {
		if item.Value == "time"{
			continue
		}
		results := make([]map[string]interface{}, 0)
		if err := queryDB.
		Select("time, `" + item.Value + "`").
		Scan(&results).
		Error; err != nil {
			response.Fail(ctx, nil, "查询数据失败")
			return
		}
		resultArr = append(resultArr, results)
		mapInfos = append(mapInfos, dto.MapVersionInfos{
			Key: item.Key,
			Value: item.Value,
		})
	}

	response.Success(ctx, gin.H{
		"resultArr": resultArr,
		"mapInfos": mapInfos,
		"startTime": start,
		"endTime": end,
		}, "查找成功")
}

// @title    ShowStationMapSystem
// @description   获取有数据站名及其映射版本、制度信息
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowStationMapSystem(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// 安全等级在四级以下的用户不能获取点集数据
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	var tableInfos []map[string]interface{}

	db := common.GetDB()

	if err := db.
	Model(&model.DataTableInfo{}).
	Select("data_table_infos.station_name, data_table_infos.data_table_name, data_table_infos.system, data_table_infos.map_ver_id, map_versions.version_name").
	Joins("LEFT JOIN map_versions ON data_table_infos.map_ver_id = map_versions.id").
	Group("data_table_infos.station_name, data_table_infos.data_table_name, data_table_infos.system, data_table_infos.map_ver_id, map_versions.version_name").
	Scan(&tableInfos).Error; err != nil {
	response.Fail(ctx, nil, "查询站点信息出错")
	return
}

	response.Success(ctx, gin.H{"tableInfos": tableInfos}, "查找成功")
}

// @title    ShowDataBackStage
// @description   后台获取点集数据
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowDataBackStage(ctx *gin.Context) {
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// 安全等级在四级以下的用户不能获取点集数据
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// 获取站名
	stationName := ctx.DefaultQuery("station_name", "")
	if stationName == "" {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	// 获取制度
	system := ctx.DefaultQuery("system", "")
	if system == "" {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	// 获取映射版本ID
	mapId := ctx.DefaultQuery("map_ver_id", "")
	if mapId == "" {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	temp, err := strconv.ParseUint(mapId, 10, 32)
	if err != nil {
		response.Fail(ctx, nil, "参数错误")
		return
	}
	mapVerId := uint(temp)

	// 获取数据库指针
	db := common.GetDB()

	var mapVer model.MapVersion
	err = db.
	Model(&model.MapVersion{}).
	Where("id = ?", mapVerId).
	Find(&mapVer).
	Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(ctx, nil, "映射版本不存在")
		return
	} else if err != nil {
		response.Fail(ctx, nil, "查找映射版本信息错误")
		return
	}

	var mapDetails []model.MapVersionDetail

	if err := db.
		Model(&model.MapVersionDetail{}).
		Where("ver_id = ? and `table` = ?", mapVerId, "列字段映射").
		Find(&mapDetails).
		Error; err != nil {
			response.Fail(ctx, nil, "查询映射信息错误")
			return
	}

	var tableInfo model.DataTableInfo

	err = db.
		Model(&model.DataTableInfo{}).
		Where("station_name = ? and map_ver_id = ? and system = ?", stationName, mapVerId, system).
		First(&tableInfo).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(ctx, nil, "站点数据不存在")
		return
	} else if err != nil {
		response.Fail(ctx, nil, "查询站点数据错误")
		return
	}

	db = db.Table(tableInfo.DataTableName)

	// 读取参数请求
	start := ctx.Params.ByName("start")

	if start != "" && start != "null" {
		start, err := time.Parse(util.ReadableTimeFormat, start)
		if err == nil{
			db = db.Where("time >= ?", start)
		} else {
			response.Fail(ctx, nil, "错误的数据开始时间")
			return
		}
	}

	end := ctx.Params.ByName("end")

	if end != "" && end != "null" {
		end, err := time.Parse(util.ReadableTimeFormat, end)

		if err == nil{
			db = db.Where("time <= ?", end)
		} else {
			response.Fail(ctx, nil, "错误的数据结束时间")
			return
		}
	}

	// 查询总数
	var total int64
	db.Count(&total)

	// 获取分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "25"))
		
	offset := (page - 1) * pageSize

	resultArr := make([]map[string]interface{}, 0)

	selectSQl := ""

	for _, item := range mapDetails {
		selectSQl += "`" + item.Value + "` AS '" + item.Key + "' ,"
	}

	selectSQl = selectSQl[: len(selectSQl) - 1]

	subQuery := common.GetDB().Table(tableInfo.DataTableName)

	if start != "null" {
		subQuery = subQuery.Where("time >= ?", start)
	}

	if end != "null" {
		subQuery = subQuery.Where("time <= ?", end)
	}

	subQuery = subQuery.
	Order("time DESC").
	Limit(10000)

	result := common.GetDB().
	Select(selectSQl).
	Table("(?) as t", subQuery).
	Limit(pageSize).
	Offset(offset).
	Scan(&resultArr)
	
	if result.Error != nil {
		log.Printf("查询失败：%v\n", err)
		response.Fail(ctx, nil, "查询失败")
		return
	}
		
	// 返回分页数据
	response.Success(ctx, gin.H{
		"resultArr": resultArr,
		"total": total,
		"versionName": mapVer.VersionName,
		"stationName": stationName,
		"system": system,
		"dataTableName": tableInfo.DataTableName,
	}, "查询成功")
}

// @title    DeleteDataBackStage
// @description   删除后台数据
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteDataBackStage(ctx *gin.Context){
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查看历史操作记录
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// 解析请求体
	var resq dto.BatchDeleteTimeId
	if err := ctx.ShouldBindJSON(&resq); err != nil {
		response.Fail(ctx, nil, "请求体参数解析错误")
		return
	}

	// 开启事务
	db := common.GetDB()
	tx := db.Begin()
	
	result := tx.Table(resq.DataTableName).Where("time IN ?", resq.Times).Delete(nil)

	if result.Error != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "删除失败，已撤销")
		return
	}

	for _, time := range resq.Times {
		if err := tx.
		Model(&model.DataHistory{}).
		Create(&model.DataHistory{
			UserId: 			user.Id,
			Time: 				time,
			StationName: 	resq.StationName,
			System: 			resq.System,
			VersionName:	resq.VersionName,
			Option:				"删除",
		}).Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "删除失败，已撤销")
			return
		}
	}
	
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "删除失败，已撤销")
		return
	}

	response.Success(ctx, gin.H{"num": result.RowsAffected}, "删除成功")
}

// @title    UpdateDataBackStage
// @description   更新后台数据
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func UpdateDataBackStage(ctx *gin.Context){
	// 获取登录用户
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)
	
	// 安全等级在四级以下的用户不能查看历史操作记录
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	mapVerId := ctx.DefaultQuery("mapVerId", "")

	versionName := ctx.DefaultQuery("versionName", "")

	stationName := ctx.DefaultQuery("stationName", "")

	system := ctx.DefaultQuery("system", "")

	if mapVerId == "" || versionName == "" || stationName == "" || system == "" {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	// 解析请求体
	resq := make(map[string]interface{})
	if err := ctx.ShouldBindJSON(&resq); err != nil {
		response.Fail(ctx, nil, "请求体参数解析错误")
		return
	}

	var mapDetails []model.MapVersionDetail

	db := common.GetDB()

	if err := db.
		Model(&model.MapVersionDetail{}).
		Where("ver_id = ? and `table` = ?", mapVerId, "列字段映射").
		Find(&mapDetails).
		Error; err != nil {
			response.Fail(ctx, nil, "查询映射信息错误")
			return
	}
	
	var tableInfo model.DataTableInfo

	err := db.
		Model(&model.DataTableInfo{}).
		Where("station_name = ? and map_ver_id = ? and system = ?", stationName, mapVerId, system).
		First(&tableInfo).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(ctx, nil, "站点数据不存在")
		return
	} else if err != nil {
		response.Fail(ctx, nil, "查询站点数据错误")
		return
	}

	updates := make(map[string]interface{})

	for _, item := range mapDetails {
		updates[item.Value] = resq[item.Key].(string)
	}

	if len(updates) == 0 {
		response.Fail(ctx, nil, "无有效更新字段")
		return
	}

	loc, _ := time.LoadLocation("Asia/Shanghai")
	parsedTime, err := time.ParseInLocation(util.ReadableTimeFormat, resq["时间"].(string), loc)

	log.Println(parsedTime)

	tx := db.Begin()

	if err := tx.Debug().
	Table(tableInfo.DataTableName).
	Where("time = ?", parsedTime).
	Updates(updates).
	Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "更新失败")
		return
	}

	if err := tx.
	Model(&model.DataHistory{}).
	Create(&model.DataHistory{
		UserId: 			user.Id,
		Time: 				model.Time(parsedTime),
		StationName: 	stationName,
		System: 			system,
		VersionName:	versionName,
		Option:				"更新",
	}).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "更新失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, nil, "更新成功")
}

// @title    Forecast
// @description   预测数据
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Forecast(ctx *gin.Context) {
	// 参数解析
	startStr := ctx.Params.ByName("start")
	endStr   := ctx.Params.ByName("end")

	start, err := time.Parse(util.ForecastTimeFormat, startStr)
	if startStr != "" && err != nil {
			response.Fail(ctx, nil, "错误的数据开始时间")
			return
	}
	end, err := time.Parse(util.ForecastTimeFormat, endStr)
	if endStr != "" && err != nil {
			response.Fail(ctx, nil, "错误的数据结束时间")
			return
	}

	stationName := ctx.DefaultQuery("station_name", "")
	system      := ctx.DefaultQuery("system", "")
	field       := ctx.DefaultQuery("field", "")
	if stationName == "" || system == "" || field == "" {
			response.Fail(ctx, nil, "参数错误")
			return
	}

	var freq string

	if system == "小时制" {
		freq = "H"
	} else {
		freq = "M"
	}

	db := common.GetDB()

	var tableInfo model.DataTableInfo
	if err := db.Table("data_table_infos").Where("station_name = ? and system = ? and active = 1", stationName, system).First(&tableInfo).Error; err != nil {
		response.Fail(ctx, nil, "数据丢失")
		return
	}

	resultArr := make([]map[string]interface{}, 0)

	if err := db.
		Table(tableInfo.DataTableName).
		Select("time, " + field).
		Order("time ASC").
		Scan(&resultArr).
		Error; err != nil {
			response.Fail(ctx, nil, "查询错误")
			return
	}

	tmp, err := os.CreateTemp("", "forecast-*.csv")
	if err != nil {
			response.Fail(ctx, nil, "临时文件创建失败")
			return
	}
	defer os.Remove(tmp.Name()) // 用完即删

	w := csv.NewWriter(tmp)
	_ = w.Write([]string{"time", field}) // header
	for _, r := range resultArr {
			// 时间按与 Python 相同的格式输出
			tStr := r["time"].(time.Time).Format(util.ForecastTimeFormat)
			vStr := r[field].(string)
			_ = w.Write([]string{tStr, vStr})
	}
	w.Flush()
	tmp.Close()

	// 调用python
	python := "/home/liu/anaconda3/envs/forecast_env/bin/python"
	script := "./foresee.py"
	cmd := exec.Command(
			python, script,
			tmp.Name(),             																// CSV 文件
			field,                  																// 预测字段
			start.Format(util.ForecastTimeFormat),               		// 开始时间
			end.Format(util.ForecastTimeFormat),                 		// 结束时间
			freq,                   						 										// 频率
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
			log.Println("Python 调用失败: "+err.Error())
			response.Fail(ctx, nil, "预测失败")
			return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
			log.Println("解析预测结果失败")
			response.Fail(ctx, nil, "预测失败")
			return
	}
	response.Success(ctx, gin.H{"result": result}, "预测成功")
}

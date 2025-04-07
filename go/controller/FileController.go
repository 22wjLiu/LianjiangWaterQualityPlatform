// @Title  PointController
// @Description  该文件用于提供操作点集的各种函数
package controller

import (
	"errors"
	"fmt"
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/util"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"strconv"
	"time"

	"lianjiang/response"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// @title    Upload
// @description   opt
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Upload(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// 安全等级在三级以下的用户无法上传文件
	if user.Level < 3 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// 获取上传文件
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Printf("获取上传文件失败: %v", err)
		response.Fail(ctx, nil, "获取上传文件失败")
		return
	}

	// 验证文件格式
	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".xls":  true,
		".xlsx": true,
		".csv":  true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		response.Fail(ctx, nil, fmt.Sprintf("无法处理后缀为 %s 的文件", extName))
		return
	}

	// 从path中获取制度
	system := ctx.Params.ByName("system")

	// 从制度映射表中获取制度映射
	sys, ok := util.SysMap.Get(system)
	if !ok {
		response.Fail(ctx, nil, fmt.Sprintf("%s的制度映射未注册", system))
		return
	}

	// 从标记映射表中获取标记映射
	opt, ok := util.OptMap.Get(sys.(string))
	if !ok {
		response.Fail(ctx, nil, fmt.Sprintf("%s制度的标记映射未注册", system))
		return
	}

	// 尝试建立对应文件夹
	dirPath := filepath.Join("./home/", sys.(string))
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Printf("创建目录失败: %v\n", err)
		response.Fail(ctx, nil, "创建存储目录失败，请重试")
		return
	}

	// 将文件存入本地
	fullPath := filepath.Join(dirPath, file.Filename)
	err = ctx.SaveUploadedFile(file, fullPath)
	if err != nil {
		log.Printf("保存文件失败: %v\n", err)
		response.Fail(ctx, nil, "文件保存失败，无法进行解析，请重试")
		return
	}

	// 读取文件
	res, err := util.Read(fullPath)
	if err != nil || res == nil {
		log.Printf("读取文件失败: %v\n", err)
	
		// 删除临时文件
		removeErr := os.Remove(fullPath)
		if removeErr != nil {
			log.Printf("删除上传文件失败: %v\n", removeErr)
		}
	
		response.Fail(ctx, nil, "读取文件内容失败，请重试")
		return
	}	

	log.Println("开始解析上传文件...")

	index, stName, start, flag, isFind:= util.FindStationNameFromFile(res, opt)

	// 未找到数据段上一行标记
	if !flag {
		response.Fail(ctx, nil, "文件内容缺少数据段上一行标记")
		return
	}

	// 未找到站名
	if !isFind {
		response.Fail(ctx, nil, "未找到站名")
		return
	}

	// 如果站名没有注册
	if !util.StationMap.Has(stName) {
		response.Fail(ctx, nil, fmt.Sprintf("站名%s映射未注册", stName))
		return
	}

	// 获取数据库指针
	db := common.GetDB()

	// 获取当前映射信息用于建表
	var mapVer model.MapVersion

	if err := db.Where("active = ?", 1).First(&mapVer).Error; err != nil{
		log.Println("获取当前映射信息错误")
		response.Fail(ctx, nil, "获取当前映射信息错误")
		return
	}

	st, _ := util.StationMap.Get(stName)

	tableName := sys.(string) + "_" + st.(string) + "_" + fmt.Sprintf("%d", mapVer.Id)

	var tableInfo model.DataTableInfo

	// 是否是新建的主数据表
	isNewDataTable := false

	var exists bool
	err = db.Raw(`
	SELECT COUNT(*) > 0 FROM information_schema.tables 
	WHERE table_schema = DATABASE() AND table_name = ?
	`, tableName).Scan(&exists).Error

	if err != nil {
		log.Printf("查找数据表[%s]失败\n", tableName)
		response.Fail(ctx, nil, "查找数据表失败")
		return
	}
	
	// 查看是否存在表
	if !exists {
		// 不存在即开始建表
		log.Printf("主数据表[%s]不存在，开始建表...\n", tableName)
		isNewDataTable = true
		var rowTime []string
		var rowStr  []string

		for _, v := range index {
			if strings.HasPrefix(v, "time")	{
				rowTime = append(rowTime, v)
			} else {
				rowStr = append(rowStr, v)
			}
		}

		sql := util.BuildCreateTableSQL_Str_T(tableName, rowTime, rowStr)

		if err := db.Exec(sql).Error; err != nil {
			log.Printf("数据表[%s]创建失败\n", tableName)
			response.Fail(ctx, nil, "数据表创建失败")
			return
		}

		log.Printf("主数据表[%s]创建成功\n", tableName)
		log.Printf("开始写入主数据表[%s]信息到表[data_table_infos]中...\n", tableName)

		// 写入数据表信息
		tableInfo.MapVerId = mapVer.Id
		tableInfo.DataTableName = tableName
		tableInfo.StationName = stName
		tableInfo.File = file.Filename
		tableInfo.System = system
		tableInfo.Active = true

		if err := db.Create(&tableInfo).Error; err != nil {
			log.Printf("写入主数据表[%s]信息失败: %v\n", tableName, err)
			response.Fail(ctx, nil, "写入数据表信息失败")
			if err := db.Exec("DROP TABLE IF EXISTS " + tableName).Error; err != nil {
				log.Printf("删除主数据表 [%s] 失败，请检查或手动删除：%v\n", tableName, err)
			}
			return
		}
	} else {
		err := db.Table("data_table_infos").
			Where("data_table_name = ? and map_ver_id = ?", tableName, mapVer.Id).
			First(&tableInfo).
			Error;
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("主数据表[%s]信息不存在\n", tableName)
			response.Fail(ctx, nil, "该主数据表信息不存在")
			return
		} else if err != nil {
			log.Printf("查找主数据表[%s]信息失败: %v\n", tableName, err)
			response.Fail(ctx, nil, "查找主数据表信息失败")
			return
		}

		var tempInfo model.DataTableInfo
		tempInfo.MapVerId = mapVer.Id
		tempInfo.DataTableName = tableName
		tempInfo.StationName = stName
		tempInfo.File = file.Filename
		tempInfo.System = system
		tempInfo.Active = true

		if err := db.Table("data_table_infos").Create(&tempInfo).Error; err != nil {
			log.Printf("写入主数据表[%s]信息失败: %v\n", tableName, err)
			response.Fail(ctx, nil, "写入数据表信息失败")
			return
		}

		tableInfo = tempInfo
	}

	isErrorHappen := false
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ defer 整体发生 panic：%v\n", r)
		}
		if !isErrorHappen {
			return
		}

		log.Println("开始删除脏数据...")

		deleteNoError := true
	
		cleanupByTable := func(table string, tableID uint) {
			err = db.Raw(`
			SELECT COUNT(*) > 0 FROM information_schema.tables 
			WHERE table_schema = DATABASE() AND table_name = ?
			`, tableName).Scan(&exists).Error

			if !exists {
				return
			}

			var count int64
			if err := db.Table(table).Where("table_id = ?", tableID).Count(&count).Error; err != nil {
				log.Printf("查询表 [%s] id=%d 时发生错误：%v\n", table, tableID, err)
				deleteNoError = false
				return
			}

			if count != 0 {
				if err := db.Table(table).Where("table_id = ?", tableID).Delete(nil).Error; err != nil {
					log.Printf("删除表 [%s] id=%d 脏数据失败，请检查或手动删除：%v\n", table, tableID, err)
					deleteNoError = false
				}
			}
		}
	
		// 清理 rowall 数据表中的脏数据
		for item := range util.RowAllMap.Iter() {
			tName := item.Val.(string)
			cleanupByTable(tName, tableInfo.Id)
		}
		
		// 清理 rowone 数据表中的脏数据
		for item := range util.RowOneMap.Iter() {
			tName := item.Val.(string)
			cleanupByTable(tName, tableInfo.Id)
		}

		// 删除主数据表
		if isNewDataTable {
			if err := tx.Exec("DROP TABLE IF EXISTS " + tableName).Error; err != nil {
				log.Printf("删除主数据表 [%s] 失败，请检查或手动删除：%v\n", tableName, err)
				deleteNoError = false
			}
		} else {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("删除主数据表 [%s] 的脏数据失败，请检查或手动删除：%v\n", tableName, err)
				deleteNoError = false
			}
		}

		// 删除 表data_table_infos 中记录
		if err := db.Table("data_table_infos").Where("id = ?", tableInfo.Id).Delete(nil).Error; err != nil {
			log.Printf("删除 表[data_table_infos] 中 id=%d 的记录失败，请检查或手动删除：%v\n", tableInfo.Id, err)
			deleteNoError = false
		}
		if deleteNoError {
			log.Println("删除脏数据成功")
		}

	}()
	
	var startTime, endTime time.Time
	sT, err := strconv.ParseFloat(res[start][0], 64)
	if err != nil {
		log.Printf("时间解析错误：位置[%d, %d], %v\n", start, 0, err)
		response.Fail(ctx, nil, fmt.Sprintf("时间解析错误，位置[%d, %d]，请检查", start, 0))
		isErrorHappen = true
		return
	}
	lastTime, err := util.ExcelFloatToTime(sT)
	if err != nil {
		log.Printf("时间解析错误：位置[%d, %d], %v\n", start, 0, err)
		response.Fail(ctx, nil, fmt.Sprintf("时间解析错误，位置[%d, %d]，请检查", start, 0))
		isErrorHappen = true
		return
	}
	lastTime = util.FixExcelTimeEdgeError(lastTime)
	startTime = lastTime

	isEndDataLine := false

	// 一行一行的遍历数据，将遍历到的数据存入数据库
	for i := start; i < len(res); i++ {

		row := make(map[string]interface{})
		isUpdate := false

		// 遍历每一列，尝试取出数据
		for j := 0; j < len(res[i]); j++ {

			if j == 0 {
				// 如果是一对一行字段
				db, isRowOne, err := util.CheckAndRecordRowOneData(db, res, i, 0, tableInfo.Id)
				if isRowOne {
					if !isEndDataLine {
						isEndDataLine = true
						eT, err := strconv.ParseFloat(res[i - 1][0], 64)
						if err != nil {
							log.Printf("时间解析错误：位置[%d, %d]，%v\n", i - 1, 0, err)
							response.Fail(ctx, nil, fmt.Sprintf("时间解析错误，位置[%d, %d]，请检查", i - 1, 0))
							isErrorHappen = true
							return
						}
						endT, err := util.ExcelFloatToTime(eT)
						if err != nil {
							log.Printf("时间解析错误：位置[%d, %d]，%v\n", i - 1, 0, err)
							response.Fail(ctx, nil, fmt.Sprintf("时间解析错误，位置[%d, %d]，请检查", i - 1, 0))
							isErrorHappen = true
							return
						}
						endTime = util.FixExcelTimeEdgeError(endT)
					}
					break;
				}
				if err != nil{
					log.Printf("内容解析错误：位置[%d, %d], %v\n", i, j, err)
					response.Fail(ctx, nil, fmt.Sprintf("内容解析错误：位置[%d, %d], 请检查", i, j))
					isErrorHappen = true
					return
				}
				// 如果是一对多行字段
				db, isRowAll, err:= util.CheckAndRecordRowAllData(db, index, res, i, 0, tableInfo.Id)
				if isRowAll {
					if !isEndDataLine {
						isEndDataLine = true
						eT, err := strconv.ParseFloat(res[i - 1][0], 64)
						if err != nil {
							log.Printf("时间解析错误：位置[%d, %d]，%v\n", i - 1, 0, err)
							response.Fail(ctx, nil, fmt.Sprintf("时间解析错误，位置[%d, %d]，请检查", i - 1, 0))
							isErrorHappen = true
							return
						}
						endT, err := util.ExcelFloatToTime(eT)
						if err != nil {
							log.Printf("时间解析错误：位置[%d, %d]，%v\n", i - 1, 0, err)
							response.Fail(ctx, nil, fmt.Sprintf("时间解析错误，位置[%d, %d]，请检查", i - 1, 0))
							isErrorHappen = true
							return
						}
						endTime = util.FixExcelTimeEdgeError(endT)
					}
					break;
				}
				if err != nil {
					log.Printf("内容解析错误：位置[%d, %d]，%v\n", i, j, err)
					response.Fail(ctx, nil, fmt.Sprintf("内容解析错误：位置[%d, %d], 请检查", i, j))
					isErrorHappen = true
					return
				}

				if isEndDataLine {
					break
				}

				// 解析数据时间列
				cT, err := strconv.ParseFloat(res[i][j], 64)
				if err != nil {
					log.Printf("时间解析错误：位置[%d, %d]，%v\n", i, j, err)
					response.Fail(ctx, nil, fmt.Sprintf("时间解析错误，位置[%d, %d]，请检查", i, j))
					isErrorHappen = true
					return
				}
				curTime, err := util.ExcelFloatToTime(cT)
				if err != nil {
					log.Printf("时间解析错误：位置[%d, %d]，%v\n", i, j, err)
					response.Fail(ctx, nil, fmt.Sprintf("时间解析错误，位置[%d, %d]，请检查", i, j))
					isErrorHappen = true
					return
				}
				curTime = util.FixExcelTimeEdgeError(curTime)

				diff := lastTime.Sub(curTime)
				if diff > time.Second {
					log.Printf("时间格式错误：时间为非递增等差排列，出现时间无法对齐错误\n")
					response.Fail(ctx, nil, "上传失败，时间为非递增等差排列")
					isErrorHappen = true
					return
				} else if diff < time.Second {
					row["time"] = curTime.Format(util.ReadableTimeFormat)
					continue
				}

				isTimeErr := false

				diff = curTime.Sub(lastTime)
				if diff < time.Hour-time.Second || diff > time.Hour+time.Second {
					log.Printf("时间格式错误：数据行时间间隔不为一个小时\n")
					log.Printf("错误发生在[%d, %d]\n", i, j)
					log.Println("开始补充空白行...\n")
					isTimeErr = true
					lastTime = lastTime.Add(time.Hour)
					break;
				}

				if !isTimeErr {
					lastTime = curTime
				}

				parsedTime := lastTime.Format(util.ReadableTimeFormat)
				var dummy map[string]interface{}
				err = db.Table(tableName).Where("time = ?", parsedTime).First(&dummy).Error
				if err == nil {
					isUpdate = true
				} else if errors.Is(err, gorm.ErrRecordNotFound) {
					isUpdate = false
				} else {
					log.Println("数据查重出错")
					log.Printf("错误发生在[%d, %d]\n", i, j)
					response.Fail(ctx, nil, fmt.Sprintf("存储数据出错，位置[%d, %d]，请检查", i, j))
					isErrorHappen = true
					return
				}

				row["time"] = parsedTime

				if !isTimeErr {
					i = i - 1
				}

				continue
			}

			// 如果该列没有字段
			if j >= len(index) || index[j] == "" {
				continue
			}

			// 尝试取出数字
			data, ok := util.StringToFloat(res[i][j])
			// 成功取出数字
			if ok {
				row[index[j]] = data
			} else {
				row[index[j]] = res[i][j]
			}

		}

		if isEndDataLine {
			continue
		}

		// 存入数据库
		if isUpdate {
			if err := tx.Table(tableName).Where("time = ?", row["time"]).Updates(row).Error; err != nil {
				log.Printf("第%d行数据更新错误：%v\n", i, err)
				response.Fail(ctx, nil, fmt.Sprintf("第%d行数据更新错误", i))
				isErrorHappen = true
				return
			}
		} else {
			if err := tx.Table(tableName).Create(&row).Error; err != nil {
				log.Printf("第%d行数据存储错误：%v\n", i, err)
				response.Fail(ctx, nil, fmt.Sprintf("第%d行数据存储错误", i))
				isErrorHappen = true
				return
			}
		}
	}

	// 创建文件历史记录
	if err := db.Create(&model.FileHistory{
		UserId:   user.Id,
		FileName: file.Filename,
		FilePath: fullPath,
		Option:   "创建",
	}).Error; err != nil {
		log.Printf("创建文件历史记录失败：%v\n", err)
		response.Fail(ctx, nil, "创建文件历史记录失败")
		isErrorHappen = true
		return
	}

	if !isEndDataLine {
		endTime = startTime
	}

	// 更新数据表开始时间和结束时间
	updateTime := map[string]interface{}{
		"start_time": 	startTime.Format(util.ReadableTimeFormat),
		"end_time": 		endTime.Format(util.ReadableTimeFormat),
	}
	if err := db.Table("data_table_infos").Where("id = ?", tableInfo.Id).Updates(&updateTime).Error; err != nil {
		log.Printf("更新数据表开始时间和结束时间失败：%v\n", err)
		response.Fail(ctx, nil, "存储数据表失败")
		isErrorHappen = true
		return
	}

	// 创建数据历史记录
	if err := db.Create(&model.DataHistory{
		UserId:      user.Id,
		StartTime:   startTime.Format(util.ReadableTimeFormat),
		EndTime:     endTime.Format(util.ReadableTimeFormat),
		StationName: stName,
		System:      system,
		Option:      "创建",
	}).Error; err != nil {
		log.Printf("创建文件历史记录失败：%v\n", err)
		response.Fail(ctx, nil, "创建文件历史记录失败")
		isErrorHappen = true
		return
	}

	if !isErrorHappen {
		if err := tx.Commit().Error; err != nil {
			log.Printf("数据储存错误：%v\n", err)
			response.Fail(ctx, nil, "数据储存错误")
			tx.Rollback()
			isErrorHappen = true
			return
		}
	}

	log.Println("解析上传文件成功")

	response.Success(ctx, gin.H{"FileName": file.Filename}, "上传成功")
}

// @title    List
// @description   提供点集文件列表
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func List(ctx *gin.Context) {

	// 取出请求
	path := ctx.DefaultQuery("path", "/")

	// 获得hour目录下的所有文件
	files, err := util.GetFiles(path)

	if err != nil {
		if path == "/month" {
			response.Fail(ctx, nil, "未上传月度制文件")
		} else if path == "/hour" {
			response.Fail(ctx, nil, "未上传小时制文件")
		} else {
			response.Fail(ctx, nil, "无法处理该文件列表获取请求")
		}
		return
	}

	response.Success(ctx, gin.H{"files": files}, "请求成功")

}

// @title    Download
// @description   下载点集文件
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Download(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在二级以下的用户不能下载文件
	if user.Level < 2 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 取出请求
	path := ctx.DefaultQuery("path", "/")
	file := ctx.DefaultQuery("file", "")

	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))
	ctx.File("./home" + path)
	response.Success(ctx, nil, "请求成功")
}

// @title    DeleteFile
// @description   删除点集文件
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteFile(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在四级以下的用户不能删除文件
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 取出请求
	path := ctx.DefaultQuery("path", "")

	// TODO 移除文件
	if os.Remove(path) != nil {
		response.Fail(ctx, nil, "路径不存在")
		return
	}

	// TODO 创建文件历史记录
	common.GetDB().Create(model.FileHistory{
		UserId:   user.Id,
		FilePath: path,
		Option:   "删除",
	})

	response.Success(ctx, nil, "删除成功")
}

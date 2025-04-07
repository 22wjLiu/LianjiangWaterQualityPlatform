// @Title  fileUtils
// @Description  各种需要使用的文件操作工具函数
package util

import (
	"strings"
	"log"
	"time"
	
	"gorm.io/gorm"
)


// @title    FindStationNameFromFile
// @description   从文件数据中找到站名并返回可录入列字段，站名，数据开始行数和是否找到站名标志
// @param     file [][]string, opt interface{}			文件数据，标记
// @return    []string, string, int, bool			可录入列字段，站名，数据开始行数，是否找到数据段上行标志
func FindStationNameFromFile(file [][]string, opt interface{}) ([]string, string, int, bool, bool){
	// 用于存储字段映射序列
	index := make([]string, len(file[0]))

	// 表示数据的起始行数
	start := 0

	// 用于标记是否遇到标记
	flag := false

	// 用于存储站名
	var stName string
	isFindStationName := false

	// 逐行遍历，尝试寻找站名并取出字段映射
	for i := 0; i < len(file); i++ {
		for j := 0; j < len(file[i]); j++ {
			
			cell := file[i][j]
			runes := []rune(cell)

			if !isFindStationName {
				for _, key := range StationNameFlagMap.Keys() {
					value, _ := StationNameFlagMap.Get(key)
					prefixRunes := []rune(value.(string))
					vLen := len(prefixRunes)
					
					// 成功找到站名
					if strings.HasPrefix(cell, value.(string)) {
						if len(runes) > vLen && runes[vLen] == ':' {
							stName = string(runes[vLen+1:])
						} else if len(runes) > vLen && runes[vLen] == '：' {
							stName = string(runes[vLen+1:])
						} else {
							stName = string(runes[vLen:])
						}
						isFindStationName = true
						break
					}
				}
				if isFindStationName {
					continue
				}
			}

			p := ""
			// 寻找最长前缀匹配
			for k := len(cell); k > 0; k-- {
				str, ok := PointMap.Get(cell[:k])
				if ok {
					p = str.(string)
					break
				}
			}
			
			// 成功匹配映射字段，则记录该字段
			if p != "" {
				index[j] = p
			}

			// 遇到标记
			if cell == opt.(string) {
				flag = true
			}
		}

		// 如果遇到数据段上一行标记，记录数据初始位置，并退出字段搜寻
		if flag {
			start = i + 1
			break
		}
	}

	return index, stName, start, flag, isFindStationName
}

// @title    CheckAndRecordRowOneData
// @description   检查并录入一对一行字段数据
// @param     db *gorm.DB, data [][]string, i int, j int, tableId int			数据库实例，数据，行，列，数据表ID
// @return    *gorm.DB, bool																							数据库实例，是否存在一对一行字段数据
func CheckAndRecordRowOneData(db *gorm.DB, data [][]string, i int, j int, tableId uint) (*gorm.DB, bool, error){
	row, ok := RowOneMap.Get(data[i][j])
	// 如果是一对一行字段
	if !ok {
		return db, false, nil
	}

	log.Printf("检测一对对行字段 [%s]\n", data[i][j])

	rowOneTime := []string{"created_at"}
	rowOneStr  := []string{"value"}

	rowOneTableName := row.(string)

	// 查看是否存在表
	if !db.Migrator().HasTable(rowOneTableName) {
		log.Printf("未检测到一对多行字段表 [%s]\n", rowOneTableName)
		log.Printf("开始创建一对多行字段表... [%s]\n", rowOneTableName)

		sql := BuildCreateTableSQL_Str_T_FId(rowOneTableName, rowOneTime, rowOneStr, "table_id", "data_table_infos")
		
		// 执行 sql 创建表
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("创建表[%s]失败: %v\n", rowOneTableName, err)
			return db, true, err
		} else {
			log.Printf("创建表[%s]成功\n", rowOneTableName)
		}
	}

	rowOne := map[string]interface{}{
		"value":		data[i][j + 1],
		"table_id": tableId,
	}

	// 存入数据库
	err := db.Table(rowOneTableName).Create(rowOne).Error
	if err != nil {
			log.Printf("插入失败: %v", err)
			return db, true, err
	} else {
			log.Println("插入成功！")
	}

	return db, true, nil
}

// @title    CheckAndRecordRowAllData
// @description   检查并录入一对多行字段数据
// @param     db *gorm.DB, data [][]string, i int, j int, tableId int			数据库实例，数据，行，列，数据表ID
// @return    *gorm.DB, bool																	数据库实例，是否存在一对多行字段数据
func CheckAndRecordRowAllData(db *gorm.DB, index []string, data [][]string, i int, j int, tableId uint) (*gorm.DB, bool, error){
	row, ok := RowAllMap.Get(data[i][j])
	// 如果是一对多行字段
	if !ok {
		return db, false, nil
	}

	log.Printf("检测一对多行字段 [%s]\n", data[i][j])

	rowAllTime := []string{"created_at"}
	var rowAllStr []string
	
	for _, key := range PointMap.Keys() {
		v,_ := PointMap.Get(key)
		rowAllStr = append(rowAllStr, v.(string))
	}

	rowAllTableName := row.(string)

	// 查看是否存在表
	if !db.Migrator().HasTable(rowAllTableName) {
		log.Printf("未检测到一对多行字段表 [%s]\n", rowAllTableName)
		log.Printf("开始创建一对多行字段表... [%s]\n", rowAllTableName)

		sql := BuildCreateTableSQL_Str_T_FId(rowAllTableName, rowAllTime, rowAllStr, "table_id", "data_table_infos")
		
		// 执行 sql 创建表
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("创建表[%s]失败: %v\n", rowAllTableName, err)
			return db, true, err
		} else {
			log.Printf("创建表[%s]成功\n", rowAllTableName)
		}
	}

	rowAll := map[string]interface{}{
		"table_id": tableId,
	}

	for k := j + 1; k < len(data[i]); k++ {
		if index[k] == "" {
			continue
		}
		rowAll[index[k]] = data[i][k]
	}

	// 存入数据库
	err := db.Table(rowAllTableName).Create(rowAll).Error
	if err != nil {
			log.Printf("插入失败: %v", err)
			return db, true, err
	} else {
			log.Println("插入成功！")
	}

	return db, true, nil
}

// @title    FixExcelTimeEdgeError
// @description   修正浮点运算造成的Excel时间解析误差
// @param     t time.Time			待修复时间
// @return    time.Time				修复后时间
func FixExcelTimeEdgeError(t time.Time) time.Time {
	min := t.Minute()
	sec := t.Second()

	// 59:59.999 → 60:00:00
	if min == 59 && sec >= 58 {
		return time.Date(
			t.Year(), t.Month(), t.Day(),
			t.Hour()+1, 0, 0, 0, t.Location(),
		)
	}

	// 00:00:01 → 00:00:00（轻微误差）
	if min == 0 && sec <= 1 {
		return time.Date(
			t.Year(), t.Month(), t.Day(),
			t.Hour(), 0, 0, 0, t.Location(),
		)
	}

	// 其他情况保留原值
	return t
}
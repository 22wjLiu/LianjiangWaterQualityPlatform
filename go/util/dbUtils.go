// @Title  fileUtils
// @Description  各种需要使用的数据库操作工具函数
package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"lianjiang/common"
	"lianjiang/model"
	
	"gorm.io/gorm"

	"github.com/spf13/viper"
)

// 系统用户ID
var sysUserId uint = 911

// 自定义错误变量
var ErrMissingBackupFile = errors.New("缺少sql备份文件")

// @title    InitSysUser
// @description   初始化系统用户
// @param     无入参
// @return    error	是否出错
func InitSysUser() error {
	// 获取数据库指针
	db := common.GetDB()

	var sysUser model.User

	err := db.Model(&model.User{}).Unscoped().Where("id = ?", sysUserId).First(&sysUser).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查找系统用户失败")
	} else if err == nil {
		return nil
	}

	sysUser =  model.User{
		Id: sysUserId,
		Name: "systemYA&&&¥¥¥¥¥",
		Email: "-",
		Password: "-",
		Level:	-1,
	}

	// 开启事务
	tx := db.Begin()

	if err := tx.Create(&sysUser).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建系统用户失败")
	}

	if err := tx.Delete(&sysUser).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("隐藏系统用户失败")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("初始化系统用户失败")
	}

	return nil
}

// @title    HasBackUpSql
// @description   判断是否存在sql备份文件
// @param     path string, sqlFiles []string			路径, sql文件
// @return    err				是否发生错误
func HasBackUpSql(path string, sqlFiles []string) error {
	for _, fileName := range sqlFiles {
			fullPath := filepath.Join(path, fileName)
			if _, err := os.Stat(fullPath); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("%w: %s", ErrMissingBackupFile, fileName)
				}
				return err // 其他错误，比如权限问题
			}
	}

	return nil
}

// @title    RecoverFromBackupSql
// @description   通过sql备份文件恢复数据
// @param     path string, sqlFiles []string, isChecked bool			路径, sql文件, 是否检查文件存在
// @return    err				是否发生错误
func RecoverFromBackupSql(path string, sqlFiles []string, isChecked bool) error {
	// 检查备份文件
	if !isChecked {
		if err := HasBackUpSql(path, sqlFiles); err != nil {
			return err
		}
	}

	for _, file := range sqlFiles {
			fullPath := filepath.Join(path, file)

			dbName := viper.GetString("datasource.database")
			tableName := strings.TrimSuffix(file, ".sql")

			// 构造执行命令（以 MySQL 为例）
			cmd := exec.Command("mysql",
			 	"-u",
			  viper.GetString("datasource.username"),
			  "-p" + viper.GetString("datasource.password"),
				dbName)
			
			sqlFile, err := os.Open(fullPath)
			if err != nil {
					return fmt.Errorf("打开SQL文件失败: %w", err)
			}
			defer sqlFile.Close()
			cmd.Stdin = sqlFile

			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("执行恢复失败 [%s → %s.%s]: %v\n输出: %s",
        file, dbName, tableName, err, output)
			}
	}

	return nil
}


// @title    BackUpMapSql
// @description   备份sql文件
// @param     path string, sqlFiles []string			路径, sql文件
// @return    error			是否发生错误
func BackUpSql(path string, sqlFiles []string) error {
	// 确保目录存在
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("创建备份目录失败: %w", err)
	}

	for _, file := range sqlFiles {
		fullPath := filepath.Join(path, file)
		dbName := viper.GetString("datasource.database")
		tableName := strings.TrimSuffix(file, ".sql")

		// 构造 mysqldump 命令
		cmd := exec.Command("mysqldump",
			 "-u",
			 viper.GetString("datasource.username"),
			 "-p" + viper.GetString("datasource.password"),
			 dbName,
			 tableName)

		// 创建输出文件
		outFile, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("创建备份文件失败 [%s.%s -> %s]: %w", dbName, tableName, fullPath, err)
		}

		cmd.Stdout = outFile

		// 执行命令
		err = cmd.Run()
		outFile.Close()
		if err != nil {
				return fmt.Errorf("执行 mysqldump 失败 [%s.%s -> %s]: %w", dbName, tableName, fullPath, err)
		}

		log.Printf("备份表 [%s.%s] → 文件 [./%s] 成功\n", dbName, tableName, fullPath)
	}

	return nil
}

// @title    DbConditionsEqual
// @description   数据库操作设置等于查询条件
// @param    *gorm.DB      数据库实例指针
// @return   void					 没有出参
func DbConditionsEqual(db *gorm.DB, conditions map[string]interface{}) *gorm.DB{
	for field, value := range conditions {
		if str, ok := value.(string); ok && str == "" {
			continue
		}
		db = db.Where("`" + field + "` = ?", value)
	}
	return db
}

// @title    DbConditionsLike
// @description   数据库操作设置Like查询条件
// @param    *gorm.DB      数据库实例指针
// @return   void					 没有出参
func DbConditionsLike(db *gorm.DB, conditions map[string]interface{}) *gorm.DB{
	for field, value := range conditions {
		if str, ok := value.(string); ok && str == "" {
			continue
		}
		db = db.Where("`" + field + "` like ?", "%" + value.(string) + "%")
	}
	return db
}


// @title    BuildCreateTableSQL_Str_T_FId
// @description   动态创建建表sql，包含主键Id、时间字段(TIMESTAMP)、字符串字段（VACHAR(30)）和一个外键
// @param    tableName string, timeCols []string, strCols []string, fIdName string, fTableName string      表名，字符串字段，时间字段，外键名和外键表
// @return   string					 建表sql
func BuildCreateTableSQL_Str_T_FId(tableName string, timeCols []string, strCols []string, fIdName string, fTableName string) string {
	var colDefs []string

	// 主键
	colDefs = append(colDefs, "`id` BIGINT(20) UNSIGNED PRIMARY KEY AUTO_INCREMENT")

	// 时间字段
	for _, col := range timeCols {
		colDefs = append(colDefs, fmt.Sprintf("`%s` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL", col))
	}


	// 字符串字段
	for _, col := range strCols {
		colDefs = append(colDefs, fmt.Sprintf("`%s` VARCHAR(30) DEFAULT '' NOT NULL", col))
	}

	// 外键字段
	if fIdName != "" && fTableName != "" {
		colDefs = append(colDefs, fmt.Sprintf("`%s` BIGINT(20) UNSIGNED NOT NULL", fIdName))
		colDefs = append(colDefs, fmt.Sprintf("FOREIGN KEY (`%s`) REFERENCES `%s`(`id`) ON DELETE CASCADE ON UPDATE CASCADE", fIdName, fTableName))
	}

	// 生成sql
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n  %s\n);",
			tableName, strings.Join(colDefs, ",\n  "))
	return sql
}

// @title    BuildCreateTableSQL_Str_T
// @description   动态创建建表sql，包含主键Id、时间字段(TIMESTAMP)、字符串字段（VACHAR(30))
// @param    tableName string, timeCols []string, strCols []string, fIdName string, fTableName string      表名，字符串字段，时间字段，外键名和外键表
// @return   string					 建表sql
func BuildCreateTableSQL_Str_T(tableName string, timeCols []string, strCols []string) string {
	var colDefs []string
	var uniqueKeyDefs []string

	// 主键
	colDefs = append(colDefs, "`id` BIGINT(20) UNSIGNED PRIMARY KEY AUTO_INCREMENT")

	// 时间字段
	for _, col := range timeCols {
		colDefs = append(colDefs, fmt.Sprintf("`%s` TIMESTAMP NOT NULL", col))
		uniqueKeyDefs = append(uniqueKeyDefs, fmt.Sprintf("`%s`", col))
	}

	// 字符串字段
	for _, col := range strCols {
		colDefs = append(colDefs, fmt.Sprintf("`%s` VARCHAR(30) DEFAULT '' NOT NULL", col))
	}

	// 生成sql
	uniqueKeySql := fmt.Sprintf("UNIQUE (%s)", strings.Join(uniqueKeyDefs, ", "))

	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n  %s,\n  %s\n);",
			tableName, strings.Join(colDefs, ",\n  "), uniqueKeySql)
	return sql
}
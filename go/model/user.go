// @Title  user
// @Description  定义跟帖
package model

import "gorm.io/gorm"

// user			定义用户
type User struct {
	gorm.Model        // gorm的模板
	Name       string `gorm:"type:varchar(20);not null;unique"` // 用户名称
	Email      string `gorm:"type:varchar(50);not null;unique"` // 邮箱
	Password   string `gorm:"size:255;not null"`                // 密码
	Level      int    `gorm:"type:int;not null"`                // 用户权限等级
}

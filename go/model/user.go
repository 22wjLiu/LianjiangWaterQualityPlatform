// @Title  user
// @Description  定义用户表
package model

import (
	"gorm.io/gorm"
)

// user			定义用户
type User struct {
	Id		   uint   `json:"id" gorm:"type:uint;not null;unique"`																			// 用户ID
	CreatedAt  Time   `json:"created_at" gorm:"type:timestamp;not null"`  														// 创建时间
	UpdatedAt  Time   `json:"updated_at" gorm:"type:timestamp;not null;default:current_timestamp"`  	// 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`																				// 删除时间
	Name       string `json:"name" gorm:"type:varchar(20);not null;unique"` 													// 用户名称
	Email      string `json:"email" gorm:"type:varchar(50);not null;unique"`													// 邮箱
	Password   string `json:"password" gorm:"size:255;not null"`                											// 密码
	Level      int    `json:"level" gorm:"type:int;not null"`                													// 用户权限等级
}

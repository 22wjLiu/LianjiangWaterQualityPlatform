// @Title  user_dto
// @Description  用于封装用户信息
package dto

import (
	"lianjiang/model"
)

// UserDto			定义了用户的基本信息
type UserDto struct {
	Id	  		uint   	`json:"id"`			// 用户ID
	CreatedAt 	model.Time	`json:"created_at"`	// 创建时间
	UpdatedAt 	model.Time	`json:"updated_at"`	// 更新时间
	Name  		string 	`json:"name"`  		// 用户名
	Email 		string 	`json:"email"` 		// 邮箱
	Level 		int    	`json:"level"` 		// 权限等级
}

// BatchDeleteId	定义了要删除的一组用户ID
type BatchDeleteId struct {
	Ids []uint `json:"ids"`
}

// @title    ToUserDto
// @description   用户信息封装
// @param    user model.User       接收一个用户类
// @return   UserDto			   返回一个用户的基本信息类
func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:  user.Name,
		Email: user.Email,
		Level: user.Level,
	}
}

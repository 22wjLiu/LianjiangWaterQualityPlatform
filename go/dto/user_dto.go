// @Title  user_dto
// @Description  用于封装用户信息
package dto

import (
	"lianjiang/model"
)

// UserDto			定义了用户的基本信息
type UserDto struct {
	Name  string `json:"name"`  // 用户名称
	Email string `json:"email"` // 邮箱
	Level int    `json:"level"` // 权限等级
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

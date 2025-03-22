// @Title  user
// @Description  定义跟帖
package vo

import "gorm.io/gorm"

// user			定义用户
type UserRequest struct {
	gorm.Model        // gorm的模板
	Name       string // 用户名称
	Email      string // 邮箱
	Password   string // 密码
	Verify     string // 验证码
}

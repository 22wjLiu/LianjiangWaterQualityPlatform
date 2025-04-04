// @Title  mapHistory
// @Description  定义映射操作历史记录
package model

// MapHistory			定义映射操作历史记录
type MapHistory struct {
	UserId    uint   `json:"user_id" gorm:"type:uint;not null"`        // 用户Id
	CreatedAt Time   `json:"created_at" gorm:"type:timestamp;not null"` // 操作时间
	Id        string `json:"id" gorm:"type:varchar(50);not null"`      // 映射表
	Key       string `json:"key" gorm:"type:varchar(50);not null"`     // 主键
	Value     string `json:"value" gorm:"type:varchar(50);not null"`   // 对应值
	Option    string `json:"option" gorm:"type:varchar(20);not null;"` // 操作方法
}

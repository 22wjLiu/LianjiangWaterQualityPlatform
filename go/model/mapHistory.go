// @Title  mapHistory
// @Description  定义映射操作历史记录
package model

// MapHistory			定义映射操作历史记录
type MapHistory struct {
	Id        uint 	 	`json:"id" gorm:"type:uint;not null"`      							// ID
	UserId    uint   	`json:"user_id" gorm:"type:uint;not null"`        			// 用户ID
	VerName		string	`json:"ver_name" gorm:"type:varchar(50);"`		 					// 版本名
	CreatedAt Time    `json:"created_at" gorm:"type:datetime;autoCreateTime"` // 创建时间
	Table			string	`json:"table" gorm:"type:varchar(50);"`     						// 映射类型
	Key       string 	`json:"key" gorm:"type:varchar(50);"`     							// 主键
	Value     string 	`json:"value" gorm:"type:varchar(50);"`   							// 对应值
	Option    string 	`json:"option" gorm:"type:varchar(20);not null;"` 			// 操作方法
}

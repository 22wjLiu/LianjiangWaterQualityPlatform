// @Title  mapHistory
// @Description  定义映射操作历史记录
package model

// MapHistory			定义映射操作历史记录
type MapHistory struct {
	Id        uint 	 	`json:"id" gorm:"type:uint;not null"`      					// ID
	UserId    uint   	`json:"user_id" gorm:"type:uint;not null"`        	// 用户ID
	VerId			uint	 	`json:"map_id" gorm:"type:uint;not null"`		 				// 版本ID
	CreatedAt Time   	`json:"created_at" gorm:"type:timestamp;not null"` 	// 操作时间
	Table			string	`json:"table" gorm:"type:varchar(50);not null"`     // 映射类型
	Key       string 	`json:"key" gorm:"type:varchar(50);not null"`     	// 主键
	Value     string 	`json:"value" gorm:"type:varchar(50);not null"`   	// 对应值
	Option    string 	`json:"option" gorm:"type:varchar(20);not null;"` 	// 操作方法
	User 			*User  	`json:"-" gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` 		// 外键 UserId -> User.Id
	MapVersion *MapVersion `json:"-" gorm:"foreignKey:VerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 外键 MapId -> MapBackUp.Id
}

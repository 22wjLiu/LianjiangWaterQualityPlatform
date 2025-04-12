// @Title  dataHistory
// @Description  定义数据操作历史记录
package model

// DataHistory			定义数据操作历史记录
type DataHistory struct {
	Id        	uint 	 `json:"id" gorm:"type:uint;not null"`      						// ID
	UserId      uint   `json:"user_id" gorm:"type:uint;not null"`             // 用户Id
	CreatedAt 	Time   `json:"created_at" gorm:"type:datetime;autoCreateTime"`// 创建时间
	Time   			Time 	 `json:"time" gorm:"type:datetime;not null"`   					// 时间
	StationName string `json:"station_name" gorm:"type:varchar(50);not null"` // 站名
	System      string `json:"system" gorm:"type:varchar(10);not null"`       // 制度
	VersionName string `json:"version_name" gorm:"type:varchar(50);not null"` // 映射版本名
	Option      string `json:"option" gorm:"type:varchar(20);not null;"`      // 操作方法
}

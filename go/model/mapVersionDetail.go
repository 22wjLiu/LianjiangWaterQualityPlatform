// @Title  mapBackup
// @Description  定义映射版本详情

package model

// MapVersionDetail			定义映射版本详情
type MapVersionDetail struct {
	Id 		  		uint   		`json:"id" gorm:"primaryKey;type:uint;not null"` 			  						// ID
	Table     	string 		`json:"table" gorm:"type:varchar(50);not null"`   									// 映射表
	Key       	string 		`json:"key" gorm:"type:varchar(50);not null"`     									// 主键
	Value     	string 		`json:"value" gorm:"type:varchar(50);not null"`   									// 对应值
	VerId				uint			`json:"map_id" gorm:"type:uint;not null"`														// 版本ID
	MapVersion  *MapVersion `json:"-" gorm:"foreignKey:VerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 外键 VerId -> MapVersion.Id
}
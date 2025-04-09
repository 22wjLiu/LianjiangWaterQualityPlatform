// @Title  mapBackup
// @Description  定义映射版本

package model

// MapVersion			定义映射版本表
type MapVersion struct {
	Id 		  			uint   	`json:"id" gorm:"primaryKey;type:uint;not null"` 																	// ID
	VersionName   string 	`json:"version_name" gorm:"type:varchar(255);unique;default:'latest';not null"` 	// 备份版本名
	CreatedAt   	Time 		`json:"created_at" gorm:"type:timestamp;not null;autoCreateTime"`									// 创建时间
	UpdatedAt 		Time    `json:"updated_at" gorm:"type:timestamp;not null;default:current_timestamp"`  		// 更新时间
	Active      	bool    `json:"active" gorm:"type:boolean;default:false;not null"`  											// 是否使用中
}

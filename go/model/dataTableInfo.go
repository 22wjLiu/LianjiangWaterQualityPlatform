// @Title  DataTableInfo
// @Description  定义数据表详情
package model

// DataTableInfo			定义数据表详情
type DataTableInfo struct {
	Id        			uint 	 	`json:"id" gorm:"type:uint;not null"`      									// ID
	MapVerId				uint		`json:"map_ver_id" gorm:"type:uint;not null"`								// 使用的映射版本ID
	DataTableName   string 	`json:"table_name" gorm:"type:varchar(50);not null"`   			// 数据表名
	StationName   	string 	`json:"table_name" gorm:"type:varchar(50);not null"`   			// 站名
	System      		string `json:"system" gorm:"type:varchar(50);not null"`      			  // 制度
	File    	 			string 	`json:"file" gorm:"type:varchar(50);not null"`     					// 数据文件
	Active      		bool    `json:"active" gorm:"type:boolean;default:false;not null"`  // 是否使用中
	StartTime       Time 		`json:"start_time" gorm:"type:datetime;default:null"`  			// 开始日期
	EndTime         Time 		`json:"end_time" gorm:"type:datetime;default:null"`    			// 终止日期
	MapVersion			*MapVersion `json:"-" gorm:"foreignKey:MapVerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 外键 MapVerId -> MapVersion.Id
}
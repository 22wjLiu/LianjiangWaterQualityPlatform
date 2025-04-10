// @Title  DataTableInfo
// @Description  定义数据表详情
package model

// DataTableInfo			定义数据表详情
type DataTableInfo struct {
	Id        			uint 	 	`json:"id" gorm:"type:uint;not null"`      									// ID
	MapVerId				uint		`json:"map_ver_id" gorm:"type:uint;not null"`								// 映射版本ID
	FileId    	 		uint 	`json:"file" gorm:"type:varchar(50);not null"`     						// 文件ID
	DataTableName   string 	`json:"data_table_name" gorm:"type:varchar(50);not null"`   // 数据表名
	StationName   	string 	`json:"station_name" gorm:"type:varchar(50);not null"`   		// 站名
	System      		string `json:"system" gorm:"type:varchar(10);not null"`      			  // 制度
	Active      		bool    `json:"active" gorm:"type:boolean;default:false;not null"`  // 是否使用中
	CreatedAt 			Time   `json:"created_at" gorm:"type:datetime;autoCreateTime"`  		// 创建时间
	StartTime       Time 		`json:"start_time" gorm:"type:datetime;default:null"`  			// 开始日期
	EndTime         Time 		`json:"end_time" gorm:"type:datetime;default:null"`    			// 终止日期
	MapVersion			*MapVersion `json:"-" gorm:"foreignKey:MapVerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 外键 MapVerId -> MapVersion.Id
	FileInfo				*FileInfo 	`json:"-" gorm:"foreignKey:FileId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` 	// 外键 FileId -> FileInfo.Id
}
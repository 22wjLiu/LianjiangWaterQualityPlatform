// @Title  fileHistory
// @Description  定义文件操作历史记录
package model

// FileHistory			定义文件操作历史记录
type FileHistory struct {
	Id        uint 	 `json:"id" gorm:"type:uint;not null"`      				 	// ID
	UserId    uint   `json:"user_id" gorm:"type:uint;not null"`          	// 用户Id
	CreatedAt Time   `json:"created_at" gorm:"type:timestamp;not null"`  	// 操作时间
	System    string `json:"system" gorm:"type:varchar(10);not null"`     // 制度
	FileName  string `json:"file_name" gorm:"type:varchar(50);not null"` 	// 文件名
	FileType  string `json:"file_type" gorm:"type:varchar(20);not null"`	// 文件类型
	FilePath  string `json:"file_path" gorm:"type:varchar(100);not null"` // 文件路径
	Option    string `json:"option" gorm:"type:varchar(20);not null;"`   // 操作方法
}

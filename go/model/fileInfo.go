// @Title  fileInfo
// @Description  定义文件表详情
package model

// FileInfo			定义文件表详情
type FileInfo struct {
	Id        uint 	 `json:"id" gorm:"type:uint;not null"`      				 // ID
	CreatedAt Time   `json:"created_at" gorm:"type:timestamp;not null"`  // 创建时间
	UpdatedAt Time   `json:"updated_at" gorm:"type:timestamp;not null;default:current_timestamp"`  // 更新时间
	FileName  string `json:"file_name" gorm:"type:varchar(50);not null"` // 文件名
	FileType  string `json:"file_type" gorm:"type:varchar(50);not null"` // 文件类型
	FilePath  string `json:"file_path" gorm:"type:varchar(100);not null"` // 文件路径
}
// @Title  RowOne
// @Description  定义点集描述信息
package model

import (
	"time"
)

// RowOne			定义点集描述信息
type RowOne struct {
	StartTime   time.Time `json:"start_time" gorm:"type:timestamp;primary_key;autoIncrement:false"`  // 记录的开始日期
	EndTime     time.Time `json:"end_time" gorm:"type:timestamp;primary_key;autoIncrement:false"`    // 记录的终止日期
	StationName string    `json:"station_name" gorm:"type:char(50);primary_key;autoIncrement:false"` // 站名
	Detail      string    `json:"detail" gorm:"type:char(50);"`                                      // 站名
	CreatedAt   Time      `json:"created_at" gorm:"type:timestamp"`                                  // 记录的存入日期
}

// @Title  Data
// @Description  用于定义行字段传回点
package vo

import "time"

// Data			定义传回的点
type Data struct {
	StartTime time.Time `json:"start_time"` // 初始时间
	EndTime   time.Time `json:"end_time"`   // 初始时间
	Field     string    `json:"field"`      // 字段值
}

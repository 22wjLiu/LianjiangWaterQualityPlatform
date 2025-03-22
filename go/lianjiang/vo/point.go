// @Title  Point
// @Description  用于定义传回点
package vo

import "time"

// Point			定义传回的点
type Point struct {
	Time  time.Time `json:"time"`  // 时间
	Field float64   `json:"field"` // 字段值
}

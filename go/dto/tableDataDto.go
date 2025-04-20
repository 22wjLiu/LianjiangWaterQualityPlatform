// @Title  tableDataDto.go
// @Description  用于封装表格数据
package dto

type StationNameData struct {
	System string `json:"system" gorm:"column:system"`
	Name string `json:"name" gorm:"column:station_name"`
}
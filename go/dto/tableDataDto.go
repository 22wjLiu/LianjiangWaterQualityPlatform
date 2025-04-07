// @Title  tableDataDto.go
// @Description  用于封装表格数据
package dto

type StationNameData struct {
	Value string `json:"value" gorm:"column:station_name"`
}
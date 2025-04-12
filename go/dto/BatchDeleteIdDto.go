// @Title  batchDeleteIdDto
// @Description  用于封装要删除的一组ID
package dto

import (
	"lianjiang/model"
)

// BatchDeleteId	定义了要删除的一组ID
type BatchDeleteId struct {
	Ids []uint `json:"ids"`
}

// BatchDeleteTimeId	定义了要删除的一组时间ID
type BatchDeleteTimeId struct {
	VersionName		string 	`json:versionName`
	StationName		string 	`json:stationName`
	System				string 	`json:system`
	DataTableName string 	`json:dataTableName`
	Times []model.Time `json:"times"`
}
// @Title  mapHistoryDto
// @Description  用于封装映射历史信息
package dto

import (
	"lianjiang/model"
)

type MapHistoryWithVerName struct {
	Id          string 			`json:"id"`
	UserId      uint   			`json:"user_id"`
	CreatedAt   model.Time  `json:"created_at"`
	Table				string			`json:"table"`
	Key         string 			`json:"key"`
	Value       string 			`json:"value"`
	Option      string 			`json:"option"`
	VersionName string 			`json:"version_name"`
}



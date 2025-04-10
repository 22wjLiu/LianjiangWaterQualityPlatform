// @Title  mapVersionInfos
// @Description  用于封装用来映射信息
package dto

// MapVersionInfos	定义了映射信息
type MapVersionInfos struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

// MapTables	定义了映射类型
type MapTables struct {
	Value string `json:"value"`
}

// CreateMapDetail	定义了创建映射详情类型
type CreateMapDetail struct {
	Table string `json:table`
	Key		string `json:value`
	Value string `json:"value`
	Formula string `json:"formula"`
}
// @Title  batchDeleteIdDto
// @Description  用于封装要删除的一组ID
package dto

// BatchDeleteId	定义了要删除的一组ID
type BatchDeleteId struct {
	Ids []uint `json:"ids"`
}
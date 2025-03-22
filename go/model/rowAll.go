// @Title  RowAll
// @Description  定义点集描述信息
package model

import (
	"time"
)

// RowAll			定义点集描述信息
type RowAll struct {
	StartTime                       time.Time `json:"start_time" gorm:"type:timestamp;primary_key;autoIncrement:false"`  // 记录的开始日期
	EndTime                         time.Time `json:"end_time" gorm:"type:timestamp;primary_key;autoIncrement:false"`    // 记录的终止日期
	StationName                     string    `json:"station_name" gorm:"type:char(50);primary_key;autoIncrement:false"` // 站点名称
	Cod                             string    `json:"cod" gorm:"type:char(50);"`                                         // 化学需氧量
	FiveDaysNiochemicalOxygenDemand string    `json:"five_days_niochemical_oxygen_demand" gorm:"type:char(50);"`         // 五日生化需氧量
	Se                              string    `json:"se" gorm:"type:char(50);"`                                          // 硒
	As                              string    `json:"as" gorm:"type:char(50);"`                                          // 砷
	Hg                              string    `json:"hg" gorm:"type:char(50);"`                                          // 汞
	Fluoride                        string    `json:"fluoride" gorm:"type:char(50);"`                                    // 氟化物
	Petroleum                       string    `json:"petroleum" gorm:"type:char(50);"`                                   // 石油类
	FecalColiform                   string    `json:"fecal_coliform" gorm:"type:char(50);"`                              // 粪大肠菌群
	Temperature                     string    `json:"temperature" gorm:"type:char(50);"`                                 // 水温
	PH                              string    `json:"ph" gorm:"type:char(50);"`                                          // pH值
	DO                              string    `json:"do" gorm:"type:char(50);"`                                          // 溶解氧
	EC                              string    `json:"ec" gorm:"type:char(50);"`                                          // 电导率
	Turbidity                       string    `json:"turbidity" gorm:"type:char(50);"`                                   // 浊度
	CODMII                          string    `json:"codmii" gorm:"type:char(50);"`                                      // 高锰酸盐指数
	NH_N                            string    `json:"nh_n" gorm:"type:char(50);column:nh_n"`                             // 氨氮
	TP                              string    `json:"tp" gorm:"type:char(50);"`                                          // 总磷
	TN                              string    `json:"tn" gorm:"type:char(50);"`                                          // 总氮
	CODcr                           string    `json:"codcr" gorm:"type:char(50);column:codcr"`                           // CODcr
	CN                              string    `json:"cn" gorm:"type:char(50);"`                                          // 氰化物
	VolatilePenol                   string    `json:"volatile_penol" gorm:"type:char(50);"`                              // 挥发酚
	Cr                              string    `json:"cr" gorm:"type:char(50);"`                                          // 六价铬
	Cu                              string    `json:"cu" gorm:"type:char(50);"`                                          // 铜
	Zn                              string    `json:"zn" gorm:"type:char(50);"`                                          // 锌
	Pb                              string    `json:"pb" gorm:"type:char(50);"`                                          // 铅
	Cd                              string    `json:"cd" gorm:"type:char(50);"`                                          // 镉
	LAS                             string    `json:"las" gorm:"type:char(50);"`                                         // 阴离子表面活性剂
	SOx                             string    `json:"sox" gorm:"type:char(50);column:sox"`                               // 硫化物
	CumulativeDischarge             string    `json:"cumulative_discharge" gorm:"type:char(50);"`                        // 累计流量
	WaterDischarge                  string    `json:"water_discharge" gorm:"type:char(50);"`                             // 水流量
	WaterLevel                      string    `json:"water_level" gorm:"type:char(50);"`                                 // 水位
	PeriodCumulativeFlow            string    `json:"period_cumulative_flow" gorm:"type:char(50);"`                      // 时段累积流量
	SectionalMeanVelocity           string    `json:"sectional_mean_velocity" gorm:"type:char(50);"`                     // 断面平均流速
	SectionalArea                   string    `json:"sectional_area" gorm:"type:char(50);"`                              // 断面面积
	TotalCumulativeFlow             string    `json:"total_cumulativeflow" gorm:"type:char(50);"`                        // 总累积流量
	CurrentInstantaneousVelocity    string    `json:"current_instantaneous_velocity" gorm:"type:char(50);"`              // 当前瞬时流速
	InstantaneousDelivery           string    `json:"instantaneous_delivery" gorm:"type:char(50);"`                      // 瞬时流量
	CreatedAt                       Time      `json:"created_at" gorm:"type:timestamp"`                                  // 记录的存入日期
}

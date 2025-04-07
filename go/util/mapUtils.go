// @Title  fileUtils
// @Description  各种需要使用的映射操作工具函数
package util

import (
		"errors"
		"fmt"
		"log"
		"lianjiang/common"
		"lianjiang/model"

		"gorm.io/gorm"

		Map "github.com/orcaman/concurrent-map"
)

// 映射表sql文件存储路径
var MapSqlPath = "./home/backup"

// 映射表sql文件名
var MapSqlFiles = []string{"map_versions.sql", "map_version_details.sql"}

// 点集字段映射表
var PointMap Map.ConcurrentMap = Map.New()

// 行唯一字段映射表
var RowOneMap Map.ConcurrentMap = Map.New()

// 行多字段映射表
var RowAllMap Map.ConcurrentMap = Map.New()

// 制度映射表
var SysMap Map.ConcurrentMap = Map.New()

// 站名注册表
var StationMap Map.ConcurrentMap = Map.New()

// 数据注册表
var DataMap Map.ConcurrentMap = Map.New()

// 数据段上一行标记映射表
var OptMap Map.ConcurrentMap = Map.New()

// 站名标记映射表
var StationNameFlagMap Map.ConcurrentMap = Map.New()

// 映射表
var MapMap = map[string]*Map.ConcurrentMap{
	"列字段映射":    &PointMap,
	"行字段一对多映射": &RowAllMap,
	"行字段一对一映射": &RowOneMap,
	"时间制映射":    &SysMap,
	"站名映射":     &StationMap,
	"数据符号映射":   &DataMap,
	"数据段上一行标记映射":     &OptMap,
	"站名标记映射":		&StationNameFlagMap,
}

// @title    InitDefaultMapMap
// @description   设置映射表为默认值
// @param     void			没有入参
// @return    void			没有回参
func InitDefaultMapMap() {
	// 初始化点集字段映射表
	PointMap.Set("监测断面", "monitoring_section")
	PointMap.Set("监测指标", "monitoring_indicators")
	PointMap.Set("监测时间", "monitoring_time")
	PointMap.Set("时间", "time")
	PointMap.Set("水温", "temperature")
	PointMap.Set("pH", "pH")
	PointMap.Set("化学需氧量", "cod")
	PointMap.Set("五日生化需氧量", "5_day_biochemical_oxygen_demand")
	PointMap.Set("硒", "se")
	PointMap.Set("砷", "as")
	PointMap.Set("汞", "hg")
	PointMap.Set("氟化物", "fluoride")
	PointMap.Set("石油类", "petroleum")
	PointMap.Set("粪大肠菌群", "fecalColiform")
	PointMap.Set("溶解氧", "dissolved_oxygen")
	PointMap.Set("电导率", "conductivity")
	PointMap.Set("浊度", "turbidity")
	PointMap.Set("高锰酸盐指数", "CODMII")
	PointMap.Set("氨氮", "ammonia_nitrogen")
	PointMap.Set("总磷", "total_phosphorus")
	PointMap.Set("总氮", "total_nitrogen")
	PointMap.Set("CODcr", "CODcr")
	PointMap.Set("氰化物", "cyanide")
	PointMap.Set("挥发酚", "volatile_phenol")
	PointMap.Set("六价铬", "hexavalent_chromium")
	PointMap.Set("铜", "cu")
	PointMap.Set("锌", "zn")
	PointMap.Set("铅", "pb")
	PointMap.Set("镉", "cd")
	PointMap.Set("阴离子表面活性剂", "LAS")
	PointMap.Set("硫化物", "SOx")
	PointMap.Set("累计流量", "cumulative_discharge")
	PointMap.Set("水流量", "water_discharge")
	PointMap.Set("总累积流量", "total_cumulative_flow")
	PointMap.Set("水位", "water_level")
	PointMap.Set("时段累积流量", "period_cumulative_Flow")
	PointMap.Set("断面平均流速", "sectional_mean_velocity")
	PointMap.Set("当前瞬时流速", "current_instantaneous_velocity")
	PointMap.Set("瞬时流量", "instantaneous_delivery")
	PointMap.Set("断面面积", "sectional_area")

	// 初始行唯一字段映射表
	RowOneMap.Set("水质类别", "water_quality_classification")
	RowOneMap.Set("主要污染物", "key_pollutant")
	RowOneMap.Set("定类项目", "classification_project")
	RowOneMap.Set("定类项目", "classification_project")

	// 初始化行多字段映射表
	RowAllMap.Set("最小值", "minimum_value")
	RowAllMap.Set("最大值", "maximum_value")
	RowAllMap.Set("平均值", "average_value")
	RowAllMap.Set("分项类别", "item_category")

	// 初始化制度映射表
	SysMap.Set("小时制", "hour")
	SysMap.Set("月度制", "month")

	// 文件内容的标记点映射表
	OptMap.Set("hour", "时间")
	OptMap.Set("month", "监测断面")

	// 站名标记映射表
	StationNameFlagMap.Set("flag1", "自动站名称")

	// 站名注册表
	StationMap.Set("海门湾桥闸", "haimen_bay_bridge_gate")
	StationMap.Set("汕头练江水站", "lian_jiang_water_station")
	StationMap.Set("青洋山桥", "lian_jiang_water_station")
	StationMap.Set("新溪西村", "xinxi_village")
	StationMap.Set("万兴桥", "wanxing_bridge")
	StationMap.Set("流仙学校", "liuxian_school")
	StationMap.Set("仙马闸", "xianma_brake")
	StationMap.Set("华侨学校", "huaqiao_school")
	StationMap.Set("港洲桥", "gangzhou_bridge")
	StationMap.Set("云陇", "yunlong")
	StationMap.Set("北港水", "beixiangshui")
	StationMap.Set("官田水", "guantianshui")
	StationMap.Set("北港河闸", "beixiang_penstock")
	StationMap.Set("峡山大溪", "xiashan_stream")
	StationMap.Set("井仔湾闸", "jingzai_wan_sluice")
	StationMap.Set("东北支流", "northeast_branch")
	StationMap.Set("西埔桥闸", "xipu_bridge_sluice")
	StationMap.Set("五福桥", "wufu_bridge")
	StationMap.Set("成田大寮", "narita_daliao")
	StationMap.Set("新坛港", "xitan_port")
	StationMap.Set("瑶池港", "yaochi_port")
	StationMap.Set("护城河闸", "moat_locks")
	StationMap.Set("和平桥", "peace_bridge")
}

// @title    InitMapMap
// @description   初始化映射表
// @param     void			没有入参
// @return    error			是否发生错误
func InitMapMap() error {
	db := common.GetDB()

	var mapVer model.MapVersion

	result_1 := db.Where("active = ?", 1).First(&mapVer)
	if result_1.Error != nil {
		if errors.Is(result_1.Error, gorm.ErrRecordNotFound) {
			log.Println("开始初始化映射表为默认值")
			InitDefaultMapMap()

			tx := db.Begin()

			mapVer.Active = true

			if err := tx.Create(&mapVer).Error; err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.
				Model(&model.MapVersion{}).
				Where("id = ?", mapVer.Id).
				Update("version_name", fmt.Sprintf("版本%d", mapVer.Id)).
				Error; err != nil {
					tx.Rollback()
					return err
			}

			var details []model.MapVersionDetail

			for table, v := range MapMap {
				for _, key := range v.Keys() {
					value, _ := v.Get(key)
					details = append(details, model.MapVersionDetail{
						Table:  table,
						Key:    key,
						VerId:  mapVer.Id,
						Value:  value.(string),
					})
				}
			}
	
			if err := tx.Create(&details).Error; err != nil {
				tx.Rollback()
				return err
			}

			return tx.Commit().Error
		} else {
			return result_1.Error
		}

		return nil
	}

	var details []model.MapVersionDetail

	result_2 :=  db.Where("ver_id = ?", mapVer.Id).Find(&details)

	if result_2.RowsAffected == 0 {

		log.Println("当前版本映射表发生错误，尝试使用映射表sql备份进行恢复")

		if err := HasBackUpSql(MapSqlPath, MapSqlFiles); err != nil {
			if errors.Is(err, ErrMissingBackupFile) {
					// 清除脏数据，重新初始化
					if err := db.Delete(&mapVer).Error; err != nil {
							return err
					}
					log.Println("未检测到映射表sql备份")
					return InitMapMap()
			}
			return err // 其他错误直接返回
		}

		log.Println("检测到映射表sql备份，开始使用映射表sql备份进行恢复...")

		if err := RecoverFromBackupSql(MapSqlPath, MapSqlFiles, true); err != nil {
			return err
		}

		log.Println("恢复成功")

		if err := db.Where("ver_id = ?", mapVer.Id).Find(&details).Error; err != nil {
			return err
		}
	} else if result_2.Error != nil {
		return result_2.Error
	}

	for _, v := range details {
		mp, ok := MapMap[v.Table]
		if !ok {
				return fmt.Errorf("映射表初始化错误：未知的映射类型 '%v'", v.Table)
		}
		mp.Set(v.Key, v.Value)
	}

	return nil
}
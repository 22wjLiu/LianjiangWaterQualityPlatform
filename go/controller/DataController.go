// @Title  DataController
// @Description  该文件用于提供操作数据的各种函数
package controller

import (
	"bytes"
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/util"
	"os/exec"
	"strconv"
	"time"
	"unicode"
	"log"

	"lianjiang/response"

	"github.com/gin-gonic/gin"
)

// @title    DeleteData
// @description   删除点集数据
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteData(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在四级以下的用户不能删除数据
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 获取path中的start
	start := ctx.Params.ByName("start")

	if start == "" {
		start = "2000-01-01"
	}

	// TODO 获取path中的end
	end := ctx.Params.ByName("end")

	if end == "" {
		end = time.Now().Format("2006-01-02")
	}

	// TODO 获取path中的time
	t := ctx.Params.ByName("time")

	time, ok := util.TimeMap[t]

	if !ok {
		response.Fail(ctx, nil, "时间字段"+t+"不存在")
		return
	}

	// TODO 取出请求
	sys := ctx.DefaultQuery("system", "")
	name := ctx.DefaultQuery("name", "")

	// TODO 尝试取出制度
	var system interface{}

	if sys != "" {
		if !util.SysMap.Has(sys) {
			response.Fail(ctx, nil, "时间制度"+sys+"不存在")
			return
		}
		system, _ = util.SysMap.Get(sys)
	} else {
		system = ""
	}

	// TODO 尝试取出站名
	var stationName interface{}

	if name != "" {
		if !util.StationMap.Has(name) {
			response.Fail(ctx, nil, "站名"+name+"不存在")
			return
		}
		stationName, _ = util.StationMap.Get(name)
	} else {
		stationName = ""
	}

	// TODO 组合数组
	systems, stationNames := make([]string, 0), make([]string, 0)

	// TODO 如果为空，取出所有值
	if stationName.(string) == "" {
		stationNames = util.StationMap.Keys()
		for i, v := range stationNames {
			s, _ := util.StationMap.Get(v)
			stationNames[i] = s.(string)
		}
	} else {
		stationNames = append(stationNames, stationName.(string))
	}

	if system.(string) == "" {
		systems = util.SysMap.Keys()
		for i, v := range systems {
			s, _ := util.SysMap.Get(v)
			systems[i] = s.(string)
		}
	} else {
		systems = append(systems, system.(string))
	}

	// TODO 删除对应数据
	db := common.GetDB()
	for _, sys := range systems {
		for _, sta := range stationNames {
			if db.Migrator().HasTable(sys + "_" + sta) {
				db.Table(sys+"_"+sta).Where(time+" >= ? and "+time+" <= ?", start, end).Delete(model.Point{})
			}
		}
	}
	// TODO 创建数据历史记录
	db.Create(&model.DataHistory{
		UserId:      user.Id,
		Option:      "删除",
		StartTime:   start,
		EndTime:     end,
		StationName: name,
		System:      sys,
		Time:        t,
	})

	response.Success(ctx, nil, "删除成功")
}

// @title    RecoverData
// @description   恢复点集数据
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func RecoverData(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在四级以下的用户不能删除数据
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 获取path中的start
	start := ctx.Params.ByName("start")

	if start == "" {
		start = "2000-01-01"
	}

	// TODO 获取path中的end
	end := ctx.Params.ByName("end")

	if end == "" {
		end = time.Now().Format("2006-01-02")
	}

	// TODO 获取path中的time
	t := ctx.Params.ByName("time")

	time, ok := util.TimeMap[t]

	if !ok {
		response.Fail(ctx, nil, "时间字段"+t+"不存在")
		return
	}

	// TODO 取出请求
	sys := ctx.DefaultQuery("system", "")
	name := ctx.DefaultQuery("name", "")

	// TODO 尝试取出制度
	var system interface{}

	if sys != "" {
		if !util.SysMap.Has(sys) {
			response.Fail(ctx, nil, "时间制度"+sys+"不存在")
			return
		}
		system, _ = util.SysMap.Get(sys)
	} else {
		system = ""
	}

	// TODO 尝试取出站名
	var stationName interface{}

	if name != "" {
		if !util.StationMap.Has(name) {
			response.Fail(ctx, nil, "站名"+name+"不存在")
			return
		}
		stationName, _ = util.StationMap.Get(name)
	} else {
		stationName = ""
	}

	// TODO 组合数组
	systems, stationNames := make([]string, 0), make([]string, 0)

	// TODO 如果为空，取出所有值
	if stationName.(string) == "" {
		stationNames = util.StationMap.Keys()
		for i, v := range stationNames {
			s, _ := util.StationMap.Get(v)
			stationNames[i] = s.(string)
		}
	} else {
		stationNames = append(stationNames, stationName.(string))
	}

	if system.(string) == "" {
		systems = util.SysMap.Keys()
		for i, v := range systems {
			s, _ := util.SysMap.Get(v)
			systems[i] = s.(string)
		}
	} else {
		systems = append(systems, system.(string))
	}

	// TODO 恢复对应数据
	db := common.GetDB()
	for _, sys := range systems {
		for _, sta := range stationNames {
			if db.Migrator().HasTable(sys + "_" + sta) {
				db.Table(sys+"_"+sta).Where(time+" >= ? and "+time+" <= ?", start, end).Update("deleted_at", nil)
			}
		}
	}
	response.Success(ctx, nil, "恢复成功")
}

// @title    ShowData
// @description   获取点集数据
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowData(ctx *gin.Context) {

	// TODO 获取path中的name
	n := ctx.Params.ByName("name")

	name, ok := util.StationMap.Get(n)

	if !ok {
		response.Fail(ctx, nil, "不存在站名"+n)
		return
	}

	// TODO 获取path中的system
	s := ctx.Params.ByName("system")

	system, ok := util.SysMap.Get(s)

	if !ok {
		response.Fail(ctx, nil, "不存在制度"+s)
		return
	}

	tableName := system.(string) + "_" + name.(string)

	// TODO 获取path中的fields
	f := ctx.QueryArray("fields")

	fields := make([]string, len(f), len(f))
	for i, v := range f {
		field, ok := util.PointMap.Get(v)
		if !ok {
			response.Fail(ctx, nil, "不存在字段"+v)
			return
		}
		fields[i] = util.StringToSql(field.(string))
	}

	log.Println("所有字段为：", fields)

	fields = append(fields, "time")

	db := common.GetDB()

	// TODO 查看是否存在该表
	if !db.Migrator().HasTable(tableName) {
		response.Fail(ctx, nil, "不存在对应表")
		return
	}

	// TODO 取出请求
	start := ctx.DefaultQuery("start", "2000-01-01")
	end := ctx.DefaultQuery("end", time.Now().Format("2006-01-02"))

	var total int64

	db.Table(tableName).Where("time >= ? and time <= ?", start, end).Count(&total)

	// TODO 查找对应数组

	resultArr := make([]map[string]interface{}, 0)

	db.Table(tableName).Select(fields).Where("time >= ? and time <= ?", start, end).Scan(&resultArr)

	response.Success(ctx, gin.H{"resultArr": resultArr}, "查找成功")
}

// @title    ShowRowAllData
// @description   获取一对多行字段点集数据
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowRowAllData(ctx *gin.Context) {

	// TODO 获取path中的name
	n := ctx.Params.ByName("name")

	name, ok := util.StationMap.Get(n)

	if !ok {
		response.Fail(ctx, nil, "不存在站名"+n)
		return
	}

	// TODO 获取path中的key
	k := ctx.Params.ByName("key")

	key, ok := util.RowAllMap.Get(k)

	if !ok {
		response.Fail(ctx, nil, "不存在字段"+k)
		return
	}

	// TODO 获取path中的fields
	f := ctx.QueryArray("fields")

	fields := make([]string, len(f), len(f))
	for i, v := range f {
		field, ok := util.PointMap.Get(v)
		if !ok {
			response.Fail(ctx, nil, "不存在字段"+v)
			return
		}
		fields[i] = util.StringToSql(field.(string))
	}

	fields = append(fields, "start_time")
	fields = append(fields, "end_time")

	db := common.GetDB()

	// TODO 查看是否存在该表
	if !db.Migrator().HasTable(key.(string)) {
		response.Fail(ctx, nil, "不存在对应表")
		return
	}

	// TODO 取出请求
	start := ctx.DefaultQuery("start", "2000-01-01")
	end := ctx.DefaultQuery("end", time.Now().Format("2006-01-02"))

	// TODO 搜索数据量
	var total int64

	db.Table(key.(string)).Where("start_time >= ? and end_time <= ?", start, end).Count(&total)

	// TODO 查找对应数组
	resultArr := make([]map[string]interface{}, 0)

	db.Table(key.(string)).Select(fields).Where("start_time >= ? and end_time <= ?", name.(string), start, end).Scan(&resultArr)

	response.Success(ctx, gin.H{"resultArr": resultArr}, "查找成功")
}

// @title    ShowRowOneData
// @description   获取一对一行字段点集数据
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ShowRowOneData(ctx *gin.Context) {

	// TODO 获取path中的name
	n := ctx.Params.ByName("name")

	name, ok := util.StationMap.Get(n)

	if !ok {
		response.Fail(ctx, nil, "不存在站名"+n)
		return
	}

	// TODO 获取path中的key
	k := ctx.Params.ByName("key")

	key, ok := util.RowAllMap.Get(k)

	if !ok {
		response.Fail(ctx, nil, "不存在字段"+k)
		return
	}

	db := common.GetDB()

	// TODO 查看是否存在该表
	if !db.Migrator().HasTable(key.(string)) {
		response.Fail(ctx, nil, "不存在对应表")
		return
	}

	// TODO 取出请求
	start := ctx.DefaultQuery("start", "2000-01-01")
	end := ctx.DefaultQuery("end", time.Now().Format("2006-01-02"))

	// TODO 搜索数据量
	var total int64

	db.Table(key.(string)).Where("start_time >= ? and end_time <= ?", start, end).Count(&total)

	// TODO 查找对应数组
	resultArr := make([]map[string]interface{}, 0)
	db.Table(key.(string)).Select([]string{"detail", "start_time", "end_time"}).Where("station_name = ? and start_time >= ? and end_time <= ?", name.(string), start, end).Scan(&resultArr)

	response.Success(ctx, gin.H{"resultArr": resultArr}, "查找成功")
}

// @title    Forecast
// @description   进行数据预测
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Forecast(ctx *gin.Context) {
	// TODO 读取数据
	Temperature := ctx.Query("Temperature")
	PH := ctx.Query("PH")
	Turbidity := ctx.Query("Turbidity")
	DO := ctx.Query("DO")

	if Temperature == "" || PH == "" || Turbidity == "" || DO == "" {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	// TODO python main.py
	cmd := exec.Command("python", "main.py", "--Temperature", Temperature, "--PH", PH, "--Turbidity", Turbidity, "--DO", DO)
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		response.Fail(ctx, nil, "参数错误")
		return
	}
	res := out.String()

	var data []float64

	before := 0

	for i, s := range []rune(res) {
		if s != '.' && !unicode.IsDigit(s) {
			data1, _ := strconv.ParseFloat(res[before:i], 64)
			before = i + 1
			data = append(data, data1)
		}
	}

	response.Success(ctx, gin.H{"data": data}, "查找成功")
}

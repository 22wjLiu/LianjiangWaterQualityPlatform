// @Title  PointController
// @Description  该文件用于提供操作点集的各种函数
package controller

import (
	"fmt"
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/util"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"

	"lianjiang/response"

	"github.com/gin-gonic/gin"
)

// @title    Upload
// @description   用户上传点集文件
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Upload(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在三级以下的用户无法上传文件
	if user.Level < 3 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	file, err := ctx.FormFile("file")

	//TODO 数据验证
	if err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 验证文件格式
	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".xls":  true,
		".xlsx": true,
		".csv":  true,
	}

	// TODO 格式验证
	if _, ok := allowExtMap[extName]; !ok {
		response.Fail(ctx, nil, "文件后缀有误")
		return
	}

	// TODO 从path中获取制度
	system := ctx.Params.ByName("system")

	// TODO 从制度表中获取制度映射
	sys, ok := util.SysMap.Get(system)

	// TODO 如果制度表中未注册
	if !ok {
		response.Fail(ctx, nil, "制度未注册")
		return
	}

	// TODO 从标记表中获取标记映射
	opt, ok := util.OptMap.Get(sys.(string))

	// TODO 如果未注册标记
	if !ok {
		response.Fail(ctx, nil, "标记未注册")
		return
	}

	// TODO 尝试建立对应文件夹
	err = util.Mkdir("./home/" + sys.(string))

	if err != nil {
		response.Fail(ctx, nil, "创建路径失败，系统错误")
		return
	}

	// TODO 将文件存入本地
	ctx.SaveUploadedFile(file, "./home/"+sys.(string)+"/"+file.Filename)

	// TODO 解析文件
	res, err := util.Read("./home/" + sys.(string) + "/" + file.Filename)

	// TODO 解析有误
	if err != nil || res == nil {
		response.Fail(ctx, nil, "文件解析有误")
		return
	}

	// TODO 用于存储字段映射序列
	index := make([]string, len(res[0]))

	// TODO start表示数据的起始行数
	start := 0

	// TODO flag 用于标记是否遇到标记
	flag := false

	// TODO 初始时间和终止时间
	startTime, endTime := time.Unix(int64((100000-25569)*24*60*60)-8*60*60, 0), time.Unix(0, 0)

	// TODO 用于建立数据库表的模板point
	var point model.Point
	point.System = sys.(string)

	// TODO 用于存储站名
	var stName string

	// TODO 逐行遍历，尝试寻找站名并取出字段映射
	for i := 0; i < len(res); i++ {
		for j := 0; j < len(res[i]); j++ {
			// TODO 成功找到站名
			if len(res[i][j]) > 18 && res[i][j][0:18] == "自动站名称：" {
				stName = res[i][j][18:]
				continue
			}
			p := ""
			// TODO 寻找最长前缀匹配
			for k := 1; k <= len(res[i][j]); k++ {
				str, ok := util.PointMap.Get(res[i][j][0:k])
				if ok {
					p = str.(string)
				}
			}
			// TODO 成功匹配映射字段，则记录该字段
			if p != "" {
				index[j] = p
			}
			// TODO 遇到标记
			if res[i][j] == opt.(string) {
				flag = true
			}
		}
		// TODO 如果遇到标记，记录数据初始位置，并退出字段搜寻
		if flag {
			start = i + 1
			break
		}
	}

	// TODO 未找到标记
	if !flag {
		response.Fail(ctx, nil, "文件内容缺少标记")
		return
	}

	// TODO 获取数据库指针
	db := common.GetDB()

	// TODO 一行一行的遍历数据，将遍历到的数据存入数据库
	for i := start; i < len(res); i++ {

		var p model.Point

		// 遍历每一列，尝试取出数据
		for j := 0; j < len(res[i]); j++ {
			row, ok := util.RowOneMap.Get(res[i][j])
			// 如果是唯一字段
			if ok {
				// TODO 时间有误
				if endTime.Before(startTime) {
					break
				}
				var rowOne model.RowOne
				// TODO 存入该字段
				rowOne.Detail = res[i][j+1]
				// TODO 存入时间
				rowOne.StartTime = startTime
				rowOne.EndTime = endTime
				// TODO 存入站名
				rowOne.StationName = point.StationName
				// TODO 查看是否存在表
				if !db.Migrator().HasTable(row.(string)) {
					// TODO 在第一次存入数据前，先尝试建立数据表
					db.AutoMigrate(&rowOne)
					// TODO 表名修正
					db.Migrator().RenameTable(&rowOne, row.(string))
				}
				// TODO 存入数据库
				db.Table(row.(string)).Create(&rowOne)
				break
			}

			row, ok = util.RowAllMap.Get(res[i][j])
			// 如果是多字段
			if ok {
				// TODO 时间有误
				if endTime.Before(startTime) {
					break
				}
				var rowAll model.RowAll
				// TODO 存入时间
				rowAll.StartTime = startTime
				rowAll.EndTime = endTime
				// TODO 存入站名
				rowAll.StationName = point.StationName
				for k := j + 1; k < len(res[i]); k++ {
					if index[k] == "" {
						continue
					}
					// TODO 利用反射机制写入结构体
					reflect.ValueOf(&rowAll).Elem().FieldByName(index[k]).SetString(res[i][k])
				}
				// TODO 查看是否存在表
				if !db.Migrator().HasTable(row.(string)) {
					// TODO 在第一次存入数据前，先尝试建立数据表
					db.AutoMigrate(&rowAll)
					// TODO 表名修正
					db.Migrator().RenameTable(&rowAll, row.(string))
				}
				// TODO 存入数据库
				db.Table(row.(string)).Create(&rowAll)
				break
			}

			// TODO 如果该列没有字段
			if j >= len(index) || index[j] == "" {
				continue
			}

			tp, _ := reflect.TypeOf(p).FieldByName(index[j])

			// TODO 利用反射机制判断结构体字段类型
			switch tp.Type.String() {
			case "string":
				// TODO 利用反射机制写入结构体
				reflect.ValueOf(&p).Elem().FieldByName(index[j]).SetString(res[i][j])

			case "float64":

				// TODO 尝试取出数字
				data, ok := util.StringToFloat(res[i][j])
				// TODO 成功取出数字
				if ok {
					// TODO 利用反射机制写入结构体
					reflect.ValueOf(&p).Elem().FieldByName(index[j]).SetFloat(data)
				}
			case "time.Time":
				// TODO 成功取出数字
				data, err := strconv.ParseFloat(res[i][j], 64)
				// TODO 如果出现了数据读出损坏，尝试修复数据
				if err != nil || data < 40000.0 || data > 60000.0 {
					// TODO 如果是递增或者递减，则测算出损坏数据
					if i > start+3 {
						var t1, t2, t3 float64

						// TODO 取出前三位数据
						t1, err = strconv.ParseFloat(res[i-1][0], 64)
						if err != nil {
							continue
						}

						t2, err = strconv.ParseFloat(res[i-2][0], 64)
						if err != nil {
							continue
						}

						t3, err = strconv.ParseFloat(res[i-3][0], 64)
						if err != nil {
							continue
						}

						// TODO 不满足递增或者递减，滤过这条数据
						if (t3-t2)-(t2-t1) > 0.001 {
							continue
						}

						// TODO 满足则计算预测值
						data = t1 + t1 - t2
					}
					if i < len(res)-3 {
						var t1, t2, t3 float64
						// TODO 取出后三位数据
						t1, err = strconv.ParseFloat(res[i+1][0], 64)
						if err != nil {
							continue
						}

						t2, err = strconv.ParseFloat(res[i+2][0], 64)
						if err != nil {
							continue
						}

						t3, err = strconv.ParseFloat(res[i+3][0], 64)
						if err != nil {
							continue
						}

						// TODO 不满足递增或者递减，滤过这条数据
						if (t3-t2)-(t2-t1) > 0.001 {
							continue
						}

						// TODO 满足则计算处预测值
						data = t1 - (t2 - t1)
					} else {
						continue
					}
				}
				// TODO 计算正确时间
				reflect.ValueOf(&p).Elem().FieldByName(index[j]).Set(reflect.ValueOf(time.Unix(int64((data-25569)*24*60*60)-8*60*60, 0)))
			}
		}
		// TODO 查看第一次取数据是否找到站名
		if flag {
			if stName == "" {
				if p.StationName != "" {
					stName = p.StationName
				} else {
					response.Fail(ctx, nil, "未能在文件内找到站名")
					return
				}
			}
			// TODO 如果站名没有注册
			if !util.StationMap.Has(stName) {
				response.Fail(ctx, nil, "站名"+stName+"未注册")
				return
			}
			st, _ := util.StationMap.Get(stName)
			point.StationName = st.(string)
			// TODO 查看是否存在表
			if !db.Migrator().HasTable(point.System + "_" + point.StationName) {
				// TODO 在第一次存入数据前，先尝试建立数据表
				db.AutoMigrate(&point)
				// TODO 表名修正
				db.Migrator().RenameTable(&point, point.System+"_"+point.StationName)
			}
			flag = false
		}
		// TODO 时间错误
		if p.Time.Before(time.Unix(int64((40000.0-25569)*24*60*60)-8*60*60, 0)) || time.Unix(int64((60000.0-25569)*24*60*60)-8*60*60, 0).Before(p.Time) {
			fmt.Printf("第%d行的时间有误\n", i+1)
			continue
		}
		// TODO 更新初始时间和终止时间
		if p.Time.Before(startTime) {
			startTime = p.Time
		}
		if endTime.Before(p.Time) {
			endTime = p.Time
		}
		// TODO 存入数据库
		db.Table(point.System + "_" + point.StationName).Create(&p)
	}
	// TODO 创建文件历史记录
	db.Create(&model.FileHistory{
		UserId:   user.ID,
		FileName: file.Filename,
		FilePath: "/" + sys.(string) + "/" + file.Filename,
		Option:   "创建",
	})
	// TODO 创建数据历史记录
	db.Create(&model.DataHistory{
		UserId:      user.ID,
		Option:      "创建",
		StartTime:   startTime.String(),
		EndTime:     endTime.String(),
		StationName: stName,
		System:      system,
	})
	response.Success(ctx, gin.H{"FileName": file.Filename}, "更新成功")
}

// @title    List
// @description   提供点集文件列表
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func List(ctx *gin.Context) {

	// TODO 取出请求
	path := ctx.DefaultQuery("path", "/")

	// TODO 获得hour目录下的所有文件
	files, err := util.GetFiles(path)

	if err != nil {
		response.Fail(ctx, nil, "不存在该文件夹")
		return
	}

	response.Success(ctx, gin.H{"files": files}, "请求成功")

}

// @title    Download
// @description   下载点集文件
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Download(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在二级以下的用户不能下载文件
	if user.Level < 2 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 取出请求
	path := ctx.DefaultQuery("path", "/")
	file := ctx.DefaultQuery("file", "")

	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))
	ctx.File("./home" + path)
	response.Success(ctx, nil, "请求成功")
}

// @title    DeleteFile
// @description   删除点集文件
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func DeleteFile(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	// TODO 安全等级在四级以下的用户不能删除文件
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 取出请求
	path := ctx.DefaultQuery("path", "")

	// TODO 移除文件
	if os.Remove(path) != nil {
		response.Fail(ctx, nil, "路径不存在")
		return
	}

	// TODO 创建文件历史记录
	common.GetDB().Create(model.FileHistory{
		UserId:   user.ID,
		FilePath: path,
		Option:   "删除",
	})

	response.Success(ctx, nil, "删除成功")
}

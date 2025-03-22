// @Title  util
// @Description  收集各种需要使用的工具函数
package util

import (
	"context"
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/vo"
	"log"
	"math"
	"math/rand"
	"net/smtp"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	Map "github.com/orcaman/concurrent-map"

	"os"

	"github.com/tealeg/xlsx"

	"github.com/extrame/xls"
	"github.com/jordan-wright/email"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 点集字段映射表
var PointMap Map.ConcurrentMap = Map.New()

// 行唯一字段映射表
var RowOneMap Map.ConcurrentMap = Map.New()

// 行多字段映射表
var RowAllMap Map.ConcurrentMap = Map.New()

// 制度映射表
var SysMap Map.ConcurrentMap = Map.New()

// 文件内容的标记点映射表
var OptMap Map.ConcurrentMap = Map.New()

// 站名注册表
var StationMap Map.ConcurrentMap = Map.New()

// 数据注册表
var DataMap Map.ConcurrentMap = Map.New()

// 用户字段映射表
var UserMap = map[string]string{
	"名称": "name",
	"邮箱": "email",
	"ID": "id",
}

// 映射表
var MapMap = map[string]*Map.ConcurrentMap{
	"列字段映射":    &PointMap,
	"行字段一对多映射": &RowAllMap,
	"行字段一对一映射": &RowOneMap,
	"时间制映射":    &SysMap,
	"标记映射":     &OptMap,
	"站名映射":     &StationMap,
	"数据符号映射":   &DataMap,
}

// 时间映射表
var TimeMap = map[string]string{
	"创建日期": "created_at",
	"记录日期": "time",
	"删除日期": "deleted_at",
}

func init() {
	// TODO 初始化点集字段映射表
	PointMap.Set("监测断面", "StationName")
	PointMap.Set("监测指标", "Time")
	PointMap.Set("监测时间", "Time")
	PointMap.Set("时间", "Time")
	PointMap.Set("水温", "Temperature")
	PointMap.Set("pH", "PH")
	PointMap.Set("化学需氧量", "Cod")
	PointMap.Set("五日生化需氧量", "FiveDaysNiochemicalOxygenDemand")
	PointMap.Set("硒", "Se")
	PointMap.Set("砷", "As")
	PointMap.Set("汞", "Hg")
	PointMap.Set("氟化物", "Fluoride")
	PointMap.Set("石油类", "Petroleum")
	PointMap.Set("粪大肠菌群", "FecalColiform")
	PointMap.Set("溶解氧", "DO")
	PointMap.Set("电导率", "EC")
	PointMap.Set("浊度", "Turbidity")
	PointMap.Set("高锰酸盐指数", "CODMII")
	PointMap.Set("氨氮", "NH_N")
	PointMap.Set("总磷", "TP")
	PointMap.Set("总氮", "TN")
	PointMap.Set("CODcr", "CODcr")
	PointMap.Set("氰化物", "CN")
	PointMap.Set("挥发酚", "VolatilePenol")
	PointMap.Set("六价铬", "Cr")
	PointMap.Set("铜", "Cu")
	PointMap.Set("锌", "Zn")
	PointMap.Set("铅", "Pb")
	PointMap.Set("镉", "Cd")
	PointMap.Set("阴离子表面活性剂", "LAS")
	PointMap.Set("硫化物", "SOx")
	PointMap.Set("累计流量", "CumulativeDischarge")
	PointMap.Set("水流量", "WaterDischarge")
	PointMap.Set("总累积流量", "TotalCumulativeFlow")
	PointMap.Set("水位", "WaterLevel")
	PointMap.Set("时段累积流量", "PeriodCumulativeFlow")
	PointMap.Set("断面平均流速", "SectionalMeanVelocity")
	PointMap.Set("当前瞬时流速", "CurrentInstantaneousVelocity")
	PointMap.Set("瞬时流量", "InstantaneousDelivery")
	PointMap.Set("断面面积", "SectionalArea")

	// TODO 初始行唯一字段映射表
	RowOneMap.Set("水质类别", "water_quality_classification")
	RowOneMap.Set("主要污染物", "key_pollutant")

	// TODO 初始化行多字段映射表
	RowAllMap.Set("分项类别", "item_category")

	// TODO 初始化制度映射表
	SysMap.Set("小时制", "hour")
	SysMap.Set("月度制", "month")

	// TODO 文件内容的标记点映射表
	OptMap.Set("hour", "时间")
	OptMap.Set("month", "监测断面")

	// TODO 站名注册表
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

// @title    StringToSql
// @description   将model字段转化为数据库字段
// @param     point string			输入字符串
// @return    sql string			sql格式的字段
func StringToSql(point string) (sql string) {
	if len(point) <= 6 {
		return strings.ToLower(point)
	}

	var builder strings.Builder

	for i, val := range point {
		// TODO 是否大写
		if unicode.IsUpper(val) {
			if i > 0 {
				builder.WriteRune('_')
			}
		}
		builder.WriteRune(val)
	}
	return strings.ToLower(builder.String())
}

// @title    Read
// @description   读取文件内容
// @param     file_path string		文件位置
// @return    res [][]string, err error		res为读出的内容，err为可能出现的错误
func Read(file_path string) (res [][]string, err error) {

	extName := path.Ext(file_path)

	if extName == ".csv" {
		return ReadCsv(file_path)
	} else if extName == ".xls" {
		return ReadXls(file_path)
	} else if extName == ".xlsx" {
		return ReadXlsx(file_path)
	}
	return nil, nil
}

// @title    ReadCsv
// @description   读取Csv文件内容
// @param     file_path string		文件位置
// @return    res [][]string, err error		res为读出的内容，err为可能出现的错误
func ReadCsv(file_path string) (res [][]string, err error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// TODO 初始化csv-reader
	reader := csv.NewReader(file)
	// TODO 设置返回记录中每行数据期望的字段数，-1 表示返回所有字段
	reader.FieldsPerRecord = -1
	// TODO 允许懒引号（忘记遇到哪个问题才加的这行）
	reader.LazyQuotes = true
	// TODO 返回csv中的所有内容
	records, read_err := reader.ReadAll()
	if read_err != nil {
		return nil, read_err
	}
	return records, nil
}

// @title    ReadXls
// @description   读取Xls文件内容
// @param     file_path string		文件位置
// @return    res [][]string, err error		res为读出的内容，err为可能出现的错误
func ReadXls(file_path string) (res [][]string, err error) {
	if xlFile, err := xls.Open(file_path, "utf-8"); err == nil {
		fmt.Println(xlFile.Author)
		// TODO 第一个sheet
		sheet := xlFile.GetSheet(0)
		if sheet.MaxRow != 0 {
			temp := make([][]string, sheet.MaxRow)
			for i := 0; i < int(sheet.MaxRow); i++ {
				row := sheet.Row(i)
				data := make([]string, 0)
				if row.LastCol() > 0 {
					for j := 0; j < row.LastCol(); j++ {
						col := row.Col(j)
						data = append(data, col)
					}
					temp[i] = data
				}
			}
			res = append(res, temp...)
		}
	} else {
		return nil, err
	}
	return res, nil
}

// @title    ReadXlsx
// @description   读取Xlsx文件内容
// @param     file_path string		文件位置
// @return    res [][]string, err error		res为读出的内容，err为可能出现的错误
func ReadXlsx(file_path string) (res [][]string, err error) {
	if xlFile, err := xlsx.OpenFile(file_path); err == nil {
		for index, sheet := range xlFile.Sheets {
			// TODO 第一个sheet
			if index == 0 {
				temp := make([][]string, len(sheet.Rows))
				for k, row := range sheet.Rows {
					var data []string
					for _, cell := range row.Cells {
						data = append(data, cell.Value)
					}
					temp[k] = data
				}
				res = append(res, temp...)
			}
		}
	} else {
		return nil, err
	}
	return res, nil
}

// @title    SecondToTime
// @description   把秒级的时间戳转为time格式
// @param     sec int64	秒
// @return    time.Time    Time类型
func SecondToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}

// @title    GetFiles
// @description   获取一个目录下的所有文件
// @param     folder string	指定目录
// @return    []string    所有文件的文件名
func GetFiles(folder string) ([]vo.File, error) {
	files, err := ioutil.ReadDir("./home" + folder)
	if err != nil {
		return nil, err
	}
	res := make([]vo.File, 0)
	for _, file := range files {
		// TODO 尝试读出所有文件的相关信息
		var f vo.File
		f.Name = file.Name()
		if folder[len(folder)-1] == '/' {
			f.Path = folder + file.Name()
		} else {
			f.Path = folder + "/" + file.Name()
		}
		f.Size = file.Size()
		f.LastWriteTime = file.ModTime()
		if file.IsDir() {
			f.Type = "Dir"
		} else {
			f.Type = path.Ext(file.Name())
		}
		res = append(res, f)
	}
	return res, nil
}

// @title    PathExists
// @description   判断文件夹是否存在
// @param     path string	指定目录
// @return    bool, error    查看文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// @title    Mkdir
// @description   建立文件夹
// @param     path string	指定路径
// @return   error    查看是否出错
func Mkdir(dir string) error {
	exist, err := PathExists(dir)
	if err != nil {
		return err
	}

	if !exist {
		// TODO 创建目录
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Printf("创建目录失败：%v\n", err)
			return err
		}
	}
	return nil
}

// @title    StringToFloat
// @description   从字符串中提取各式各样的浮点数
// @param     s string		一串字符串
// @return    float64, bool		表示解析出来的浮点数，ok表示解析是否成功
func StringToFloat(s string) (float64, bool) {
	// TODO 优先查看数据注册表
	data, ok := DataMap.Get(s)
	if ok {
		return data.(float64), ok
	}
	k := len(s)
	// TODO 尝试取出前缀数字，以此来滤过符号单位
	for k >= 0 {
		_, err := strconv.ParseFloat(s[0:k], 64)
		if err != nil {
			k--
		} else {
			break
		}
	}
	// TODO 成功取出数字
	if k > 0 {
		data, err := strconv.ParseFloat(s[0:k], 64)
		if err != nil {
			return 0, false
		}
		// TODO 查看是否有科学计数法
		if k+4 <= len(s) && s[k:(k+4)] == "×10" {
			// TODO 尝试读出后缀数字
			data1, ok := StringToFloat(s[(k + 4):])
			if !ok {
				data1 = 0
			} else if data1 == 0 {
				data1 = 1
			}
			data *= math.Pow(10, data1)
		}
		return data, true
	}
	return 0, false
}

// @title    RandomString
// @description   生成一段随机的字符串
// @param     n int		字符串的长度
// @return    string    一串随机的字符串
func RandomString(n int) string {
	var letters = []byte("qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	// TODO 不断用随机字母填充字符串
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// @title    VerifyEmailFormat
// @description   用于验证邮箱格式是否正确的工具函数
// @param     email string		一串字符串，表示邮箱
// @return    bool    返回是否合法
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// @title    isEmailExist
// @description   查看email是否在数据库中存在
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	return user.ID != 0
}

// @title    isNameExist
// @description   查看name是否在数据库中存在
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsNameExist(db *gorm.DB, name string) bool {
	var user model.User
	db.Where("name = ?", name).First(&user)
	return user.ID != 0
}

var ctx context.Context = context.Background()

// @title    SendEmailValidate
// @description   发送验证邮件
// @param    em []string       接收一个邮箱字符串
// @return   string, error     返回验证码和error值
func SendEmailValidate(em []string) (string, error) {
	// 邮箱发件配置
	from := "studentclubmanager@qq.com"
	fromHeader := "watercleareyes <studentclubmanager@qq.com>"
	authCode := "vwigictgofdtdbaf"
	smtpServer := "smtp.qq.com:465"

	auth := smtp.PlainAuth("", from, authCode, "smtp.qq.com")

	tlsConfig := &tls.Config{
		InsecureSkipVerify:	true,
		ServerName:		"smtp.qq.com",
	}

	mod := `From: %s
To: %s

	尊敬的%s，您好！

	您于 %s 提交的邮箱验证，本次验证码为%s，为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。
	此邮箱为系统邮箱，请勿回复。
`
	// 连接 SMTP 服务器
	conn, err := tls.Dial("tcp", smtpServer, tlsConfig)
	if err != nil {
		log.Printf("TLS 连接失败：%v\n", err)
		return "", err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, "smtp.qq.com")
	if err != nil {
		log.Printf("SMTP 客户端创建失败：%v\n", err)
		return "", err
	}
	defer client.Quit()

	// 使用 auth 进行认证
	if err = client.Auth(auth); err != nil {
		log.Printf("认证失败：%v\n", err)
		return  "", err
	}

	// 设置发件人和收件人
	if err = client.Mail(from); err != nil {
		log.Printf("发件人设置失败：%v\n", err)
		return "", err
	}
	if err = client.Rcpt(em[0]); err != nil {
		log.Printf("收件人设置失败：%v\n", err)
		return "" , err
	}

	// 写入邮件内容
	wc, err := client.Data()
	if err != nil{
		log.Printf("数据写入失败：%v\n", err)
		return "", err
	}
	defer wc.Close()
	//  生成6位随机验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	t := time.Now().Format("2006-01-02 15:04:05i")
	//  设置文件发送的内容
	content := fmt.Sprintf(mod, fromHeader, em[0], em[0], t, vCode)
	_, err = wc.Write([]byte(content))
	if err != nil{
		log.Printf("消息发生失败：%v", err)
	}

	log.Printf("消息发送成功")

	return vCode, err
}

// @title    SendEmailPass
// @description   发送密码邮件
// @param    em []string       接收一个邮箱字符串
// @return   string, error     返回验证码和error值
func SendEmailPass(em []string) string {
	mod := `
	尊敬的%s，您好！

	您于 %s 提交的邮箱验证，已经将密码重置为%s，为了保证账号安全。切勿向他人泄露，并尽快更改密码，感谢您的理解与使用。
	此邮箱为系统邮箱，请勿回复。
`
	e := email.NewEmail()
	e.From = "watercleareyes <studentclubmanager@qq.com>"
	e.To = em
	// TODO 生成8位随机密码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := fmt.Sprintf("%08v", rnd.Int31n(100000000))
	t := time.Now().Format("2006-01-02 15:04:05")

	db := common.GetDB()

	// TODO 创建密码哈希
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "密码加密失败"
	}

	// TODO 更新密码
	err = db.Model(&model.User{}).Where("email = ?", em[0]).Updates(model.User{
		Password: string(hasedPassword),
	}).Error

	if err != nil {
		return "密码更新失败"
	}

	// TODO 设置文件发送的内容
	content := fmt.Sprintf(mod, em[0], t, password)
	e.Text = []byte(content)
	// TODO 设置服务器相关的配置
	err = e.Send("smtp.qq.com:25", smtp.PlainAuth("", "studentclubmanager@qq.com", "gcovigagucxncjih", "smtp.qq.com"))

	if err != nil {
		return "邮件发送失败"
	}

	return "密码已重置"
}

// @title    IsEmailPass
// @description   验证邮箱是否通过
// @param    em []string       接收一个邮箱字符串
// @return   string, error     返回验证码和error值
func IsEmailPass(email string, vertify string) bool {
	client := common.GetRedisClient(0)
	V, err := client.Get(ctx, email).Result()
	if err != nil {
		return false
	}
	return V == vertify
}

// @title    SetRedisEmail
// @description   设置验证码，并令其存活五分钟
// @param    email string, v string       接收一个邮箱和一个验证码
// @return   void
func SetRedisEmail(email string, v string) {
	client := common.GetRedisClient(0)

	err := client.Set(ctx, email, v, 300*time.Second).Err()
	if err != nil {
		log.Println("Redis存储验证码失败：", err)
	}
}

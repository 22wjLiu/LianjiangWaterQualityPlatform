// @Title  commonUtils
// @Description  各种需要使用的其他工具函数
package util

import (
	"context"
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"lianjiang/common"
	"lianjiang/model"
	"log"
	"math"
	"math/rand"
	"net/smtp"
	"path"
	"regexp"
	"strconv"
	"time"
	"os"

	"github.com/tealeg/xlsx"

	"github.com/extrame/xls"
	"github.com/jordan-wright/email"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 可读时间格式定义
var ReadableTimeFormat = "2006-01-02 15:04:05"

// 预测时间格式定义
var ForecastTimeFormat = "2006-01-02 15:04"

// 转换Excel时间格式定义
var ExcelTimeFormat = "2006/1/2 15:04"

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
						if cell.Type() == xlsx.CellTypeDate {
							t, err := cell.GetTime(false) // 精度到秒
							if err == nil {
								val := t.Format(ExcelTimeFormat)
								data = append(data, val)
								continue
							}
						}
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
	return user.Id != 0
}

// @title    isNameExist
// @description   查看name是否在数据库中存在
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsNameExist(db *gorm.DB, name string) bool {
	var user model.User
	db.Where("name = ?", name).First(&user)
	return user.Id != 0
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

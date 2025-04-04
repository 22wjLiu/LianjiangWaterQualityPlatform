// @Title  UserController
// @Description  该文件用于提供操作用户的各种方法
package controller

import (
	"lianjiang/common"
	"lianjiang/model"
	"lianjiang/response"
	"lianjiang/util"
	"lianjiang/vo"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @title    Register
// @description   用户注册
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Register(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = vo.UserRequest{}
	ctx.Bind(&requestUser)
	// TODO 获取参数
	email := requestUser.Email
	password := requestUser.Password
	name := requestUser.Name
	// TODO 数据验证
	if !util.VerifyEmailFormat(email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, 201, 201, nil, "密码不能少于6位")
		return
	}

	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, email, password)

	// TODO 判断email是否存在
	if util.IsEmailExist(DB, email) {
		response.Response(ctx, 201, 201, nil, "用户已经存在")
		return
	}

	// TODO 判断email是否通过验证
	if !util.IsEmailPass(email, requestUser.Verify) {
		response.Response(ctx, 201, 201, nil, "邮箱验证码错误")
		return
	}

	if util.IsNameExist(DB, name) {
		response.Response(ctx, 201, 201, nil, "用户名称已被注册")
		return
	}

	// TODO 创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "加密错误")
		return
	}
	newUser := model.User{
		Email:    email,
		Password: string(hasedPassword),
		Level:    1,
	}
	DB.Create(&newUser)

	// TODO 发放token给前端
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	// TODO 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

// @title    Login
// @description   用户登录
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Login(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	ctx.Bind(&requestUser)
	// TODO 获取参数
	email := requestUser.Email
	password := requestUser.Password
	// TODO 数据验证
	if !util.VerifyEmailFormat(email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, 201, 201, nil, "密码不能少于6位")
		return
	}
	// TODO 判断邮箱是否存在
	var user model.User
	DB.Where("email = ?", email).First(&user)
	if user.Id == 0 {
		response.Response(ctx, 201, 201, nil, "用户不存在")
		return
	}
	// TODO 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil{
		response.Fail(ctx, nil, "密码错误")
		return
	}
	// TODO 发放token给前端
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	// TODO 返回结果
	response.Success(ctx, gin.H{
		"token": token,
		"level": user.Level, 
	}, "登录成功")
}

// @title    Security
// @description   进行密码找回的函数
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Security(ctx *gin.Context) {
	// TODO 数据验证
	DB := common.GetDB()
	var requestUser = vo.UserRequest{}
	ctx.Bind(&requestUser)
	if !util.VerifyEmailFormat(requestUser.Email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	// TODO 判断email是否存在
	if !util.IsEmailExist(DB, requestUser.Email) {
		response.Response(ctx, 201, 201, nil, "用户不存在")
		return
	}

	// TODO 判断email是否通过验证
	if !util.IsEmailPass(requestUser.Email, requestUser.Verify) {
		response.Response(ctx, 201, 201, nil, "邮箱验证码错误")
		return
	}

	err := util.SendEmailPass([]string{requestUser.Email})

	// TODO 返回结果
	response.Success(ctx, nil, err)
}

// @title    VerifyEmail
// @description   进行邮箱验证码发送的函数
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func VerifyEmail(ctx *gin.Context) {
	email := ctx.Params.ByName("id")
	// TODO 数据验证
	if !util.VerifyEmailFormat(email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	v, err := util.SendEmailValidate([]string{email})
	if err != nil {
		response.Response(ctx, 201, 201, nil, "邮箱验证码发送失败")
		return
	}
	// 验证码存入redis 并设置过期时间5分钟
	util.SetRedisEmail(email, v)

	// TODO 返回结果
	response.Success(ctx, gin.H{"email": email}, "验证码请求成功")
}

// @title    UpdatePass
// @description   进行密码修改的函数
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func UpdatePass(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var pairString = vo.PairString{}
	ctx.Bind(&pairString)

	// TODO 获取参数
	oldPass := pairString.First
	newPass := pairString.Second

	// TODO 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPass)); err != nil {
		response.Fail(ctx, nil, "密码错误")
		return
	}

	// TODO 创建密码哈希
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)

	if err != nil {
		response.Response(ctx, 201, 201, nil, "加密错误")
		return
	}

	db := common.GetDB()

	// TODO 更新密码
	user.Password = string(hasedPassword)

	db.Save(&user)

	response.Success(ctx, nil, "密码修改成功")
}

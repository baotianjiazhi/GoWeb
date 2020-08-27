package service

import (
	"bluebell/dao/mysql"
	"bluebell/model"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

type SignUpService struct {
	Username string	`form:"username" json:"username" binding:"required,max=20,min=8"`
	Password string	`form:"password" json:"password" binding:"required,max=18,min=8"`
	RePassword string `form:"re_password" json:"re_password" binding:"required,max=18,min=8,eqfield=Password"`
}

type SignInService struct {
	Username string	`form:"username" json:"username" binding:"required,max=20,min=8"`
	Password string `form:"password" json:"password" binding:"required,max=18,min=8"`
}

func (servicer *SignUpService)SignUp() (user *model.User, err error) {
	// 判断用户存不存在
	err = mysql.QueryUserByUsername(servicer.Username)
	if err != nil {
		return &model.User{}, err
	}
	// 生成UID
	userID := snowflake.GenID()
	// 构造一个User实例
	user = &model.User{
		UserID: userID,
		Username: servicer.Username,
		Password: servicer.Password,
	}
	// 保存进数据库
	u, err := mysql.CreateUser(user)
	if err != nil {
		return &model.User{}, err
	}
	return u, err
}

func (servicer *SignInService) Login() (user *model.User, err error) {
	// 通过用户名判断用户存不存在
	user = &model.User{
		Username: servicer.Username,
		Password: servicer.Password,
	}
	// 传递的是指针，可以拿到UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成JWT Token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}


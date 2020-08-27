package controller

import (
	"bluebell/dao/mysql"
	"bluebell/serializer"
	"bluebell/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	var servicer service.SignUpService
	// 1. 获取参数和参数校验
	if err := c.ShouldBindJSON(&servicer); err == nil {
		// 2. 业务处理
		if user, err := servicer.SignUp(); err != nil {
			if errors.Is(err, mysql.ErrUserNotExist){
				serializer.ResponseError(c, serializer.CodeUserExist)
				return
			}
			serializer.ResponseError(c, serializer.CodeCreateUserFault)
			return
		} else {
			// 3. 返回响应
			serializer.ResponseSuccess(c, user)
		}
	} else {
		errs, ok := err.(validator.ValidationErrors)
		if !ok{
			serializer.ResponseError(c, serializer.CodeInvalidParam)
			return
		}
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		serializer.ResponseErrorWithMsg(c, serializer.CodeInvalidPassword, removeTopStruct(errs.Translate(trans)))
		return
	}

}


func SignInHandler(c *gin.Context) {
	var servicer service.SignInService
	if err := c.ShouldBindJSON(&servicer); err == nil {
		// 业务逻辑
		if user, err := servicer.Login(); err != nil {
			zap.L().Error("logic.login failed", zap.String("username", servicer.Password), zap.Error(err))
			if errors.Is(err, mysql.ErrUserNotExist) {
				serializer.ResponseError(c, serializer.CodeUserNotExist)
				return
			}
			serializer.ResponseError(c, serializer.CodeInvalidPassword)
			return
		} else {
			serializer.ResponseSuccess(c, gin.H{
				"username": user.Username, // id值>53-1 int64类型的最大值是63-1
				"user_id": fmt.Sprintf("%d", user.UserID),
				"token": user.Token,
			})
		}
	} else {
		errs, ok := err.(validator.ValidationErrors)
		if !ok{
			serializer.ResponseError(c, serializer.CodeInvalidParam)
		}
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		serializer.ResponseErrorWithMsg(c, serializer.CodeInvalidPassword, removeTopStruct(errs.Translate(trans)))
		return
	}
}
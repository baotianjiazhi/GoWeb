package controller

import (
	"bluebell/serializer"
	"bluebell/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 投票
func PostVoteController(c *gin.Context) {
	// 参数校验
	var servicer service.VoteData
	if err := c.ShouldBindJSON(&servicer); err != nil {
		errs, ok := err.(validator.ValidationErrors)  // 类型断言
		if !ok {
			serializer.ResponseError(c, serializer.CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))  // 翻译并去除掉错误提示中的结构体标识
		serializer.ResponseErrorWithMsg(c, serializer.CodeInvalidParam, errData)
		return
	}
	// 具体实现的业务逻辑
	userID, err := GetCurrentUser(c)
	if err != nil {
		serializer.ResponseError(c, serializer.CodeNeedLogin)
		return
	}
	if err := service.VoteForPost(userID, &servicer); err != nil {
		zap.L().Error("service.VoteForPost(userID, &servicer) failed", zap.Error(err))
		serializer.ResponseError(c, serializer.CodeServerBusy)
		return
	}
	serializer.ResponseSuccess(c, nil)
}

package serializer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Data interface{} `json:"data"`
	Err interface{} `json:"err"`
	Status Rescode `json:"status"`
	Msg interface{} `json:"msg"`
}

func ResponseError(c *gin.Context, code Rescode) {
	c.JSON(http.StatusOK, &Response{
		Status: code,
		Msg: code.getMsg(),
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		Status: CodeSuccess,
		Msg: CodeSuccess.getMsg(),
		Data: data,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code Rescode, msg interface{}) {
	c.JSON(http.StatusOK, &Response{
		Status: code,
		Msg: msg,
	})
}
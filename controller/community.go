package controller

import (
	"bluebell/serializer"
	"bluebell/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// -- 跟社区相关
// CommunityHandler 查询所有社区
func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区(community_id, community_name) 以列表的形式返回
	data, err := service.GetCommunityList()
	if err != nil {
		zap.L().Error("service.GetCommunityList() failed", zap.Error(err))
		serializer.ResponseError(c, serializer.CodeServerBusy)
		return
	}
	serializer.ResponseSuccess(c, data)
}

// CommunityDetailHandler 根据ID查询社区详情
func CommunityDetailHandler(c *gin.Context) {
	community_id := c.Param("id")
	id, err := strconv.ParseInt(community_id, 10, 64)
	if err != nil {
	serializer.ResponseError(c, serializer.CodeInvalidParam)
	return
	}

	data, err := service.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("service.GetCommunityList() failed", zap.Error(err))
		serializer.ResponseError(c, serializer.CodeServerBusy)
		return
	}
	serializer.ResponseSuccess(c, data)

}

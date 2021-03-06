package controller

import (
	"bluebell/model"
	"bluebell/serializer"
	"bluebell/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)


// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	p := new(model.Post)
	// 1. 获取参数及参数请求
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p)", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		serializer.ResponseError(c, serializer.CodeInvalidParam)
		return
	}
	// 从c中取到当前发请求的用户ID值
	userID, err := GetCurrentUser(c)
	if err != nil {
		serializer.ResponseError(c, serializer.CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2. 创建帖子
	if err := service.CreatePost(p); err != nil {
		zap.L().Error("service.createPost(p) failed", zap.Error(err))
		serializer.ResponseError(c, serializer.CodeServerBusy)
		return
	}
	// 3.返回响应
	serializer.ResponseSuccess(c, nil)
}

// GetPostHandler 获取帖子详情
func GetPostHandler(c *gin.Context) {
	// 1. 获取参数
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		serializer.ResponseError(c, serializer.CodeInvalidParam)
		return
	}
	// 2. 根据id取出帖子数据（查数据库）
	data, err := service.GetPostByID(pid)
	if err != nil {
		zap.L().Error("service.GetPostByID(pid) failed", zap.Error(err))
		serializer.ResponseError(c, serializer.CodeServerBusy)
		return
	}
	fmt.Println(data)
	// 3. 返回响应
	serializer.ResponseSuccess(c, data)
}

// GetPostList 获取帖子列表函数
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := GetPageInfo(c)
	// 获取数据
	data, err := service.GetPostList(page, size)
	if err != nil {
		zap.L().Error("service.GetPostListHandler() err", zap.Error(err))
		serializer.ResponseError(c, serializer.CodeServerBusy)
		return
	}

	// 返回响应
	serializer.ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表函数
// 根据前端传来的参数（按分数或按创建时间排序）动态的获取帖子列表
// 1.获取参数
// 2.去redis查询id列表
// 3.根据id去数据库查询帖子详细信息
func GetPostListHandler2(c *gin.Context) {
	var servicer service.PostListService
	servicer = service.PostListService{
		Page: 1,
		Size: 10,
		Order: model.OrderTime,
	}
	// GET请求参数:/api/v1/posts2?page=1&size=10&order=time
	if err := c.ShouldBindQuery(&servicer); err != nil {
		zap.L().Error("c.ShouldBindQuery(&servicer) err", zap.Error(err))
		serializer.ResponseError(c, serializer.CodeInvalidParam)
		return
	}
	data, err := service.GetPostListNew(&servicer) // 更新：合二为一
	if err != nil {
		zap.L().Error("service.GetPostListHandler() err", zap.Error(err))
		serializer.ResponseError(c, serializer.CodeServerBusy)
		return
	}

	// 返回响应
	serializer.ResponseSuccess(c, data)
}

// 根据社区去查询帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	var servicer service.ParamCommunityPostListService
//	servicer = service.ParamCommunityPostListService{
//		PostListService: &service.PostListService{
//			Page: 1,
//			Size: 10,
//			Order: model.OrderTime,
//		},
//	}
//	// GET请求参数:/api/v1/posts2?page=1&size=10&order=time
//	if err := c.ShouldBindQuery(&servicer); err != nil {
//		zap.L().Error("c.ShouldBindQuery(&servicer) err", zap.Error(err))
//		serializer.ResponseError(c, serializer.CodeInvalidParam)
//		return
//	}
//
//	data, err := service.GetCommunityPostList2(&servicer)
//	if err != nil {
//		zap.L().Error("service.GetPostListHandler() err", zap.Error(err))
//		serializer.ResponseError(c, serializer.CodeServerBusy)
//		return
//	}
//
//	// 返回响应
//	serializer.ResponseSuccess(c, data)
//}

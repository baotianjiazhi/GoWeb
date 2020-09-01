package routers

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUp() *gin.Engine {
	r := gin.New()
	r.LoadHTMLGlob("templates/**/*")
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 注册相关路由信息
	r.GET("/", controller.Index)
	// 注册业务路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", controller.SignUpHandler)
		v1.POST("/login", controller.SignInHandler)

		auth := v1.Group("/")
		auth.Use(middleware.JWTAuthMiddleware())
		{
			auth.GET("/info", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"msg": "pangpangpang",
				})
			})
			auth.GET("/community", controller.CommunityHandler)
			auth.GET("/community/:id", controller.CommunityDetailHandler)
			auth.POST("/post", controller.CreatePostHandler)
			auth.GET("/post/:id", controller.GetPostHandler)
			auth.GET("/posts", controller.GetPostListHandler)
			auth.GET("/posts2", controller.GetPostListHandler2)
			// 投票
			auth.POST("/vote", controller.PostVoteController)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}

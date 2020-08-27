package middleware

import (
	"bluebell/pkg/jwt"
	"bluebell/serializer"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	CtxtUserIDKey = "userID"
)
// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式：1.放在请求头2.放在请求体3.放在URI
		// 这里假设Token都放在请求头中的Authorization,并使用Bearer开头
		// 这里的具体实现方式根据实际的业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			serializer.ResponseError(c, serializer.CodeNeedLogin)
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			serializer.ResponseError(c, serializer.CodeInvalidToken)
			c.Abort()
			return
		}

		// parts[1]中是获取到的token，用定义好的ParseToken来解析获取到的token
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			serializer.ResponseError(c, serializer.CodeInvalidToken)
			c.Abort()
			return
		}

		// 将当前请求的userID信息保存到请求的上下文中
		c.Set(CtxtUserIDKey, mc.UserId)
		c.Next() // 后续的处理请求函数中可以通过c.Get("userID")来获取当前请求的用户信息
	}
}
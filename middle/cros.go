package middle

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors1() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		//http://localhost:8080,http://localhost:5500
		c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		//允许跨域设置可以返回其他子段，可以自定义字段
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token, Token")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")

		//允许客户端传递校验信息比如 cookie (重要)
		c.Header("Access-Control-Allow-Credentials", "true")

		//c.Header("Access-Control-Max-Age", "172800")
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
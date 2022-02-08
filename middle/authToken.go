package middle

import (
	"chatApp_backend/dao"
	"chatApp_backend/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type myClaims struct {
	Username string `json:"username"`
	UserID   string `json:"id"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2 //定义token的存活时间为2小时
//接下来还需要定义Secret：
var mySecret = []byte("夏天夏天悄悄过去") //加密，根据他在服务器端才能签发token

//token
func CreateToken(username string, id string) (string, error) {

	c := myClaims{
		username,
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), //过期时间
			Issuer:    "my-project",                               //签发人
		},
	}
	//使用指定的签名方法创造签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	//使用指定的secret签名并获得完整的编码后的token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*myClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &myClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*myClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// JWTAuthMiddleware 基于JWT的认证中间件
// 可以多检验一层userid
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		var user model.User
		//根据解析出来的id查询数据库看是否存在该用户
		err = dao.DB.First(&user, "userid=?", mc.UserID).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未注册！"})
			c.Abort()
			return
		}
		// 数据库存在该用户
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", user.Username)
		c.Set("userID", user.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

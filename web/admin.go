package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	DefaultPORT int = 8080
)

// PostSignIn 接收参数，往数据库里加一个User，返回一个token
func PostSignIn() func(c *gin.Context){
	return func(c *gin.Context) {
		// TODO: 参数校验
		var user User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		id := insertUser(&user)
		if id == -1 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed insert user..."})
		}
		id = user.Id
		privilege := user.Privilege
		token,msg := GenToken(id,privilege)
		if msg != nil {
			msgStr := "Generate token failed: " + msg.Error()
			c.JSON(http.StatusInternalServerError,gin.H{"error":msgStr})
		}
		c.JSON(http.StatusOK, gin.H{"token":token})
	}
}

// PostSignUp 接收参数，验证数据库中是否有user, 如果有那么返回用户的权限
func PostSignUp() func(c *gin.Context){
	return func(c *gin.Context) {
		var user User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		realUser := checkUser(user)
		if realUser == (User{}) {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"retry your username or password ..."})
			return
		}
		token,msg := GenToken(realUser.Id,realUser.Privilege)
		if msg != nil {
			msgStr := "Generate token failed: " + msg.Error()
			c.JSON(http.StatusInternalServerError,gin.H{"error":msgStr})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token":token})
	}
}
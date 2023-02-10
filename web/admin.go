package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	DefaultPORT int = 8080
)

// postRegister 接收参数，往数据库里加一个User
func postRegister(c *gin.Context) {
	// TODO: 参数校验
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, user)
}

// postSignIn 接收参数，验证数据库中是否有user, 如果有那么返回用户的权限
func postSignIn(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	b, msg, privilege := checkUser(user)
	if !b {
		c.JSON(http.StatusOK, msg)
	}
	c.JSON(http.StatusOK, privilege)
}

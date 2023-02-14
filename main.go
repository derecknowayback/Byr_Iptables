package main

import (
	"BYR_Iptables/web"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	web.InitDB()

	// 所有请求都经过 ipAuth 中间件
	// router.Use(web.IpAuthMiddleware())

	// admin 路由组
	adminRouter := router.Group("/admin")
	{
		adminRouter.POST("/signup",web.PostSignUp()) // 登录
		adminRouter.POST("/signin", web.PostSignIn()) // 注册
	}

	// iptables 路由组
	iptablesRouter := router.Group("/iptables")
	{
		iptablesRouter.Use(web.JWTAuthMiddleware())
		iptablesRouter.POST("/addRule",web.PostAddRule())
	}

	router.Run(":8080")
}

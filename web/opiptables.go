package web

import (
	"BYR_Iptables/iptables"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	DNAT = "dnat"
	SNAT = "snat"
)

// PostAddRule 2个参数 ip1 和 ip2, ip1
func PostAddRule() func(c *gin.Context){
	return func(c *gin.Context) {
		ip1, b1 := c.GetPostForm("ip1")
		ip2, b2 := c.GetPostForm("ip2")
		if !(b1 && b2) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ip format ..."})
			return
		}
		errDnat := iptables.Dnat(ip1,ip2)
		errSnat := iptables.Snat(ip2,ip1)
		if errDnat != nil || errSnat != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "insert failed...",
				"errDnat":errDnat.Error(),
				"errSnat":errSnat.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{"ok":"add ok"})
	}
}

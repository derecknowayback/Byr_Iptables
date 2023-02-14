package web

import (
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ip format ..."})
		}

		//errDnat := iptables.Dnat(ip1,ip2)
		//errSnat := iptables.Snat(ip2,ip1)
		errDnat := dummyNat(ip1,ip2)
		errSnat := dummyNat(ip1,ip2)
		if errDnat != nil || errSnat != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "insert failed...",
				"errDnat":errDnat.Error(),
				"errSnat":errSnat.Error(),
			})
		}
		c.JSON(http.StatusOK,gin.H{"ok":"add ok"})
	}
}

// use for test without iptables cmd
func dummyNat(ip1,ip2 string) error {
	return nil
}

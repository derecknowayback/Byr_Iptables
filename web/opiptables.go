package web

import (
	"BYR_Iptables/iptables"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

const (
	DNAT = "dnat"
	SNAT = "snat"
)

// addDnat 3个参数 ip1 和 ip2 以及 规则类型 type: dnat/snat
func add(c *gin.Context) {
	// TODO: 鉴权
	ip1Str, b1 := c.GetPostForm("ip1")
	ip2Str, b2 := c.GetPostForm("ip2")
	nat, _ := c.GetPostForm("type")
	if !(b1 && b2) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ip format ..."})
	}
	ip1 := net.ParseIP(ip1Str).To4()
	ip2 := net.ParseIP(ip2Str).To4()
	if ip1 == nil || ip2 == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Don't support Ipv6 now ..."})
	}
	var err error
	switch nat {
	case DNAT:
		err = iptables.Dnat(ip1, ip2)
	case SNAT:
		err = iptables.Snat(ip1, ip2)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "insert failed: " + err.Error()})
	}
}

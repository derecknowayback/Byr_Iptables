package web

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"time"
)

const (
	DefaultPORT int = 8080
)

// postRegister 接收参数，往数据库里加一个User
func postRegister(c *gin.Context) {
	// TODO: 参数校验
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, user)
}

// postSignIn 接收参数，验证数据库中是否有user, 如果有那么返回用户的权限
func postSignIn(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	b, msg, privilege := checkUser(user)
	if !b {
		c.JSON(http.StatusOK, msg)
	}
	c.JSON(http.StatusOK, privilege)
}

// 鉴权: IP鉴权 + JWT 鉴权

// IpAuthMiddleware 验证ip的地址
func IpAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ipStr := ctx.RemoteIP()
		if ipStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Couldn't get remote ip ..."})
		}
		ip4 := net.ParseIP(ipStr).To4()
		if ip4 == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Don't support Ipv6 now ..."})
		}
		if !checkIp(ip4) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Illegal ip address ..."})
		}
		ctx.Next()
	}
}

// 测试 ip在不在某个网段内
func checkIp(ip net.IP) bool {
	// TODO: 需要知道判断哪些ip能够放行
	return false
}

// MyClaims JWT 鉴权
type MyClaims struct {
	Userid             int           `json:"userid"`
	Privilege          UserPrivilege `json:"privilege"`
	jwt.StandardClaims               // jwt.StandardClaims包含了官方定义的字段
}

var MySecret = []byte("byriptable")

const TokenExpireDuration = time.Hour * 2

func GenToken(id int, privilege UserPrivilege) (string, error) {
	c := MyClaims{
		Userid:    id,
		Privilege: privilege,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),                          //签发时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(MySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 从token中提取我们自己的MyClaims, 注意这里token.Valid的检查似乎不包括过期检查
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix() //设定增加的时间
		return GenToken(claims.Userid, claims.Privilege)
	}
	return "", errors.New("couldn't handle this token")
}

func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		tokenStr := ctx.Request.Header.Get("Authorization") // 这里假设Token放在Header的Authorization中
		// 如果没有token的话
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Required token missed ...",
			})
			return
		}
		claims, err := ParseToken(tokenStr)
		// 如果token无效
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is invalid ...",
			})
			return
		} else if time.Now().Unix() > claims.ExpiresAt {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is expired ...", //token超市
			})
			return
		}
		ctx.Set("id", claims.Userid)
		ctx.Set("privilege", claims.Privilege)
		ctx.Next()
	}
}

package httputil

import (
	"time"

	"github.com/bbdshow/qelog/infra/jwt"
	"github.com/gin-gonic/gin"
)

type ClaimsData struct {
	UID      int64  `json:"uid"`
	Role     string `json:"role"`
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
}

var ClaimsKey = "jwt_claims_key"

// GenerateJWTToken
func GenerateJWTToken(data ClaimsData, ttl time.Duration, signingKey ...string) (token string, err error) {
	return jwt.GenerateJWTToken(jwt.NewCustomClaims(data, ttl, signingKey...))
}

func SetJWTClaims(c *gin.Context, token string, signingKey ...string) error {
	data := ClaimsData{}
	if err := jwt.GetCustomData(token, &data, signingKey...); err != nil {
		return err
	}
	c.Set(ClaimsKey, data)
	return nil
}

func GetJWTClaims(c *gin.Context) (data ClaimsData, exists bool) {
	val, ok := c.Get(ClaimsKey)
	if !ok {
		return data, ok
	}
	data, ok = val.(ClaimsData)
	if !ok {
		return data, ok
	}
	return data, true
}

// CheckClaimsUidIsEqual
func CheckClaimsUidIsEqual(c *gin.Context, uid int64) bool {
	data, exists := GetJWTClaims(c)
	if !exists {
		return exists
	}
	return data.UID == uid
}

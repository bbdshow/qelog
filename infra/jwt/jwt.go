package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	Issuer                          = "Anonymous"
	SigningKey                      = []byte("please replace signingKey, use SetSigningKey function.")
	SigningMethod jwt.SigningMethod = jwt.SigningMethodHS256

	ErrCustomClaimsInValid = errors.New("custom claims invalid")
)

func SetIssuer(issuer string) {
	Issuer = issuer
}

func SetSigningMethod(method jwt.SigningMethod) {
	SigningMethod = method
}

type CustomClaims struct {
	SigningKey []byte
	CustomData interface{} `json:"custom_data"`
	jwt.StandardClaims
}

func NewCustomClaims(data interface{}, ttl time.Duration, signingKey ...string) *CustomClaims {
	cc := &CustomClaims{
		SigningKey: SigningKey,
		CustomData: data,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Add(-1 * time.Second).Unix(),
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Issuer:    Issuer,
		},
	}
	if len(signingKey) > 0 {
		cc.SigningKey = []byte(signingKey[0])
	}
	return cc
}

func GenerateJWTToken(customClaims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(SigningMethod, customClaims)
	if customClaims.SigningKey == nil {
		customClaims.SigningKey = SigningKey
	}
	str, err := token.SignedString(customClaims.SigningKey)
	if err != nil {
		return "", err
	}
	return str, nil
}

func VerifyJWTToken(tokenStr string, signingKey ...string) (bool, error) {
	key := SigningKey
	if len(signingKey) > 0 {
		key = []byte(signingKey[0])
	}
	token, err := parseJWTToken(tokenStr, key)
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

func GetCustomData(tokenStr string, data interface{}, signingKey ...string) error {
	key := SigningKey
	if len(signingKey) > 0 {
		key = []byte(signingKey[0])
	}
	token, err := parseJWTToken(tokenStr, key)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return ErrCustomClaimsInValid
	}
	switch claims.CustomData.(type) {
	case map[string]interface{}:
		// 因为加密之前使用 JSON 编码，所以编码处理一下，返回到结构体
		byt, err := json.Marshal(claims.CustomData)
		if err != nil {
			return err
		}
		err = json.Unmarshal(byt, data)
		return err
	default:
		return nil
	}
}

func parseJWTToken(tokenStr string, signingKey []byte) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != SigningMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

package jwt

import (
	"crypto/rand"
	"errors"
	"es-3d-editor-go-back/struct"
	"fmt"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
	"io"
	"log"
	"time"
)

// JWT : HEADER 有效载荷签名
const (
	SecretKEY            string = "JWT-Secret-Key"
	DefaultExpireSeconds int    = 120 * 60 // 默认过期时长(s)
	PasswordHashBytes           = 16
)

// MyCustomClaims 有效载荷
type MyCustomClaims struct {
	UserID int `json:"UserID"`
	jwt.StandardClaims
}

// JwtPayload 令牌有效负载的解析
type JwtPayload struct {
	Username  string `json:"Username"`
	UserID    int    `json:"UserID"`
	IssuedAt  int64  `json:"Iat"`
	ExpiresAt int64  `json:"Exp"`
}

// GenerateToken "生成token"
// @param loginInfo *LoginRequest "登录信息"
// @param userID int
// @param expiredSeconds int "过期时间（秒）"
// @return tokenString string "token字符串"
func GenerateToken(loginInfo *_struct.LoginRequest, userID int, expiredSeconds int) (tokenString string, err error) {
	if expiredSeconds == 0 {
		expiredSeconds = DefaultExpireSeconds
	}

	// Create the Claims
	mySigningKey := []byte(SecretKEY)
	expireAt := time.Now().Add(time.Second * time.Duration(expiredSeconds)).Unix()
	logs.Info("Token will be expired at ", time.Unix(expireAt, 0))

	user := *loginInfo
	claims := MyCustomClaims{
		userID,
		jwt.StandardClaims{
			Issuer:    user.UserName,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireAt,
		},
	}

	// Create the token using your claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用secret对令牌进行签名
	tokenStr, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("error: failed to generate token")
	}

	return tokenStr, nil
}

// ValidateToken "验证token"
// @params tokenString string "token字符串"
// @return *JwtPayload "令牌有效负载的解析"
func ValidateToken(tokenString string) (*JwtPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKEY), nil
		})

	claims, ok := token.Claims.(*MyCustomClaims)
	if ok && token.Valid {
		log.Println("ok && token valid")
		logs.Info("%v %v", claims.UserID, claims.StandardClaims.ExpiresAt)
		logs.Info("Token was issued at ", time.Now().Unix())
		logs.Info("Token will be expired at ", time.Unix(claims.StandardClaims.ExpiresAt, 0))

		return &JwtPayload{
			Username:  claims.StandardClaims.Issuer,
			UserID:    claims.UserID,
			IssuedAt:  claims.StandardClaims.IssuedAt,
			ExpiresAt: claims.StandardClaims.ExpiresAt,
		}, nil
	} else {
		fmt.Println(err)
		return nil, errors.New("error: 无法验证token")
	}
}

// RefreshToken "更新token"
// @param tokenString string "token字符串"
// @return newTokenString string "新的token字符串"
func RefreshToken(tokenString string) (newTokenString string, err error) {
	// 获取上一个token
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKEY), nil
		})

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid {
		return "", err
	}

	mySigningKey := []byte(SecretKEY)
	expireAt := time.Now().Add(time.Second * time.Duration(DefaultExpireSeconds)).Unix() //new expired
	newClaims := MyCustomClaims{
		claims.UserID,
		jwt.StandardClaims{
			Issuer:    claims.StandardClaims.Issuer, //name of token issue
			IssuedAt:  time.Now().Unix(),            //time of token issue
			ExpiresAt: expireAt,                     // new expired
		},
	}

	// generate new token with new claims
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	// sign the token with a secret
	tokenStr, err := newToken.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("error: failed to generate new fresh json web token")
	}

	return tokenStr, nil
}

// GenerateSalt "生成SALT值"
// @return salt string
func GenerateSalt() (salt string, err error) {
	buf := make([]byte, PasswordHashBytes)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", errors.New("error: failed to generate user's salt")
	}

	return fmt.Sprintf("%x", buf), nil
}

// GeneratePassHash 生成Hash密码
// @param password string 明文密码
// @param salt string SALT值
// @return hash string Hash密码
func GeneratePassHash(password string, salt string) (hash string, err error) {
	h, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, PasswordHashBytes)
	if err != nil {
		return "", errors.New("error: failed to generate password hash")
	}

	return fmt.Sprintf("%x", h), nil
}

// CheckStatus "验证token是否过期"
// @param tokenString string "token字符串"
func CheckStatus(tokenString string) (string, int64) {
	jp, err := ValidateToken(tokenString)

	if err != nil {
		// 如果令牌已经过期
		fmt.Println("你的令牌已过期, 请重新登录! ")
		return "", -1
	}

	timeDiff := jp.ExpiresAt - time.Now().Unix()
	fmt.Println("timeDiff = ", timeDiff)
	if timeDiff <= 30 {
		// 如果令牌即将到期，请刷新令牌
		fmt.Println("您的令牌很快就要过期了")
		newToken, err := RefreshToken(tokenString)
		if err == nil {
			return newToken, timeDiff
		}
	}
	// 如果令牌有效，则不执行任何操作
	fmt.Println("Your token is good ")
	return tokenString, timeDiff
}

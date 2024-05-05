package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cxi7448/cxhttp/clGlobal"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/core/clAuth"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	TokenExpireDuration  = int64(4000) // 有效期
	TokenReflushDuration = int64(3600) // 刷新时间
)

var SecretKey = []byte("32honefzr7vnbm0k")

const JWT_PREFIX = "JWT_U_INFO_"

type UserInfo struct {
	Uid       uint64            `json:"uid,omitempty"`   // 当前用户Id
	Token     string            `json:"token,omitempty"` // redis登陆校验的token
	ExtraData map[string]string `json:"ed,omitempty"`    // 附属数据
}
type Claims struct {
	UserInfo *UserInfo `json:"uinfo,omitempty"`
	jwt.StandardClaims
	CreateTime  int64 `json:"ct,omitempty"` // 创建时间
	ReflushTime int64 `json:"st,omitempty"` // 刷新时间
}

func (this Claims) SaveToRedis() {
	redis := clGlobal.GetRedis()
	err := redis.Set(fmt.Sprintf("%v%v", JWT_PREFIX, this.UserInfo.Uid), this.UserInfo, int32(TokenExpireDuration))
	if err != nil {
		clLog.Error("存入redis失败:%v", err)
	}
}

// 从redis中校验，redis不存在的时候，表示有效，存在的时候比对签名，签名一样的有效，不一样的无效
func (this Claims) IsEffective() bool {
	redis := clGlobal.GetRedis()
	uStr := redis.Get(fmt.Sprintf("%v%v", JWT_PREFIX, this.UserInfo.Uid))
	if uStr != "" {
		uInfo := UserInfo{}
		json.Unmarshal([]byte(uStr), &uInfo)
		if uInfo.Token != this.UserInfo.Token {
			// 被其他人登陆占用了
			return false
		}
	}
	return true
}

func (this Claims) IsExpire() bool {
	return this.ExpiresAt > time.Now().Unix()
}

func (this Claims) IsReflush() bool {
	return this.ReflushTime < time.Now().Unix()
}
func (this Claims) GetUser() *clAuth.AuthInfo {
	user := clAuth.NewUser(this.UserInfo.Uid, this.UserInfo.Token)
	user.ExtraData = this.UserInfo.ExtraData
	return user
}

func (this Claims) ReflushToken() (string, error) {
	return GenToken(this.GetUser())
}

func GenToken(user *clAuth.AuthInfo) (string, error) {
	now := time.Now().Unix()
	c := Claims{
		CreateTime:  now,
		ReflushTime: now + TokenReflushDuration,
		UserInfo: &UserInfo{
			Uid:       user.Uid,
			Token:     user.Token,
			ExtraData: user.ExtraData,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now + TokenExpireDuration,
			Issuer:    "goapi",
			Subject:   "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	result, err := token.SignedString(SecretKey)
	if err != nil {
		clLog.Error("生成jwt失败:%v", err)
		return "", err
	}
	c.SaveToRedis()
	return result, nil
}

// 解析token
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return SecretKey, nil
	})
	if err != nil {
		clLog.Error("错误了:%v", err)
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

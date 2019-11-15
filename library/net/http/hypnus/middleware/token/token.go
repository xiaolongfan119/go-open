package token

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	xtime "github.com/ihornet/go-open/library/time"

	redis "github.com/ihornet/go-open/library/database/redis"
	hp "github.com/ihornet/go-open/library/net/http/hypnus"

	"github.com/ihornet/go-open/library/ecode"

	jwt "github.com/dgrijalva/jwt-go"
)

var Instance Token

var (
	now   int64
	mu    sync.Mutex
	index int64
)

type Token struct {
	Conf *TokenConfig
}

func Init(conf *TokenConfig) {
	Instance = Token{Conf: conf}
}

type TokenConfig struct {
	Secret     string
	Expiration xtime.Duration
}

func (t *Token) Verify(ctx *hp.Context) {

	stoken := ctx.Request.Header.Get("token")
	if stoken == "" {
		t.handleFailed(ctx, ecode.TokenEmpty)
		return
	}
	token, err := jwt.Parse(stoken, t.getValidationKey)
	if err != nil {
		t.handleFailed(ctx, ecode.TokenInvalid)
		return
	}

	if !token.Valid {
		t.handleFailed(ctx, ecode.TokenInvalid)
		return
	}
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return
	}

	ctx.Req.Header = make(map[string]string)
	userId := mapClaims["user"].(map[string]interface{})["id"]
	param1 := mapClaims["user"].(map[string]interface{})["param1"]
	param2 := mapClaims["user"].(map[string]interface{})["param2"]

	ctx.Req.Header["userId"] = fmt.Sprintf("%v", userId)
	ctx.Req.Header["wand"] = mapClaims["wand"].(string)

	if param1 != nil {
		ctx.Req.Header["param1"] = fmt.Sprintf("%v", param1)
	}
	if param2 != nil {
		ctx.Req.Header["param2"] = fmt.Sprintf("%v", param2)
	}
}

func (t *Token) GenToken(payload interface{}) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Duration(t.Conf.Expiration)).Unix()
	claims["user"] = payload
	claims["wand"] = getToken()
	token.Claims = claims

	// 把token已约定的加密方式和加密秘钥加密，当然也可以使用不对称加密
	tokenString, _ := token.SignedString([]byte(t.Conf.Secret))

	if redis.RedisClient != nil {
		go t.AddToken(payload, claims["wand"].(string))
	}

	return tokenString
}

// 在redis查找token, 如果不存在就错误
func (t *Token) VerifyRedis(ctx *hp.Context) {

	if redis.RedisClient == nil {
		return
	}

	wand := ctx.Req.Header["wand"]
	userId := ctx.Req.Header["userId"]
	key := fmt.Sprintf("%s:user:%s:token", hp.ServerName, userId)

	v, _ := redis.RedisClient.Get(key).Result()

	if v != wand {
		t.handleFailed(ctx, ecode.TokenInvalid2)
	}
}

// 将token加入redis黑名单
func (t *Token) DisableToken(userId int) {

	if redis.RedisClient == nil {
		return
	}

	key := fmt.Sprintf("%s:user:%s:token", hp.ServerName, userId)

	redis.RedisClient.Del(key)

	//redis.RedisClient.Set(key, "-", time.Duration(t.Conf.Expiration))
}

func (t *Token) AddToken(payload interface{}, wand string) {

	_payload := reflect.ValueOf(payload)
	isValid := _payload.FieldByName("Id").IsValid()
	if isValid {
		userId := _payload.FieldByName("Id").Interface()
		key := fmt.Sprintf("%s:user:%d:token", hp.ServerName, userId.(int))
		redis.RedisClient.Set(key, wand, time.Duration(t.Conf.Expiration))
	}

}

func (t *Token) handleFailed(ctx *hp.Context, err ecode.Code) {
	ctx.Abort()
	ctx.JSON(nil, err)
}

func (t *Token) getValidationKey(*jwt.Token) (interface{}, error) {
	return []byte(t.Conf.Secret), nil
}

func getToken() string {
	if now == 0 {
		now = time.Now().UnixNano() / 1e6
	}
	atomic.AddInt64(&index, 1)
	return fmt.Sprintf("%d-%d", now, index)
}

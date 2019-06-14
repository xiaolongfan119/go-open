package token

import (
	"go-open/library/ecode"
	hp "go-open/library/net/http/hypnus"
	xtime "go-open/library/time"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var Instance Token

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
}

func (t *Token) GenToken(payload interface{}) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Duration(t.Conf.Expiration)).Unix()
	claims["user"] = payload
	token.Claims = claims

	// 把token已约定的加密方式和加密秘钥加密，当然也可以使用不对称加密
	tokenString, _ := token.SignedString([]byte(t.Conf.Secret))
	return tokenString
}

func (t *Token) handleFailed(ctx *hp.Context, err ecode.Code) {
	ctx.Abort()
	ctx.JSON(nil, err)
}

func (t *Token) getValidationKey(*jwt.Token) (interface{}, error) {
	return []byte(t.Conf.Secret), nil
}

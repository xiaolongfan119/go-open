package token

import (
	"net/http"
	"testing"
	"time"

	hp "github.com/ihornet/go-open/v2/library/net/http/hypnus"
	xtime "github.com/ihornet/go-open/v2/library/time"
)

func TestToken(t *testing.T) {

	var context hp.Context

	token := Token{
		Conf: &TokenConfig{
			Secret:     "as13adsa",
			Expiration: xtime.Duration(time.Second * 10),
		},
	}

	user := User{Id: 10, Name: "nma"}
	stoken := token.GenToken(user)

	req := http.Request{}
	req.Header = map[string][]string{"token": []string{stoken}}
	context.Request = &req
	token.Verify(&context)

	println(context.Req.Header["userId"])
	println(stoken)
}

type User struct {
	Id   int
	Name string
}

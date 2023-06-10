package app

import (
	"context"
	"time"

	"github.com/glennliao/apijson-go/config/executor"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	executor.RegActionExecutor("custom", &ActionExecutor{})
}

type ActionExecutor struct{}

func (a *ActionExecutor) Do(ctx context.Context, req executor.ActionExecutorReq) (ret model.Map, err error) {

	body := req.Data[0]

	user, err := g.DB().Model("user").One(g.Map{
		"username": body["email"].(string),
	})

	if err != nil {
		return nil, err
	}

	password := body["password"].(string)

	if user == nil || user.IsEmpty() || bcrypt.CompareHashAndPassword([]byte(gconv.String(user.Map()["password"])), []byte(password)) != nil {
		time.Sleep(time.Millisecond * time.Duration(grand.N(100, 500)))
		return nil, gerror.New("账户或密码错误")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo":    "bar",
		"userId": user.Map()["user_id"],
		"nbf":    time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}
	return model.Map{
		"token": tokenString,
	}, nil
}

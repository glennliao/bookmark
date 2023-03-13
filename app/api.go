package app

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type api struct{}

var Api = api{}

type (
	AuthReq struct {
		g.Meta   `path:"/auth" method:"POST"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	AuthRes struct {
		Token string `json:"token"`
	}
)

type (
	AddUseReq struct {
		g.Meta `path:"/addUse" method:"POST"`
		BmId   string `json:"bmId"`
	}
	AddUseRes struct {
	}
)

func (a *api) Auth(ctx context.Context, req *AuthReq) (res *AuthRes, err error) {

	user, err := g.DB().Model("user").One(g.Map{
		"username": req.Email,
	})
	if err != nil {
		return nil, err
	}

	if user == nil || user.IsEmpty() || bcrypt.CompareHashAndPassword([]byte(gconv.String(user.Map()["password"])), []byte(req.Password)) != nil {
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

	return &AuthRes{Token: tokenString}, nil
}

func (a *api) AddUse(ctx context.Context, req *AddUseReq) (res *AddUseRes, err error) {

	m := g.DB().Model("bookmark_use").Safe().Ctx(ctx)

	user, _ := ctx.Value(UserIdKey).(*CurrentUser)

	num, err := m.UpdateAndGetAffected(g.Map{
		"times": &gdb.Counter{
			Field: "times",
			Value: gconv.Float64(1),
		},
	}, g.Map{
		"bm_id":   req.BmId,
		"user_id": user.UserId,
	})
	if err != nil {
		return nil, err
	}

	if num == 0 {
		m.Insert(g.Map{
			"times":   "1",
			"bm_id":   req.BmId,
			"user_id": user.UserId,
		})
	}

	return &AddUseRes{}, nil
}

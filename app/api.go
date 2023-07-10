package app

import (
	"context"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yitter/idgenerator-go/idgen"
	"golang.org/x/crypto/bcrypt"
)

func readConfig(ctx context.Context, key string) *gvar.Var {
	cfg, err := g.Cfg().Get(ctx, key)
	if err != nil {
		panic(err)
	}
	return cfg
}

type Api struct {
	A *apijson.ApiJson
}

// register
type (
	RegisterReq struct {
		g.Meta   `method:"POST" path:"/register"`
		Email    string `v:"required"`
		Password string `v:"required"`
		Code     string
	}

	RegisterRes struct {
		Msg string `json:"msg"`
	}

	RegisterLayoutReq struct {
		g.Meta `method:"GET" path:"/register"`
	}

	RegisterLayoutRes struct {
		CanReg bool   `json:"canReg"`
		Msg    string `json:"msg"`
	}
)

func (a *Api) Register(ctx context.Context, req *RegisterReq) (res *RegisterRes, err error) {

	canReg := readConfig(ctx, "app.canReg").Bool()
	if !canReg {
		return nil, gerror.New("注册已关闭")
	}

	// todo 防止暴力注册

	regCode := readConfig(ctx, "app.regCode").String()
	if regCode != req.Code {
		return nil, gerror.New("regCode 不正确")
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	password := string(hash)

	one, err := g.DB().Model("user").Where("username", req.Email).One()
	if err != nil {
		return nil, err
	}

	if !one.IsEmpty() {
		return nil, gerror.New("用户已存在")
	}

	userId := strconv.FormatInt(idgen.NextId(), 10)

	_, err = g.DB().Model("user").Insert(g.Map{
		"user_id":  userId,
		"username": req.Email,
		"password": password,
	})

	if err != nil {
		return nil, err
	}

	act := a.A.NewAction(ctx, http.MethodPost, model.Map{
		"tag": "Groups",
		"Groups": model.Map{
			"groupId": userId,
			"title":   "个人分组",
		},
		"GroupUser[]": []model.Map{
			{
				"groupId": userId,
				"userId":  userId,
			},
		},
	})
	act.NoAccessVerify = true
	_, err = act.Result()

	return &RegisterRes{Msg: ""}, err
}

func (a *Api) RegisterLayout(ctx context.Context, req *RegisterLayoutReq) (res *RegisterLayoutRes, err error) {

	canReg := readConfig(ctx, "app.canReg").Bool()
	msg := ""
	one, err := g.DB().Model("user").One()
	if err != nil {
		return nil, err
	}

	if one.IsEmpty() {
		msg = "第一个用户注册后即为管理员"
		canReg = true
	}
	return &RegisterLayoutRes{CanReg: canReg, Msg: msg}, nil
}

// auth
type (
	AuthReq struct {
		g.Meta   `method:"POST" path:"/auth"`
		Email    string `v:"required"`
		Password string `v:"required"`
	}
	AuthRes struct {
		Token string `json:"token"`
	}
)

func (a *Api) Auth(ctx context.Context, req *AuthReq) (res *AuthRes, err error) {

	user, err := g.DB().Model("user").One(g.Map{
		"username": strings.TrimSpace(req.Email),
	})

	if err != nil {
		return nil, err
	}

	password := strings.TrimSpace(req.Password)

	if user == nil || user.IsEmpty() || bcrypt.CompareHashAndPassword([]byte(gconv.String(user.Map()["password"])), []byte(password)) != nil {
		time.Sleep(time.Millisecond * time.Duration(grand.N(100, 500)))
		return nil, gerror.New("账户或密码错误")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Map()["user_id"],
		"iat":    time.Now().Unix(),
	})

	// Sign and get the complete encoded token as A string using the secret
	tokenString, err := token.SignedString(jwtSecret)

	return &AuthRes{Token: tokenString}, err
}

// upload
type (
	UploadReq struct {
		g.Meta `method:"POST" path:"/upload"`
		File   *ghttp.UploadFile
	}

	UploadRes struct {
		Path string `json:"path"`
	}

	IconReq struct {
		g.Meta `method:"GET" path:"/icon"`
		Name   string `json:"name"`
	}

	IconRes struct {
	}
)

func (a *Api) Upload(ctx context.Context, req *UploadReq) (res *UploadRes, err error) {

	iconSavePath, _ := filepath.Abs(readConfig(ctx, "app.iconSavePath").String())

	path, err := req.File.Save(iconSavePath, true)

	return &UploadRes{Path: path}, err
}

func (a *Api) Icon(ctx context.Context, req *IconReq) (res *IconRes, err error) {
	r := g.RequestFromCtx(ctx)
	iconSavePath, _ := filepath.Abs(readConfig(ctx, "app.iconSavePath").String())

	path := filepath.Join(iconSavePath, req.Name)
	if !strings.HasPrefix(path, iconSavePath) {

		r.Response.Status = 500
		r.Response.Write("?")
		return
	}
	r.Response.ServeFile(path, false)
	return &IconRes{}, nil
}

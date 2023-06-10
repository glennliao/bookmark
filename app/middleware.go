package app

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitToken() {
	val, err := g.Cfg().Get(nil, "jwt.secret")
	if err != nil {
		g.Log().Fatal(nil, err)
	}
	jwtSecret = val.Bytes()
	if len(jwtSecret) < 16 {
		g.Log().Fatal(nil, "jwt.secret不能小于16位")
	}
}

func Auth(r *ghttp.Request) {

	authorization := r.Request.Header.Get("Authorization")

	if authorization != "" {
		ctx := r.Context()
		token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return jwtSecret, nil
		})

		if err != nil {
			g.Log().Info(ctx, "token err", err)
			r.Response.WriteJson(g.Map{
				"code": 401,
				"msg":  "未知错误",
			})
			return
		}

		userId := ""
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId = gconv.String(claims["userId"])
			if userId == "" {
				g.Log().Info(ctx, "userId is empty")
				r.Response.WriteJson(g.Map{
					"code": 401,
					"msg":  "未知错误",
				})
				return
			}
		} else {
			g.Log().Info(ctx, "token no valid")
			r.Response.WriteJson(g.Map{
				"code": 401,
				"msg":  "未知错误",
			})
			return
		}
		ctx = context.WithValue(ctx, UserIdKey, &CurrentUser{UserId: userId})
		r.SetCtx(ctx)
	} else {
		if r.Request.URL.Path != "/api/data/auth" {
			r.Response.WriteJson(g.Map{
				"code": 401,
				"msg":  "未登录",
			})
			return
		}

	}

	r.Middleware.Next()
}

func Cors(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	corsOptions.AllowOrigin = r.Request.Header.Get("Origin")
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

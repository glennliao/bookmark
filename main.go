package main

import (
	"context"
	"github.com/glennliao/apijson-go"
	_ "github.com/glennliao/apijson-go/drivers/config/goframe"
	_ "github.com/glennliao/apijson-go/drivers/executor/goframe"
	"github.com/glennliao/apijson-go/drivers/framework_goframe"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/bookmark/app"
	_ "github.com/glennliao/bookmark/packed"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
)

func main() {

	gres.Dump()

	a := apijson.Load(app.Init)

	if len(os.Args) > 1 {
		if os.Args[1] == "adduser" {
			createUser(a)
		} else {
			g.Log().Fatal(nil, "unknown command")
		}
		return
	}

	f := framework_goframe.New(a)

	s := g.Server()

	s.Group("/api/", func(group *ghttp.RouterGroup) {
		group.Middleware(app.Cors, app.Auth, response)
		group.Bind(app.Api)
	})

	s.Group("/api/data", func(group *ghttp.RouterGroup) {
		group.Middleware(app.Cors, app.Auth)
		f.Bind(group)
	})

	s.AddStaticPath("/", "dist")
	s.AddStaticPath("/assets", "dist/assets")

	s.Run()
}

func createUser(a *apijson.ApiJson) {
	email := gcmd.Scan("What's email?\n")
	email = strings.TrimSpace(email)
	if email == "" {
		g.Log().Fatal(nil, "email不为空")
	}
	password := gcmd.Scan("What's password?\n")
	password = strings.TrimSpace(password)
	if email == "" {
		g.Log().Fatal(nil, "password不为空")
	}

	ctx := context.Background()

	pw, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	act := a.NewAction(ctx, http.MethodPost, model.Map{
		"User": g.Map{
			"username": email,
			"password": gconv.String(pw),
		},
		"tag": "User",
	})
	act.NoAccessVerify = true
	_, err := act.Result()

	if err != nil {
		g.Log().Fatal(ctx, err)
	}

}

func response(r *ghttp.Request) {
	r.Middleware.Next()
	res := r.GetHandlerResponse()
	err := r.GetError()
	if err != nil {
		r.Response.WriteJson(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	} else {
		r.Response.WriteJson(g.Map{
			"code": 200,
			"data": res,
		})
	}
}

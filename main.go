package main

import (
	"context"
	_ "embed"
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
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/grand"
	"net/http"
	"os"
	"strings"
)

//go:embed config.toml.example
var configExample string

func main() {

	initConfigFile()

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

	s.Group("/api/data", func(group *ghttp.RouterGroup) {
		group.Middleware(app.Cors, app.Auth)

		group.POST("/auth", f.CommonResponse(func(ctx context.Context, req model.Map) (res model.Map, err error) {
			return a.NewAction(ctx, http.MethodPost, req).Result()
		}, framework_goframe.InDataMode))

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

	act := a.NewAction(ctx, http.MethodPost, model.Map{
		"User": g.Map{
			"username": email,
			"password": password,
		},
		"tag": "User",
	})
	act.NoAccessVerify = true
	_, err := act.Result()

	if err != nil {
		g.Log().Fatal(ctx, err)
	}

}

func initConfigFile() {
	if !gfile.Exists("config.toml") {
		configExample = strings.Replace(configExample, "{{secret}}", grand.S(32), 1)
		gfile.PutContents("config.toml", configExample)
	}
}

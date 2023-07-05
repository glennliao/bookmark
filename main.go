package main

import (
	"context"
	_ "embed"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/glennliao/apijson-go"
	_ "github.com/glennliao/apijson-go/drivers/goframe"
	"github.com/glennliao/apijson-go/drivers/goframe/web"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/bookmark/app"
	_ "github.com/glennliao/bookmark/packed"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/grand"
	"golang.org/x/crypto/bcrypt"
)

//go:embed config.toml.example
var configExample string

var (
	Main = &gcmd.Command{
		Name:        "server",
		Brief:       "start http server",
		Description: "",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {

			a := apijson.Load(app.App)

			w := web.New(a)

			s := g.Server()

			s.Group("/api/data", func(group *ghttp.RouterGroup) {
				group.Middleware(app.Cors, app.Auth)
				group.POST("/auth", w.CommonResponse(func(ctx context.Context, req model.Map) (res model.Map, err error) {
					return a.NewAction(ctx, http.MethodPost, req).Result()
				}, web.InDataMode))
				w.Bind(group)
			})

			s.Group("/api", func(group *ghttp.RouterGroup) {
				iconSavePath, _ := filepath.Abs("./runtime/icon")

				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(app.Cors, app.Auth)
					group.POST("/upload", func(req *ghttp.Request) {
						file := req.GetUploadFile("file")
						filePath, err := file.Save(iconSavePath, true)
						if err != nil {
							req.Response.WriteJson(g.Map{
								"code": 500,
								"msg":  err.Error(),
							})
							return
						}

						req.Response.WriteJson(g.Map{
							"code": 200,
							"data": g.Map{
								"path": filePath,
							},
						})

					})
				})

				group.GET("/icon", func(req *ghttp.Request) {
					name := req.GetQuery("name").String()
					path := filepath.Join(iconSavePath, name)
					if !strings.HasPrefix(path, iconSavePath) {
						req.Response.Status = 500
						req.Response.Write("?")
						return
					}
					req.Response.ServeFile(path, false)
				})
			})

			if g.Res().Contains("dist") {
				s.AddStaticPath("/", "dist")
				s.AddStaticPath("/assets", "dist/assets")
			}

			s.Run()
			return nil
		},
	}
	Init = &gcmd.Command{
		Name:        "init",
		Brief:       "init config file",
		Description: "",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			if !gfile.Exists("config.toml") {
				configExample = strings.Replace(configExample, "{{secret}}", grand.S(32), 1)
				gfile.PutContents("config.toml", configExample)
			}
			return
		},
	}

	CreateUser = &gcmd.Command{
		Name:  "createUser",
		Brief: "create user",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			createUser(apijson.Load(app.App))
			return nil
		},
	}
)

func main() {
	err := Main.AddCommand(Init, CreateUser)
	if err != nil {
		panic(err)
	}
	Main.Run(gctx.New())
}

func createUser(a *apijson.ApiJson) {
	ctx := context.Background()

	email := gcmd.Scan("What's email?\n")
	email = strings.TrimSpace(email)
	if email == "" {
		g.Log().Fatal(ctx, "email不为空")
	}
	password := gcmd.Scan("What's password?\n")
	password = strings.TrimSpace(password)
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	password = string(hash)
	if email == "" {
		g.Log().Fatal(ctx, "password不为空")
	}

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

package main

import (
	"context"
	_ "embed"
	"github.com/samber/lo"
	"strings"

	"github.com/glennliao/apijson-go"
	_ "github.com/glennliao/apijson-go/drivers/goframe"
	"github.com/glennliao/apijson-go/drivers/goframe/web"
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
)

//go:embed config.yaml.example
var configExample string

var (
	Main = &gcmd.Command{
		Name:        "server",
		Brief:       "start http server",
		Description: "",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {

			// 数据库同步前判断
			tables, err := g.DB().Tables(ctx)
			if err != nil {
				g.Log().Fatal(ctx, err)
			}

			a := apijson.Load(app.App)

			w := web.New(a)

			s := g.Server()

			s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Middleware(app.Cors, app.Response, app.Auth)
				group.Bind(app.Api{
					A: a,
				})
			})

			s.Group("/api/data", func(group *ghttp.RouterGroup) {
				group.Middleware(app.Cors, app.Auth)
				w.Bind(group)
			})

			if g.Res().Contains("dist") {
				s.AddStaticPath("/", "dist")
				s.AddStaticPath("/assets", "dist/assets")
			}

			if !lo.Contains(tables, "config") && lo.Contains(tables, "group_bookmark") {
				app.ToV0Dot11(ctx)
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
			if !gfile.Exists("config.yaml") {
				configExample = strings.Replace(configExample, "{{secret}}", grand.S(32), 1)
				gfile.PutContents("config.yaml", configExample)
			}
			return
		},
	}
)

func main() {
	err := Main.AddCommand(Init)
	if err != nil {
		panic(err)
	}
	Main.Run(gctx.New())
}

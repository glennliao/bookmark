package app

import (
	"context"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
	"github.com/glennliao/bookmark/app/util/fetchurl"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

func initFunc(a *apijson.ApiJson) {
	a.Config().Functions.Bind("fetchURL", config.Func{
		ParamList: []config.ParamItem{
			{
				Name: "url",
				Type: "string",
				Desc: "url地址",
			},
		},
		Handler: func(ctx context.Context, param model.Map) (res any, err error) {
			url := gconv.String(param["url"])

			var meta *fetchurl.UrlMeta
			if !strings.HasPrefix(strings.ToLower(url), "http") {
				meta, err = fetchurl.FetchURLMeta(ctx, "https://"+url)
				if err != nil {
					meta, err = fetchurl.FetchURLMeta(ctx, "http://"+url)
				}
			} else {
				meta, err = fetchurl.FetchURLMeta(ctx, url)
			}

			return meta, err
		},
	})

	a.Config().Functions.Bind("addUse", config.Func{
		ParamList: []config.ParamItem{
			{
				Name: "bmId",
				Type: "string",
				Desc: "书签id",
			},
		},
		Handler: func(ctx context.Context, param model.Map) (res any, err error) {

			user, _ := ctx.Value(UserIdKey).(*CurrentUser)

			bmId := util.String(param["bmId"])

			m := g.DB().Model("bookmark_use").Safe().Ctx(ctx)

			num, err := m.UpdateAndGetAffected(g.Map{
				"times": &gdb.Counter{
					Field: "times",
					Value: gconv.Float64(1),
				},
			}, g.Map{
				"bm_id":   bmId,
				"user_id": user.UserId,
			})
			if err != nil {
				return nil, err
			}

			if num == 0 {
				m.Insert(g.Map{
					"times":   1,
					"bm_id":   bmId,
					"user_id": user.UserId,
				})
			}

			return "ok", err
		},
	})
}

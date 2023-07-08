package app

import (
	"context"
	url2 "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gcache"

	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/bookmark/app/util/fetchurl"
	"github.com/glennliao/bookmark/app/util/htmlbookmark"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/yitter/idgenerator-go/idgen"
)

func initFunc(a *apijson.ApiJson) {
	a.Config().Functions.Bind("fetchURL", fetchURL)
	a.Config().Functions.Bind("import", importBookmark)
	a.Config().Functions.Bind("cateIdByBmId", cateIdByBmId())
	a.Config().Functions.Bind("latestVersion", latestVersion)
}

var cateIdByBmId = func() config.Func {

	var cache = gcache.NewAdapterMemory()
	return config.Func{
		ParamList: []config.ParamItem{
			{
				Name: "bmId",
				Type: "string",
				Desc: "",
			},
			{
				Name: "groupId",
				Type: "string",
				Desc: "",
			},
		},
		Handler: func(ctx context.Context, param model.Map) (res any, err error) {

			groupId := gconv.String(param["groupId"])
			if groupId == "" {
				groupId = ctx.Value(UserIdKey).(*CurrentUser).UserId
			}
			val, err := cache.GetOrSetFuncLock(ctx, groupId, func(ctx context.Context) (value interface{}, err error) {
				values, err := g.DB().Model("group_bookmark").Where("group_id", groupId).Fields("cate_id,bm_id").All()
				m := map[string]string{}
				for _, v := range values {
					item := gconv.MapStrStr(v)
					m[item["bm_id"]] = item["cate_id"]
				}
				return m, err
			}, time.Second*5)

			bmId := gconv.String(param["bmId"])
			return val.MapStrStr()[bmId], err
		},
	}
}

var fetchURL = config.Func{
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

		if meta != nil {
			if strings.HasPrefix(meta.Icon, "/") {
				_url, _ := url2.Parse(meta.Url)
				meta.Icon = _url.Scheme + "://" + _url.Host + meta.Icon
			}

			meta.Url = url
		}

		return meta, err
	},
}

var importBookmark = config.Func{
	ParamList: []config.ParamItem{
		{
			Name: "data",
			Type: "string",
		},
	},
	Handler: func(ctx context.Context, param model.Map) (res any, err error) {

		data := param["data"].(string)

		bookmarks, err := htmlbookmark.Parse(ctx, data)
		if err != nil {
			return nil, err
		}

		user, _ := ctx.Value(UserIdKey).(*CurrentUser)

		err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {

			one, err := g.DB().Ctx(ctx).Model("bookmark_cate").Where("group_id", user.UserId).Where("parent_id", "root").Where("title", "导入").One()

			if err != nil {
				return err
			}

			parentId := ""
			if one.IsEmpty() {

				parentId = strconv.FormatInt(idgen.NextId(), 10)
				_, err := g.DB().Ctx(ctx).Model("bookmark_cate").Insert(g.Map{
					"cate_id":   parentId,
					"title":     "导入",
					"sorts":     5,
					"parent_id": "root",
					"group_id":  user.UserId,
				})

				if err != nil {
					return err
				}

			} else {
				parentId = one.Map()["cate_id"].(string)
			}

			return process(ctx, bookmarks, parentId)
		})

		if err != nil {
			return nil, err
		}

		return g.Map{}, nil
	},
}

func process(ctx context.Context, bookmarks *htmlbookmark.Bookmarks, parentId string) error {

	bmList := make([]model.Map, 0)
	bmGroupList := make([]model.Map, 0)

	user, _ := ctx.Value(UserIdKey).(*CurrentUser)

	for _, item := range bookmarks.Children {
		if item.IsDir {

			cateId := strconv.FormatInt(idgen.NextId(), 10)

			_, err := g.DB().Ctx(ctx).Model("bookmark_cate").Insert(g.Map{
				"cate_id":   cateId,
				"title":     item.Title,
				"sorts":     5,
				"parent_id": parentId,
				"group_id":  user.UserId,
			})

			if err != nil {
				return err
			}

			err = process(ctx, item, cateId)
			if err != nil {
				return err
			}
		} else {
			bmId := strconv.FormatInt(idgen.NextId(), 10)
			bmList = append(bmList, g.Map{
				"bm_id":       bmId,
				"title":       item.Title,
				"url":         item.URL,
				"description": "@import ",
			})

			bmGroupList = append(bmGroupList, g.Map{
				"bm_id":    bmId,
				"group_id": user.UserId,
				"cate_id":  parentId,
				"sorts":    5,
			})
		}
	}

	if len(bmList) > 0 {
		_, err := g.DB().Ctx(ctx).Model("bookmark").Insert(bmList)

		if err != nil {
			return err
		}

		_, err = g.DB().Ctx(ctx).Model("group_bookmark").Insert(bmGroupList)

		if err != nil {
			return err
		}

	}

	return nil
}

var latestVersion = config.Func{
	Handler: func(ctx context.Context, param model.Map) (res any, err error) {

		url := "https://api.github.com/repos/glennliao/bookmark/releases/latest"
		resp, err := gclient.New().Get(ctx, url)
		if err != nil {
			return nil, err
		}

		json := gjson.New(resp.ReadAllString())

		return json.Get("tag_name"), nil
	},
}

package app

import (
	"context"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/bookmark/app/util/fetchurl"
	"github.com/glennliao/bookmark/app/util/htmlbookmark"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/panjf2000/ants/v2"
	"github.com/yitter/idgenerator-go/idgen"
	url2 "net/url"
	"strconv"
	"strings"
)

func initFunc(a *apijson.ApiJson) {
	a.Config().Functions.Bind("fetchURL", fetchURL)
	a.Config().Functions.Bind("import", importBookmark)

	a.Config().Functions.Bind("latestVersion", latestVersion())
	a.Config().Functions.Bind("noteTags", noteTags())
	a.Config().Functions.Bind("bmTags", bmTags())
	a.Config().Functions.Bind("cateBmNum", cateBmNum())
	a.Config().Functions.Bind("fetchMetaBatch", fetchMetaBatchFunc())
}

var fetchURL = config.Func{
	ParamList: []config.ParamItem{
		{
			Name: "url",
			Type: "string",
			Desc: "url地址",
		},
	},
	Handler: func(ctx context.Context, param model.FuncParam) (res any, err error) {
		url := param["url"].String()

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
	Handler: func(ctx context.Context, param model.FuncParam) (res any, err error) {

		data := param["data"].String()

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

func latestVersion() config.Func {
	return config.Func{
		Handler: func(ctx context.Context, param model.FuncParam) (res any, err error) {
			url := "https://api.github.com/repos/glennliao/bookmark/releases/latest"
			resp, err := gclient.New().Get(ctx, url)
			if err != nil {
				g.Log().Error(ctx, err)
				return "v0.0.0", nil
			}

			json := gjson.New(resp.ReadAllString())

			return json.Get("tag_name"), nil
		},
	}
}

func noteTags() config.Func {
	return config.Func{
		ParamList: nil,
		Batch:     false,
		Handler: func(ctx context.Context, param model.FuncParam) (res any, err error) {
			user, _ := ctx.Value(UserIdKey).(*CurrentUser)

			type TagRecord struct {
				Tags []string
			}

			var list []TagRecord

			var tagSet gset.StrSet

			err = g.DB().Model("note").Where("group_id", user.UserId).Scan(&list)
			if err != nil {
				return nil, err
			}

			if len(list) == 0 {
				return []string{}, nil
			}

			for _, record := range list {
				for _, tag := range record.Tags {
					tagSet.Add(tag)
				}

			}

			return tagSet.Slice(), err
		},
	}
}

func bmTags() config.Func {
	return config.Func{
		ParamList: nil,
		Batch:     false,
		Handler: func(ctx context.Context, param model.FuncParam) (res any, err error) {
			user, _ := ctx.Value(UserIdKey).(*CurrentUser)

			type TagRecord struct {
				Tags []string
			}

			var list []TagRecord

			var tagSet gset.StrSet

			err = g.DB().Model("bookmark").Where("group_id", user.UserId).Scan(&list)
			if err != nil {
				return nil, err
			}

			if len(list) == 0 {
				return []string{}, nil
			}

			for _, record := range list {
				for _, tag := range record.Tags {
					tagSet.Add(tag)
				}

			}

			return tagSet.Slice(), err
		},
	}
}

func cateBmNum() config.Func {
	return config.Func{
		ParamList: nil,
		Batch:     false,
		Handler: func(ctx context.Context, param model.FuncParam) (res any, err error) {
			groupId := ""
			if param["groupId"] == nil {
				groupId = ctx.Value(UserIdKey).(*CurrentUser).UserId
			} else {
				groupId = param["groupId"].String()
			}
			values, err := g.DB().Model("bookmark").Where("group_id", groupId).WhereNull("drop_at").Group("cate_id").Fields("cate_id cateId, count(1) cnt").All()

			return values.List(), err
		},
	}
}

func fetchMetaBatchFunc() config.Func {
	return config.Func{
		Handler: func(ctx context.Context, param model.FuncParam) (res any, err error) {
			all, err := g.DB().Ctx(ctx).Model("bookmark").Where("description", "@import ").All()
			if err != nil {
				return nil, err
			}

			task := func(data interface{}) {

				ctx := gctx.NeverDone(ctx)

				item := data.(g.Map)
				url := gconv.String(item["url"])
				bmId := gconv.String(item["bm_id"])
				var meta *fetchurl.UrlMeta
				if !strings.HasPrefix(strings.ToLower(url), "http") {
					meta, err = fetchurl.FetchURLMeta(ctx, "https://"+url)
					if err != nil {
						meta, err = fetchurl.FetchURLMeta(ctx, "http://"+url)
						if err != nil {
							g.Log().Error(ctx, err, url, bmId)
							return
						}
					}
				} else {
					meta, err = fetchurl.FetchURLMeta(ctx, url)
					if err != nil {
						g.Log().Error(ctx, err, url, bmId)
						return
					}
				}

				if meta != nil {
					if strings.HasPrefix(meta.Icon, "/") {
						_url, err := url2.Parse(meta.Url)
						if err != nil {
							g.Log().Error(ctx, err, url, bmId)
							return
						}
						meta.Icon = _url.Scheme + "://" + _url.Host + meta.Icon
					}

					meta.Url = url
				}

				_, err := g.DB().Model("bookmark").Ctx(ctx).Where("bm_id", bmId).Update(g.Map{
					"icon":        meta.Icon,
					"description": meta.Description,
				})
				if err != nil {
					g.Log().Error(ctx, err)
				}

			}

			p, _ := ants.NewPoolWithFunc(8, task)
			defer p.Release()

			for _, item := range all {
				err := p.Invoke(item.Map())
				if err != nil {
					return nil, err
				}
			}
			return nil, err
		},
	}
}

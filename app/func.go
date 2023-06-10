package app

import (
	"context"
	url2 "net/url"
	"strings"

	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/bookmark/app/util/fetchurl"
	"github.com/gogf/gf/v2/util/gconv"
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

			if meta != nil {
				if strings.HasPrefix(meta.Icon, "/") {
					_url, _ := url2.Parse(meta.Url)
					meta.Icon = _url.Scheme + "://" + _url.Host + meta.Icon
				}

				meta.Url = url
			}

			return meta, err
		},
	})

}

package app

import (
	"context"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/bookmark/app/util/fetchurl"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

func initFunc(a *apijson.ApiJson) {
	a.Config().Functions.BindHandlerFunc("fetchURL", func(ctx context.Context, param model.Map) (res any, err error) {
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
	})
}

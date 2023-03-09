package fetchurl

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/glennliao/getfavicon"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
)

type UrlFetchResolver func(ctx context.Context, url string, doc *goquery.Document) (*UrlMeta, error)

var urlFetchResolverList []UrlFetchResolver

func init() {
	urlFetchResolverList = append(urlFetchResolverList, gitee)
	urlFetchResolverList = append(urlFetchResolverList, github)
	urlFetchResolverList = append(urlFetchResolverList, common)
}

func gitee(ctx context.Context, url string, doc *goquery.Document) (meta *UrlMeta, err error) {
	if strings.HasPrefix(url, "https://gitee.com/") {
		meta = &UrlMeta{
			Url:  url,
			Icon: "https://gitee.com/favicon.ico",
		}

		description := strings.TrimSpace(doc.Find(".git-project-desc-text").Text())
		meta.Description = description

		title, _ := doc.Find("input[name=project_title]").Attr("value")
		meta.Title = strings.Trim(title, "/")
		return
	}
	return nil, nil
}

func github(ctx context.Context, url string, doc *goquery.Document) (meta *UrlMeta, err error) {
	if strings.HasPrefix(url, "https://github.com/") {
		meta = &UrlMeta{
			Url:  url,
			Icon: "https://github.githubassets.com/favicons/favicon.svg",
		}

		description := strings.TrimSpace(doc.Find(".Layout-sidebar .f4").Text())
		meta.Description = description

		title, _ := doc.Find("strong[itemprop=name]>a").Attr("href")
		meta.Title = strings.Trim(title, "/")
		return
	}
	return nil, nil
}

func common(ctx context.Context, url string, doc *goquery.Document) (meta *UrlMeta, err error) {

	meta = &UrlMeta{
		Url: url,
	}

	head := doc.Find("head")

	// 查找meta属性项charset
	title := head.Find("title").Text()
	description, _ := head.Find("meta[name=description]").Attr("content")

	meta.Description = description
	meta.Title = title

	if meta.Icon == "" {
		meta.Icon, err = getfavicon.Get(url)
		if err != nil {
			g.Log().Error(ctx, err)
		}
	}

	return meta, nil
}

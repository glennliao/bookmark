package fetchurl

import (
	"context"
	"crypto/tls"
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/v2/frame/g"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type UrlMeta struct {
	Title       string `json:"title" des:"标题"`
	Url         string `json:"url" des:"url"`
	Icon        string `json:"icon" des:"图标"`
	Description string `json:"description" des:"描述"`
}

func FetchURLMeta(ctx context.Context, uri string) (meta *UrlMeta, err error) {
	client := http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(uri)

	if err != nil {
		if e, ok := err.(*url.Error); ok {
			if !e.Timeout() {
				return &UrlMeta{}, err
			}
		} else {
			return &UrlMeta{}, err
		}
	}

	// for china
	if err != nil && strings.HasPrefix(uri, "https://github.com") {
		uri = strings.Replace(uri, "https://github.com/", "https://hub.nuaa.cf/", 1)
		resp, err = client.Get(uri)
		if err != nil && err != http.ErrHandlerTimeout {
			return
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		g.Log().Warningf(ctx, "uri:%s status code error: %d %s", uri, resp.StatusCode, resp.Status)
	}

	doc, err2 := goquery.NewDocumentFromReader(resp.Body)
	if err != nil || err2 != nil {
		g.Log().Error(ctx, err)
		return
	}

	for _, resolver := range urlFetchResolverList {
		ret, err := resolver(ctx, uri, doc)
		if ret != nil {
			return ret, err
		}
	}
	return &UrlMeta{
		Url: uri,
	}, nil

}

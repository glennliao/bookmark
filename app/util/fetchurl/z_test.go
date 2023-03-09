package fetchurl

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/v2/frame/g"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	ctx := context.Background()
	uri := "https://github.com/glennliao/apijson-go"
	resp, err := client.Get(uri)

	if err != nil && err != http.ErrHandlerTimeout {
		return
	}

	// for china
	if err == http.ErrHandlerTimeout && strings.HasPrefix(uri, "https://github.com") {
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

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	title, _ := doc.Find("strong[itemprop=name]>a").Attr("href")

	title = strings.Trim(title, "/")
	fmt.Println(title)
}

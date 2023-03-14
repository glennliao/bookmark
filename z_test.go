package main

import (
	"context"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/bookmark/app"
	"github.com/gogf/gf/v2/frame/g"
	"log"
	"testing"
)

func TestName(t *testing.T) {
	a := apijson.Load(app.Init)

	ctx := context.Background()

	ctx = context.WithValue(ctx, app.UserIdKey, &app.CurrentUser{UserId: "353368957608005"})

	q := a.NewQuery(ctx, model.Map{
		"Bookmark": model.Map{
			"info()": "fetchURL(url)",
		},
		"url@":       "Bookmark/url",
		"url2":       "Bookmark/url",
		"fetchURL()": "fetchURL(url@)",
	})

	ret, err := q.Result()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	g.Dump(ret)
}

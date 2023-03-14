package app

import (
	"context"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"net/http"
)

func initHook(a *apijson.ApiJson) {
	action.RegHook(action.Hook{
		For: "*",
		BeforeExecutorDo: func(ctx context.Context, n *action.Node, method string) error {
			user, ok := ctx.Value(UserIdKey).(*CurrentUser)
			if ok {
				for i := range n.Data {
					if method == http.MethodPost || method == http.MethodPut {
						n.Data[i]["updated_by"] = user.UserId
						if method == http.MethodPost {
							n.Data[i]["created_by"] = user.UserId
						}
					}
				}
			}

			return nil
		},
		AfterExecutorDo: func(ctx context.Context, n *action.Node, method string) error {
			if n.Key == TableUser && method == http.MethodPost {
				for _, item := range n.Data {
					act := a.NewAction(ctx, http.MethodPost, model.Map{
						"tag": "Groups",
						"Groups": model.Map{
							"groupId": item["user_id"],
							"title":   "个人分组",
						},
						"GroupUser[]": []model.Map{
							{
								"groupId": item["user_id"],
								"userId":  item["user_id"],
							},
						},
					})
					act.NoAccessVerify = true
					_, err := act.Result()
					if err != nil {
						return gerror.Wrap(err, "AfterExecutorDo")
					}

				}
			}

			if method == http.MethodPut && n.Key == TableBookmarkUse {
				if gconv.Int64(n.Ret["count"]) == 0 {
					m := g.DB().Model("bookmark_use").Safe().Ctx(ctx)

					user, _ := ctx.Value(UserIdKey).(*CurrentUser)
					m.Insert(g.Map{
						"times":   "1",
						"bm_id":   util.String(n.Where[0]["bm_id"]),
						"user_id": user.UserId,
					})
				}
			}

			return nil
		},
	})
}

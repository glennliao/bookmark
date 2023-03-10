package app

import (
	"context"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/errors/gerror"
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
			if n.AccessName == TableUser && method == http.MethodPost {
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

			return nil
		},
	})
}

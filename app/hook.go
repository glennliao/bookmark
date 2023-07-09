package app

import (
	"context"
	"net/http"

	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func initHook(a *apijson.ApiJson) {
	action.RegHook(action.Hook{
		For:            []string{"Note"},
		BeforeNodeExec: nil,
		AfterNodeExec:  nil,
		BeforeExecutorDo: func(ctx context.Context, n *action.Node, method string) error {
			user, _ := ctx.Value(UserIdKey).(*CurrentUser)

			for i, _ := range n.Data {
				n.Data[i]["group_id"] = user.UserId
			}

			return nil
		},
		AfterExecutorDo: func(ctx context.Context, n *action.Node, method string) error {

			return nil
		},
	})

	action.RegHook(action.Hook{
		For: []string{"*"},
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

			if method == http.MethodDelete {
				if n.Key == TableBookmarkCate {
					for _, item := range n.Where {
						cateId := item["cate_id"]

						one, err := g.DB().Model("bookmark_cate").Where("parent_id", cateId).Safe().Ctx(ctx).One()
						if err != nil {
							return err
						}

						if !one.IsEmpty() {
							return gerror.Newf("分类下有子分类，不能删除，cate_id: %v", cateId)
						}

						one, err = g.DB().Model("group_bookmark").Where("cate_id", cateId).Safe().Ctx(ctx).One()
						if err != nil {
							return err
						}

						if !one.IsEmpty() {
							return gerror.Newf("分类下有书签，不能删除，cate_id: %v", cateId)
						}
					}
				}
			}

			return nil
		},
		AfterExecutorDo: func(ctx context.Context, n *action.Node, method string) error {

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

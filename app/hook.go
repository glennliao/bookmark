package app

import (
	"context"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func initHook(a *apijson.ApiJson) {

	a.RegActionHook(action.Hook{
		For:     []string{"Note"},
		Handler: nil,
		HandlerInTransaction: func(ctx context.Context, req *action.HookReq) error {
			user, _ := ctx.Value(UserIdKey).(*CurrentUser)
			for i, _ := range req.Node.Data {
				req.Node.Data[i]["group_id"] = user.UserId
			}
			return req.Next()
		},
	})

	a.RegActionHook(action.Hook{
		For:     []string{"*"},
		Handler: nil,
		HandlerInTransaction: func(ctx context.Context, req *action.HookReq) error {
			user, ok := ctx.Value(UserIdKey).(*CurrentUser)
			if ok {
				for i := range req.Node.Data {
					if req.IsPost() || req.IsPut() {
						req.Node.Data[i]["updated_by"] = user.UserId
						if req.IsPost() {
							req.Node.Data[i]["created_by"] = user.UserId
						}
					}
				}
			}

			if req.IsDelete() {
				if req.Node.Key == TableBookmarkCate {
					for _, item := range req.Node.Where {
						cateId := item["cate_id"]

						one, err := g.DB().Model("bookmark_cate").Where("parent_id", cateId).Safe().Ctx(ctx).One()
						if err != nil {
							return err
						}

						if !one.IsEmpty() {
							return gerror.NewCodef(gcode.CodeInvalidParameter, "分类下有子分类，不能删除，cate_id: %v", cateId)
						}

						one, err = g.DB().Model("bookmark").Where("cate_id", cateId).Safe().Ctx(ctx).One()
						if err != nil {
							return err
						}

						if !one.IsEmpty() {
							return gerror.NewCodef(gcode.CodeInvalidParameter, "分类下有书签，不能删除，cate_id: %v", cateId)
						}
					}
				}
			}

			if req.IsPut() && req.Node.Key == TableBookmarkUse {

				for i := range req.Node.Where {
					req.Node.Where[i]["bm_id"] = req.Node.Data[i]["bm_id"]
				}
			}

			err := req.Next()

			if req.IsPut() && req.Node.Key == TableBookmarkUse {
				if gconv.Int64(req.Node.Ret["count"]) == 0 {
					m := g.DB().Model("bookmark_use").Safe().Ctx(ctx)

					user, _ := ctx.Value(UserIdKey).(*CurrentUser)
					m.Insert(g.Map{
						"times":   "1",
						"bm_id":   util.String(req.Node.Where[0]["bm_id"]),
						"user_id": user.UserId,
					})
				}
			}
			return err
		},
	})

}

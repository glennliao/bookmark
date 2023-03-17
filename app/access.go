package app

import (
	"context"
	"fmt"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
	"net/http"
)

const (
	TableUser          = "User"
	TableBookmark      = "Bookmark"
	TableBookmarkUse   = "BookmarkUse"
	TableBookmarkCate  = "BookmarkCate"
	TableGroups        = "Groups"
	TableGroupUser     = "GroupUser"
	TableGroupBookmark = "GroupBookmark"
)

const UserIdKey = "userId"

type CurrentUser struct {
	UserId string
}

const RoleGroupUser = "GroupUser"
const RoleGroupAdmin = "GroupAdmin"

func Role(ctx context.Context, req config.RoleReq) (string, error) {
	_, ok := ctx.Value(UserIdKey).(*CurrentUser)

	if !ok {
		return consts.UNKNOWN, nil //未登录
	}

	if req.NodeRole == "" {

		switch req.AccessName {
		case TableUser:
			return consts.OWNER, nil
		}

	} else {

		switch req.AccessName {
		case TableUser:
			if req.NodeRole == consts.LOGIN {
				return consts.OWNER, nil
			}

		case TableBookmark, TableBookmarkUse:

			if req.NodeRole == consts.LOGIN {
				req.NodeRole = consts.OWNER
			}

			if lo.Contains([]string{consts.OWNER, RoleGroupUser}, req.NodeRole) {
				return req.NodeRole, nil
			}

			return consts.DENY, nil // 非拥有的角色

		default:
			return req.NodeRole, nil
		}
	}

	return consts.LOGIN, nil

}

func AccessCondition(ctx context.Context, req config.ConditionReq, where *config.ConditionRet) (err error) {

	user, ok := ctx.Value(UserIdKey).(*CurrentUser)
	if !ok {
		where.Add("1", 2)
		return nil
	}

	switch req.AccessName {
	case TableUser:
		if req.NodeRole == consts.OWNER {
			where.Add("user_id", user.UserId)
		}
	case TableBookmarkUse:

		where.Add("bm_id", req.NodeReq["bmId"])
		delete(req.NodeReq, "bmId")
		where.Add("user_id", user.UserId)

	case TableBookmarkCate:
		where.AddRaw("group_id in (select group_id from group_user where user_id = ?)", user.UserId)
	case TableBookmark:
		if req.Method == http.MethodGet {
			if v, exists := req.NodeReq["cateId"]; exists {
				delete(req.NodeReq, "cateId")
				where.AddRaw("bm_id in (select bm_id from group_bookmark where  drop_at is null and cate_id = ? and group_id in (select group_id from group_user where user_id = ? ))", []string{v.(string), user.UserId})
			} else {
				where.AddRaw("bm_id in (select bm_id from group_bookmark where drop_at is null and group_id in (select group_id from group_user where user_id = ? ))", []string{user.UserId})
			}
		} else {
			where.AddRaw("bm_id in (select bm_id from group_bookmark where drop_at is null and group_id in (select group_id from group_user where user_id = ? ))", []string{user.UserId})
		}

	case TableGroupBookmark:
		if req.Method == http.MethodPut {

			if _, exists := req.NodeReq["dropAt"]; exists {
				where.Add("bm_id", req.NodeReq["bmId"])
				where.Add("group_id", req.NodeReq["groupId"])
				where.Add("cate_id", req.NodeReq["cateId"])
				delete(req.NodeReq, "bmId")
				delete(req.NodeReq, "groupId")
				delete(req.NodeReq, "cateId")
			} else {

			}
			// todo 添加个@where op? 用于专门指定where条件
		}
	case TableGroups:
		if req.Method == http.MethodGet {
			where.AddRaw("group_id in (select group_id from group_user where user_id = ? )", []string{user.UserId})
		}

	}
	return nil
}

var jwtSecret []byte

func init() {
	val, err := g.Cfg().Get(nil, "jwt.secret")
	if err != nil {
		g.Log().Fatal(nil, err)
	}
	jwtSecret = val.Bytes()
	if len(jwtSecret) < 16 {
		g.Log().Fatal(nil, "jwt.secret不能小于16位")
	}
}

func Auth(r *ghttp.Request) {

	authorization := r.Request.Header.Get("Authorization")

	if authorization != "" {
		ctx := r.Context()
		token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return jwtSecret, nil
		})

		if err != nil {
			g.Log().Info(ctx, "token err", err)
			r.Response.WriteJson(g.Map{
				"code": 401,
				"msg":  "未知错误",
			})
			return
		}

		userId := ""
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId = gconv.String(claims["userId"])
			if userId == "" {
				g.Log().Info(ctx, "userId is empty")
				r.Response.WriteJson(g.Map{
					"code": 401,
					"msg":  "未知错误",
				})
				return
			}
		} else {
			g.Log().Info(ctx, "token no valid")
			r.Response.WriteJson(g.Map{
				"code": 401,
				"msg":  "未知错误",
			})
			return
		}
		ctx = context.WithValue(ctx, UserIdKey, &CurrentUser{UserId: userId})
		r.SetCtx(ctx)
	} else {
		if r.Request.URL.Path != "/api/data/auth" {
			r.Response.WriteJson(g.Map{
				"code": 401,
				"msg":  "未登录",
			})
			return
		}

	}

	r.Middleware.Next()
}

func Cors(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	corsOptions.AllowOrigin = r.Request.Header.Get("Origin")
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

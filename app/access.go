package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
)

const (
	TableUser         = "User"
	TableBookmark     = "Bookmark"
	TableBookmarkUse  = "BookmarkUse"
	TableBookmarkCate = "BookmarkCate"
	TableGroups       = "Groups"
	TableGroupUser    = "GroupUser"
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
		return consts.UNKNOWN, nil // 未登录
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
		where.Add("user_id", user.UserId)

	case TableBookmarkCate:
		where.AddRaw("group_id in (select group_id from group_user where user_id = ?)", user.UserId)
	case TableBookmark:
		if req.Method == http.MethodGet {
			if v, exists := req.NodeReq["q"]; exists {
				delete(req.NodeReq, "q")

				q := fmt.Sprintf("%s", gconv.String(v))

				words := strings.Split(q, " ")

				sql := ""
				var params []string
				for i, word := range words {

					if i > 0 {
						sql += " and "
					}
					sql += fmt.Sprintf("(lower(title) like ? or lower(url) like ? or lower(remark) like ?)")
					word = strings.ToLower(word)
					params = append(params, "%"+word+"%", "%"+word+"%", "%"+word+"%")
				}

				if len(params) > 0 {
					where.AddRaw("("+sql+")", params)
				}

			}

			where.AddRaw("drop_at is null and group_id in (select group_id from group_user where user_id = ? )", []string{user.UserId})
		} else {
			where.AddRaw("drop_at is null and group_id in (select group_id from group_user where user_id = ? )", []string{user.UserId})
		}

	case "Note":
		where.Add("group_id", user.UserId)

		if req.Method == http.MethodGet {
			if _, exists := req.NodeReq["tag"]; exists {
				where.AddRaw("json_extract(tags,'$') like ?", "%"+gconv.String(req.NodeReq["tag"])+"%")
				delete(req.NodeReq, "tag")
			}
		}

	case TableGroups:
		if req.Method == http.MethodGet {
			where.AddRaw("group_id in (select group_id from group_user where user_id = ? )", []string{user.UserId})
		}

	}
	return nil
}

package app

import (
	"context"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/bookmark/app/inits"
	"github.com/glennliao/table-sync/tablesync"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	config.RegAccessListProvider("json", func(ctx context.Context) []config.AccessConfig {
		var accessList []config.AccessConfig
		err := gconv.Scan(inits.AccessConfigJson, &accessList)
		if err != nil {
			panic(err)
		}
		return accessList
	})

	config.RegRequestListProvider("json", func(ctx context.Context) []config.Request {
		var requestList []config.Request
		err := gconv.Scan(inits.RequestConfigJson, &requestList)
		if err != nil {
			panic(err)
		}
		for i, request := range requestList {
			if _, ok := request.Structure[request.Tag]; !ok {
				requestList[i].Structure = map[string]*config.Structure{
					request.Tag: {
						Must:    nil,
						Refuse:  nil,
						Unique:  nil,
						Insert:  nil,
						Update:  nil,
						Replace: nil,
						Remove:  nil,
					},
				}
			}
		}
		return requestList
	})
}

func Init(ctx context.Context, a *apijson.ApiJson) {

	a.Config().Access.ConditionFunc = AccessCondition
	a.Config().Access.DefaultRoleFunc = Role
	a.Config().AccessListProvider = "json"
	a.Config().RequestListProvider = "json"
	a.Config().Access.AddRole([]string{RoleGroupAdmin, RoleGroupUser})

	initDB(ctx)
	initHook(a)
	initRowKeyGen(a)
	initFunc(a)
}

func initDB(ctx context.Context) {
	db := g.DB()
	syncer := tablesync.Syncer{Tables: inits.Tables()}
	err := syncer.Sync(ctx, db)
	if err != nil {
		panic(err)
	}
}

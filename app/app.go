package app

import (
	"context"

	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/config"
	jsonConfig "github.com/glennliao/apijson-go/drivers/json/config"
	"github.com/glennliao/bookmark/app/inits"
	"github.com/glennliao/table-sync/tablesync"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	config.RegAccessListProvider("json", jsonConfig.AccessListProvider(nil, inits.AccessConfigJson))
	config.RegRequestListProvider("json", jsonConfig.RequestListProvider(nil, inits.RequestConfigJson))
}

func App(ctx context.Context, a *apijson.ApiJson) {

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

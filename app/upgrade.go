package app

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"time"
)

// ToV0Dot11 升级到v0.11, 处理数据库数据调整
func ToV0Dot11(ctx context.Context) {
	all, err := g.DB().Model("group_bookmark").All()
	if err != nil {
		g.Log().Fatal(ctx, err)
	}

	if len(all) > 0 {
		g.Log().Info(ctx, "正在升级处理中")

		c := g.DB().GetConfig()
		if c.Type == "sqlite" {
			g.Log().Info(ctx, "正在备份sqlite数据库")
			err := gfile.CopyFile(c.Name, c.Name+"."+time.Now().Format("0102150405")+".bak")
			if err != nil {
				g.Log().Fatal(ctx, err)
			}
		}

		for _, record := range all {
			_, err := g.DB().Model("bookmark").Where("bm_id", record.Map()["bm_id"]).Update(g.Map{
				"drop_at":  record.Map()["drop_at"],
				"cate_id":  record.Map()["cate_id"],
				"group_id": record.Map()["group_id"],
				"sorts":    record.Map()["sorts"],
			})
			if err != nil {
				g.Log().Fatal(ctx, err)
			}
		}

		g.DB().Exec(ctx, "drop table group_bookmark")

		g.Log().Info(ctx, "升级处理完成")

	}
}

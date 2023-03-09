package inits

import (
	"github.com/glennliao/table-sync/tablesync"
	"time"
)

type Bookmark struct {
	tablesync.TableMeta `charset:"utf8mb4" comment:"书签"`
	Id                  uint64 `ddl:"primaryKey"`
	BmId                string `ddl:"size:32;comment:书签id"`
	Title               string `ddl:"size:128;comment:标题"`
	Url                 string `ddl:"size:512;comment:书签地址"`
	Icon                string `ddl:"size:256;comment:图标地址"`
	Description         string `ddl:"size:2048;comment:书签描述(来自网站)"`
	Remark              string `ddl:"size:1024;comment:书签描述(来自用户)"`
	CreatedAt           *time.Time
	CreatedBy           string `ddl:"size:32"`
	UpdatedAt           *time.Time
	UpdatedBy           string `ddl:"size:32"`
	DeletedAt           *time.Time
}

type BookmarkUse struct {
	tablesync.TableMeta `charset:"utf8mb4" comment:"书签使用次数"`
	Id                  uint64 `ddl:"primaryKey"`
	BmId                string `ddl:"size:32;comment:书签id"`
	UserId              string `ddl:"size:128;comment:标题"`
	Times               uint64 `ddl:"comment:使用次数"`
	UpdatedAt           *time.Time
}

type BookmarkCate struct {
	tablesync.TableMeta `charset:"utf8mb4" comment:"书签目录"`
	Id                  uint64 `ddl:"primaryKey"`
	CateId              string `ddl:"size:32;comment:目录id"`
	Title               string `ddl:"size:128;comment:群组名"`
	ParentId            string `ddl:"size:32;comment:父级目录id"`
	GroupId             string `ddl:"size:32;comment:群组id"`
	Sorts               int32
	CreatedAt           *time.Time
	CreatedBy           string `ddl:"size:32"`
	UpdatedAt           *time.Time
	UpdatedBy           string `ddl:"size:32"`
	DeletedAt           *time.Time
}

type GroupBookmark struct {
	tablesync.TableMeta `charset:"utf8mb4" comment:"书签目录"`
	Id                  uint64 `ddl:"primaryKey"`
	CateId              string `ddl:"size:32;comment:目录id"`
	BmId                string `ddl:"size:32;comment:书签id"`
	GroupId             string `ddl:"size:32;comment:群组id"`
	Sorts               int32
	CreatedAt           *time.Time
	CreatedBy           string `ddl:"size:32"`
	UpdatedAt           *time.Time
	UpdatedBy           string `ddl:"size:32"`
	DeletedAt           *time.Time
}

type Groups struct {
	tablesync.TableMeta `charset:"utf8mb4" comment:"群组"`
	Id                  uint64 `ddl:"primaryKey"`
	GroupId             string `ddl:"size:32;comment:群组id"`
	Title               string `ddl:"size:128;comment:群组名"`
	CreatedAt           *time.Time
	CreatedBy           string `ddl:"size:32"`
	UpdatedAt           *time.Time
	UpdatedBy           string `ddl:"size:32"`
	DeletedAt           *time.Time
}

type GroupUser struct {
	tablesync.TableMeta `charset:"utf8mb4" comment:"群组用户"`
	Id                  uint64 `ddl:"primaryKey"`
	GroupId             string `ddl:"size:32;comment:群组id"`
	UserId              string `ddl:"size:32;comment:用户id"`
	IsAdmin             uint8  `ddl:"comment:是否管理员"`
	CreatedAt           *time.Time
	CreatedBy           string `ddl:"size:32"`
	UpdatedAt           *time.Time
	UpdatedBy           string `ddl:"size:32"`
	DeletedAt           *time.Time
}

type User struct {
	tablesync.TableMeta `charset:"utf8mb4" comment:"用户"`
	Id                  uint64 `ddl:"primaryKey"`
	UserId              string `ddl:"size:32;comment:用户id"`
	Username            string `ddl:"size:512;comment:用户名"`
	Password            string `ddl:"size:64;comment:密码"`
	CreatedAt           *time.Time
	CreatedBy           string `ddl:"size:32"`
	UpdatedAt           *time.Time
	UpdatedBy           string `ddl:"size:32"`
	DeletedAt           *time.Time
}

func Tables() []tablesync.Table {
	return []tablesync.Table{
		// app
		Bookmark{},
		BookmarkUse{},
		BookmarkCate{},
		Groups{},
		GroupUser{},
		GroupBookmark{},
		User{},
	}
}

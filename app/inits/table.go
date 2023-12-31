package inits

import (
	"time"

	"github.com/glennliao/table-sync/tablesync"
)

type Created struct {
	CreatedAt *time.Time
	CreatedBy string `ddl:"size:32"`
}

type Updated struct {
	UpdatedAt *time.Time
	UpdatedBy string `ddl:"size:32"`
}

type Deleted struct {
	DeletedAt *time.Time
}

type (
	Groups struct {
		tablesync.TableMeta `charset:"utf8mb4" comment:"群组"`
		Id                  uint64 `ddl:"primaryKey"`
		GroupId             string `ddl:"size:32;comment:群组id"`
		Title               string `ddl:"size:128;comment:群组名"`
		Created
		Updated
		Deleted
	}

	GroupUser struct {
		tablesync.TableMeta `charset:"utf8mb4" comment:"群组用户"`
		Id                  uint64 `ddl:"primaryKey"`
		GroupId             string `ddl:"size:32;comment:群组id"`
		UserId              string `ddl:"size:32;comment:用户id"`
		IsAdmin             uint8  `ddl:"comment:是否管理员"`
		Created
		Updated
		Deleted
	}

	User struct {
		tablesync.TableMeta `charset:"utf8mb4" comment:"用户"`
		Id                  uint64 `ddl:"primaryKey"`
		UserId              string `ddl:"size:32;comment:用户id"`
		Username            string `ddl:"size:512;comment:用户名"`
		Password            string `ddl:"size:64;comment:密码"`
		Created
		Updated
		Deleted
	}
)

type (
	Bookmark struct {
		tablesync.TableMeta `charset:"utf8mb4" comment:"书签"`
		Id                  uint64     `ddl:"primaryKey"`
		BmId                string     `ddl:"size:32;comment:书签id"`
		Title               string     `ddl:"size:128;comment:标题"`
		Url                 string     `ddl:"size:512;comment:书签地址"`
		Icon                string     `ddl:"size:256;comment:图标地址"`
		Description         string     `ddl:"size:2048;comment:书签描述(来自网站)"`
		Remark              string     `ddl:"size:1024;comment:书签描述(来自用户)"`
		Tags                string     `ddl:"type:json;comment:标签"`
		EncodeKey           string     `ddl:"size:128;comment:加密key"`
		IsPublic            uint8      `ddl:"comment:是否公开"`
		CateId              string     `ddl:"size:32;comment:分类id"`
		GroupId             string     `ddl:"size:32;comment:群组id"`
		FromGroupId         string     `ddl:"size:32;comment:分享来源"`
		DropAt              *time.Time `ddl:"comment:放置到回收站时间"`
		Sorts               int32

		Created
		Updated
		Deleted
	}

	BookmarkUse struct {
		tablesync.TableMeta `charset:"utf8mb4" comment:"书签使用次数"`
		Id                  uint64 `ddl:"primaryKey"`
		BmId                string `ddl:"size:32;comment:书签id"`
		UserId              string `ddl:"size:128;comment:标题"`
		Times               uint64 `ddl:"comment:使用次数"`
		UpdatedAt           *time.Time
	}

	BookmarkCate struct {
		tablesync.TableMeta `charset:"utf8mb4" comment:"书签目录"`
		Id                  uint64 `ddl:"primaryKey"`
		CateId              string `ddl:"size:32;comment:目录id"`
		Title               string `ddl:"size:128;comment:群组名"`
		ParentId            string `ddl:"size:32;comment:父级目录id"`
		GroupId             string `ddl:"size:32;comment:群组id"`
		Sorts               int32
		Created
		Updated
		Deleted
	}
)

type (
	Note struct {
		tablesync.TableMeta `charset:"utf8mb4" comment:"笔记"`
		Id                  uint64     `ddl:"primaryKey"`
		GroupId             string     `ddl:"size:32;comment:groupId"`
		NoteId              string     `ddl:"size:32;comment:@rowKey noteId"`
		Content             string     `ddl:"type:json;comment:内容"`
		Tags                string     `ddl:"type:json;comment:标签"`
		IsPublic            uint8      `ddl:"comment:是否公开"`
		EncodeKey           string     `ddl:"size:128;comment:加密key"`
		DropAt              *time.Time `ddl:"comment:放置到回收站时间"`
		Created
		Updated
		Deleted
	}

	NoteHistory struct {
		tablesync.TableMeta `charset:"utf8mb4" comment:"笔记历史"`
		Id                  uint64 `ddl:"primaryKey"`
		NoteId              string `ddl:"size:32;comment:@rowKey noteId"`
		Content             string `ddl:"type:json;comment:内容"`
		IsPublic            uint8  `ddl:"comment:是否公开"`
		EncodeKey           string `ddl:"size:128;comment:加密key"`
		Created
	}
)

type (
	AppList struct {
		tablesync.TableMeta `charset:"utf8mb4" comment:"应用"`
		Id                  uint64 `ddl:"primaryKey"`
		AppId               string `ddl:"size:32;comment:appId"`
		Title               string `ddl:"size:128;comment:标题"`
		Url                 string `ddl:"size:512;comment:书签地址"`
		Icon                string `ddl:"size:256;comment:图标地址"`
		Remark              string `ddl:"size:1024;comment:备注"`
		From                string `ddl:"size:32;comment:来源"`
		UserId              string `ddl:"size:128;comment:标题"`
		Created
		Updated
		Deleted
	}
)

type (
	Tag struct {
		tablesync.TableMeta `charset:"utf8mb4" comment:"标签"`
		Id                  uint64 `ddl:"primaryKey"`
		Tag                 string `ddl:"size:32;comment:标签内容"`
		Created
	}

	Config struct {
		tablesync.TableMeta `charset:"utf8mb4" comment:"设置表"`
		Id                  uint64 `ddl:"primaryKey"`
		For                 string `ddl:"size:32;comment:system/user/group"`
		ForId               string `ddl:"size:32;"`
		Key                 string `ddl:"size:64;"`
		Value               string `ddl:"type:json"`
		Created
		Updated
		Deleted
	}
)

func Tables() []tablesync.Table {
	return []tablesync.Table{
		// app
		Bookmark{},
		BookmarkUse{},
		BookmarkCate{},

		Groups{},
		GroupUser{},
		User{},

		Note{},
		NoteHistory{},

		AppList{},

		Tag{},
		Config{},
	}
}

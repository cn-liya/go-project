package proto

import (
	"project/model/entity"
)

type SystemActionMenuSyncArgs struct {
	Menus []*SysMenu `json:"menus" binding:"required,dive"`
}
type SysMenu struct {
	Name  string     `json:"name" binding:"required"`
	Title string     `json:"title" binding:"required"`
	Sub   []*SysMenu `json:"sub,omitempty"`
}

type SystemSyncData struct {
	ActionDelete []string
	ActionCreate []*entity.SystemAction
	MenuTree     []*entity.SysMenuTree
	MenuKeys     []string
}

type UpdateSystemActionArgs struct {
	KeyName string `json:"key_name" binding:"required"`
	Title   string `json:"title" binding:"max=16"`
	Sort    int    `json:"sort" binding:"min=0,max=9999"`
}

type SystemActionMenuListResp struct {
	Menus   []*entity.SysMenuTree `json:"menus"`
	Actions []*SystemActionGroup  `json:"actions"`
}

type SystemActionGroup struct {
	KeyName string                 `json:"key_name"`
	Title   string                 `json:"title"`
	Sort    int                    `json:"sort"`
	Sub     []*entity.SystemAction `json:"sub"`
}

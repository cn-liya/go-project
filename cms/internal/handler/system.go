package handler

import (
	"github.com/gin-gonic/gin"
	"project/cms/internal/proto"
	"project/model/entity"
	"project/pkg/logger"
	"project/pkg/types"
	"project/pkg/web/errcode"
	"sort"
	"strings"
)

func dfsMenuTree(tree []*proto.SysMenu) ([]*entity.SysMenuTree, []string) {
	var nodes []*entity.SysMenuTree
	var keys []string
	for _, v := range tree {
		node := &entity.SysMenuTree{
			Name:  v.Name,
			Title: v.Title,
		}
		if len(v.Sub) > 0 {
			n, k := dfsMenuTree(v.Sub)
			node.Sub = n
			keys = append(keys, k...)
		} else {
			keys = append(keys, v.Name)
		}
		nodes = append(nodes, node)
	}
	return nodes, keys
}

func (h *Handler) SystemActionMenuSync(g *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r proto.SystemActionMenuSyncArgs
		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(errcode.InvalidParam.Response())
			return
		}
		groups, err := h.service.SystemActionKeyNames(c, entity.SystemActionGroup)
		if err != nil {
			logger.FromContext(c).Error("service.SystemActionKeyNames error", entity.SystemActionGroup, err)
			c.JSON(errcode.FromError(err))
			return
		}
		actions, err := h.service.SystemActionKeyNames(c, entity.SystemActionRoute)
		if err != nil {
			logger.FromContext(c).Error("service.SystemActionKeyNames error", entity.SystemActionRoute, err)
			c.JSON(errcode.FromError(err))
			return
		}
		groupSet := make(map[string]struct{})
		actionMap := types.SliceToMap(actions)
		routes := g.Routes()
		actionCreate := make([]*entity.SystemAction, 0, len(routes))
		sortNum := 1
		for _, v := range routes {
			paths := strings.Split(v.Path[1:], "/")
			if len(paths) < 2 || paths[0] == "user" {
				continue
			}
			groupSet[paths[0]] = struct{}{}
			if _, ok := actionMap[v.Path]; ok { // 存在则跳过
				delete(actionMap, v.Path)
			} else { // 不存在则插入
				actionCreate = append(actionCreate, &entity.SystemAction{
					KeyName: v.Path,
					Level:   entity.SystemActionRoute,
					Sort:    sortNum,
				})
			}
			sortNum++
		}
		groupMap := types.SliceToMap(groups)
		sortNum = 1
		for key := range groupSet {
			if _, ok := groupMap[key]; ok { // 存在则跳过
				delete(groupMap, key)
			} else { // 不存在则插入
				actionCreate = append(actionCreate, &entity.SystemAction{
					KeyName: key,
					Level:   entity.SystemActionGroup,
					Sort:    sortNum,
				})
			}
			sortNum++
		}
		actionDelete := append(types.Keys(actionMap), types.Keys(groupMap)...) // 多余的删除
		trees, keys := dfsMenuTree(r.Menus)
		err = h.service.SystemActionMenuSave(c, &proto.SystemSyncData{
			ActionDelete: actionDelete,
			ActionCreate: actionCreate,
			MenuTree:     trees,
			MenuKeys:     keys,
		})
		if err != nil {
			logger.FromContext(c).Error("service.SystemActionMenuSave error", nil, err)
			c.JSON(errcode.FromError(err))
			return
		}
		h.SystemActionMenuList(c)
	}
}

func (h *Handler) SystemActionMenuList(c *gin.Context) {
	menus, err := h.service.SystemMenuTrees(c)
	if err != nil {
		logger.FromContext(c).Error("service.SystemMenuTrees error", nil, err)
		c.JSON(errcode.FromError(err))
		return
	}
	data, err := h.service.SystemActionList(c)
	if err != nil {
		logger.FromContext(c).Error("service.SystemActionList error", nil, err)
		c.JSON(errcode.FromError(err))
		return
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Sort < data[j].Sort
	})
	var groups []*proto.SystemActionGroup
	routes := make(map[string][]*entity.SystemAction)
	for _, v := range data {
		if v.Level == entity.SystemActionGroup {
			groups = append(groups, &proto.SystemActionGroup{
				KeyName: v.KeyName,
				Title:   v.Title,
				Sort:    v.Sort,
			})
		} else {
			group := v.KeyName[1 : strings.Index(v.KeyName[1:], "/")+1]
			routes[group] = append(routes[group], v)
		}
	}
	for _, v := range groups {
		v.Sub = routes[v.KeyName]
	}
	c.JSON(OK, &proto.SystemActionMenuListResp{
		Menus:   menus,
		Actions: groups,
	})
}

func (h *Handler) SystemActionUpdate(c *gin.Context) {
	var r proto.UpdateSystemActionArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.FromError(err))
		return
	}
	err := h.service.SystemActionUpdate(c, &entity.SystemAction{
		KeyName: r.KeyName,
		Title:   r.Title,
		Sort:    r.Sort,
	})
	if err != nil {
		logger.FromContext(c).Error("service.SystemActionUpdate error", &r, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, Empty)
}

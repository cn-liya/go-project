package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"github.com/gin-gonic/gin"
	"project/cms/internal/proto"
	"project/model/cache"
	"project/model/entity"
	"project/pkg/ginutil/bizresp"
	"project/pkg/logger"
	"project/pkg/random"
	"project/pkg/types"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) captchaSign(exp, code string) string {
	mac := hmac.New(sha1.New, []byte(h.captcha))
	mac.Write([]byte(exp))
	mac.Write([]byte(strings.ToUpper(code)))
	return base32.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (h *Handler) Captcha(c *gin.Context) {
	code := random.Strings(4)
	bin := h.drawer.Draw(code)
	exp := strconv.FormatInt(time.Now().Unix()+65, 10)
	key := exp + "." + h.captchaSign(exp, code)
	c.JSON(OK, &proto.CaptchaResp{
		SessionKey:  key,
		Base64Image: bin,
	})
}

func (h *Handler) UserLogin(c *gin.Context) {
	var r proto.LoginArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(bizresp.InvalidParam.Reply())
		return
	}

	sli := strings.Split(r.SessionKey, ".")
	if len(sli) != 2 {
		c.JSON(bizresp.InvalidParam.Reply())
		return
	}
	exp, _ := strconv.ParseInt(sli[0], 10, 64)
	if time.Now().Unix() > exp {
		c.JSON(bizresp.CaptchaExpired.Reply())
		return
	}
	sign := h.captchaSign(sli[0], r.Captcha)
	if sli[1] != sign {
		c.JSON(bizresp.CaptchaWrong.Reply())
		return
	}

	user, err := h.service.AdminUserFindByName(c, r.Username)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserFindByName error", r.Username, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	if user.ID == 0 || !user.CheckPassword(r.Password) {
		c.JSON(bizresp.UserOrPwdWrong.Reply())
		return
	}
	pwdExp := user.PwdExp - time.Now().Unix()
	if pwdExp < 0 {
		c.JSON(bizresp.PasswordExpired.Reply())
		return
	}
	if user.Status != entity.StatusOn {
		c.JSON(bizresp.AccountDisabled.Reply())
		return
	}

	store := &cache.AdminToken{
		ID:       user.ID,
		Username: user.Username,
	}
	resp := &proto.LoginResp{Username: user.Username}
	if user.IsSuper() {
		store.IsSuper = true
		resp.IsSuper = true
	} else if user.AdminRole != nil {
		store.Actions = user.AdminRole.Actions
		resp.Actions = user.AdminRole.Actions
		resp.Menus = user.AdminRole.Menus
	}
	token, err := h.service.AdminTokenSet(c, store)
	if err != nil {
		logger.FromContext(c).Error("service.AdminTokenSet error", user, err)
		c.JSON(bizresp.WithErr(err))
		return
	}

	if day := (pwdExp + 43200) / 86400; day < 2 {
		resp.PwdTip = "密码即将过期，请立即修改！"
	} else if day < 16 {
		resp.PwdTip = "密码将于" + strconv.FormatInt(day, 10) + "天内过期，请尽快修改！"
	}
	resp.Token = token
	c.JSON(OK, resp)
}

func (h *Handler) UserLogout(c *gin.Context) {
	err := h.service.AdminTokenDel(c, c.GetHeader("Authorization"))
	if err != nil {
		logger.FromContext(c).Error("service.AdminTokenDel error", nil, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) UserPassword(c *gin.Context) {
	var r proto.UserPasswordArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(bizresp.InvalidParam.Reply())
		return
	}
	auth := getAuth(c)
	user, err := h.service.AdminUserFindByID(c, auth.ID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserFindByID error", auth.ID, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	if user.CheckPassword(r.Password) {
		c.JSON(bizresp.NeedModified.WithMsg("新密码不能与原密码相同"))
		return
	}
	err = h.service.AdminUserUpdate(c, &entity.AdminUser{
		ID:       auth.ID,
		Password: r.Password,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserUpdate error", auth.ID, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminRoleList(c *gin.Context) {
	var r proto.ListArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(bizresp.InvalidParam.Reply())
		return
	}
	total, list, err := h.service.AdminRoleList(c, &r)
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoleList error", &r, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	items := make([]*proto.AdminRoleItem, 0, len(list))
	for _, v := range list {
		items = append(items, &proto.AdminRoleItem{
			ID:      v.ID,
			Name:    v.Name,
			Actions: v.Actions,
			Menus:   v.Menus,
		})
	}
	c.JSON(OK, &proto.AdminRoleListResp{
		Total: total,
		List:  items,
	})
}

func (h *Handler) AdminRoleOption(c *gin.Context) {
	data, err := h.service.AdminRoles(c)
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoles error", nil, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	opt := make([]*proto.OptionItem, 0, len(data))
	for _, d := range data {
		opt = append(opt, &proto.OptionItem{
			ID:   d.ID,
			Name: d.Name,
		})
	}
	c.JSON(OK, &proto.OptionResp{List: opt})
}

func (h *Handler) AdminRoleCreate(c *gin.Context) {
	var r proto.AdminRoleCreateArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(bizresp.InvalidParam.Reply())
		return
	}
	allActions, err := h.service.SystemActionKeyNames(c, entity.SystemActionRoute)
	if err != nil {
		logger.FromContext(c).Error("service.SystemActionKeyNames error", nil, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	allMenus, err := h.service.SystemMenuKeyNames(c)
	if err != nil {
		logger.FromContext(c).Error("service.SystemMenuKeyNames error", nil, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	actionMap := types.SliceToMap(allActions)
	var actions []string
	for _, v := range r.Actions {
		if _, ok := actionMap[v]; ok {
			actions = append(actions, v)
			delete(actionMap, v)
		}
	}
	menuMap := types.SliceToMap(allMenus)
	var menus []string
	for _, v := range r.Menus {
		if _, ok := menuMap[v]; ok {
			menus = append(menus, v)
			delete(menuMap, v)
		}
	}
	err = h.service.AdminRoleCreate(c, &entity.AdminRole{
		Name:    r.Name,
		Actions: actions,
		Menus:   menus,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoleCreate error", nil, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminRoleUpdate(c *gin.Context) {
	var r proto.AdminRoleUpdateArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(bizresp.InvalidParam.Reply())
		return
	}
	role, err := h.service.AdminRoleFindByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoleFindByID error", r.ID, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	if role.ID == 0 {
		c.JSON(bizresp.InvalidAssociatedID.Reply())
		return
	}
	allActions, err := h.service.SystemActionKeyNames(c, entity.SystemActionRoute)
	if err != nil {
		logger.FromContext(c).Error("service.SystemActionKeyNames error", nil, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	allMenus, err := h.service.SystemMenuKeyNames(c)
	if err != nil {
		logger.FromContext(c).Error("service.SystemMenuKeyNames error", nil, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	actionMap := types.SliceToMap(allActions)
	var actions []string
	for _, v := range r.Actions {
		if _, ok := actionMap[v]; ok {
			actions = append(actions, v)
			delete(actionMap, v)
		}
	}
	menuMap := types.SliceToMap(allMenus)
	var menus []string
	for _, v := range r.Menus {
		if _, ok := menuMap[v]; ok {
			menus = append(menus, v)
			delete(menuMap, v)
		}
	}
	err = h.service.AdminRoleUpdate(c, &entity.AdminRole{
		ID:      r.ID,
		Name:    r.Name,
		Actions: actions,
		Menus:   menus,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoleUpdate error", nil, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminUserList(c *gin.Context) {
	var r proto.AdminUserListArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(bizresp.InvalidParam.Reply())
		return
	}
	total, list, err := h.service.AdminUserList(c, &r)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserList error", &r, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	items := make([]*proto.AdminUserItem, 0, len(list))
	for _, v := range list {
		items = append(items, &proto.AdminUserItem{
			ID:       v.ID,
			Username: v.Username,
			RoleID:   v.RoleID,
			IsSuper:  v.IsSuper(),
			Status:   v.Status,
		})
	}
	c.JSON(OK, &proto.AdminUserListResp{
		Total: total,
		List:  items,
	})
}

func (h *Handler) AdminUserCreate(c *gin.Context) {
	var r proto.AdminUserCreateArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(bizresp.InvalidParam.Reply())
		return
	}
	role, err := h.service.AdminRoleFindByID(c, r.RoleID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoleFindByID error", r.RoleID, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	if role.ID == 0 {
		c.JSON(bizresp.InvalidAssociatedID.Reply())
		return
	}
	ok, err := h.service.AdminUserCreate(c, &entity.AdminUser{
		Username: r.Username,
		Password: r.Password,
		RoleID:   r.RoleID,
		Status:   entity.StatusOn,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserCreate error", r.RoleID, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	if !ok {
		c.JSON(bizresp.UniqueKeyExist.WithMsg("用户名已存在"))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminUserPassword(c *gin.Context) {
	var r proto.AdminUserPasswordArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(bizresp.InvalidParam.Reply())
		return
	}
	user, err := h.service.AdminUserFindByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserFindByID error", r.ID, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	if user.ID == 0 {
		c.JSON(bizresp.InvalidAssociatedID.Reply())
		return
	}
	if user.IsSuper() {
		c.JSON(bizresp.OperationDeny.Reply())
		return
	}
	err = h.service.AdminUserUpdate(c, &entity.AdminUser{
		ID:       r.ID,
		Password: r.Password,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserUpdate error", &r, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminUserRole(c *gin.Context) {
	var r proto.AdminUserRoleArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(bizresp.InvalidParam.Reply())
		return
	}
	user, err := h.service.AdminUserFindByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserFindByID error", r.ID, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	if user.ID == 0 {
		c.JSON(bizresp.InvalidAssociatedID.Reply())
		return
	}
	if user.IsSuper() {
		c.JSON(bizresp.OperationDeny.Reply())
		return
	}
	if user.RoleID == r.RoleID {
		c.JSON(OK, Empty)
		return
	}
	role, err := h.service.AdminRoleFindByID(c, r.RoleID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoleFindByID error", r.RoleID, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	if role.ID == 0 {
		c.JSON(bizresp.InvalidAssociatedID.Reply())
		return
	}
	err = h.service.AdminUserUpdate(c, &entity.AdminUser{
		ID:     r.ID,
		RoleID: r.RoleID,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserUpdate error", &r, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

func (h *Handler) AdminUserStatus(c *gin.Context) {
	var r proto.SwitchStatusArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(bizresp.InvalidParam.Reply())
		return
	}
	user, err := h.service.AdminUserFindByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserFindByID error", r.ID, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	if user.ID == 0 {
		c.JSON(bizresp.InvalidAssociatedID.Reply())
		return
	}
	if user.IsSuper() {
		c.JSON(bizresp.OperationDeny.Reply())
		return
	}
	if user.Status == r.Status {
		c.JSON(OK, Empty)
		return
	}
	err = h.service.AdminUserUpdate(c, &entity.AdminUser{
		ID:     r.ID,
		Status: r.Status,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserUpdate error", &r, err)
		c.JSON(bizresp.WithErr(err))
		return
	}
	c.JSON(OK, Empty)
}

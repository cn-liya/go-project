package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"github.com/gin-gonic/gin"
	"project/cms/internal/acl"
	"project/cms/internal/proto"
	"project/pkg/logger"
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

type CaptchaResp struct {
	SessionKey  string `json:"session_key"`
	Base64Image []byte `json:"base64_image"`
}

func (h *Handler) Captcha(c *gin.Context) {
	code, bin := h.drawer.Generate(4)
	exp := strconv.FormatInt(time.Now().Unix()+70, 10)
	key := exp + "." + h.captchaSign(exp, code)
	c.JSON(OK, &CaptchaResp{
		SessionKey:  key,
		Base64Image: bin,
	})
}

type LoginArgs struct {
	SessionKey string `json:"session_key" binding:"required"`
	Captcha    string `json:"captcha" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
}
type LoginResp struct {
	Token     string        `json:"token"`
	Username  string        `json:"username"`
	Authority acl.Authority `json:"authority"`
}

func (h *Handler) UserLogin(c *gin.Context) {
	var r LoginArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWrapErr(err))
		return
	}

	sli := strings.Split(r.SessionKey, ".")
	if len(sli) != 2 {
		c.JSON(RespWithMsg(InvalidParam, "Invalid SessionKey"))
		return
	}
	exp, _ := strconv.ParseInt(sli[0], 10, 64)
	if time.Now().Unix() > exp {
		c.JSON(RespWithMsg(InvalidParam, "验证码过期"))
		return
	}
	sign := h.captchaSign(sli[0], r.Captcha)
	if sli[1] != sign {
		c.JSON(RespWithMsg(InvalidParam, "验证码错误"))
		return
	}

	admin, err := h.service.TakeAdminByUsername(c, r.Username)
	if err != nil {
		logger.FromContext(c).Error("service.TakeAdminByUsername error", r.Username, err)
		c.JSON(RespWrapErr(err))
		return
	}
	if admin.ID == 0 || !acl.CheckPassword(r.Password, admin.Password) {
		c.JSON(RespWithMsg(InvalidParam, "用户名或密码错误"))
		return
	}
	if admin.Status != 1 {
		c.JSON(RespWithMsg(Forbidden, "账号已禁用请联系管理员"))
		return
	}
	token, err := h.service.SetAdminToken(c, admin)
	if err != nil {
		logger.FromContext(c).Error("service.SetAdminToken error", admin, err)
		c.JSON(RespWrapErr(err))
		return
	}
	if admin.Username == acl.SuperUser {
		admin.Authority = acl.AllAuthority
	}
	c.JSON(OK, &LoginResp{
		Token:     token,
		Username:  admin.Username,
		Authority: admin.Authority,
	})
}

func (h *Handler) UserLogout(c *gin.Context) {
	err := h.service.DelAdminToken(c, c.GetHeader("Authorization"))
	if err != nil {
		logger.FromContext(c).Error("service.DelAdminToken error", nil, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, Empty)
}

type ModifyUserPasswordArgs struct {
	Password string `json:"password" binding:"required,min=6"`
}

func (h *Handler) ModifyUserPassword(c *gin.Context) {
	var r ModifyUserPasswordArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithMsg(InvalidParam, err.Error()))
		return
	}
	v, _ := c.Get("user")
	user := v.(*proto.AdminToken)
	err := h.service.UpdateAdmin(c, &acl.Admin{
		ID:       user.ID,
		Password: r.Password,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UpdateAdmin error", user.ID, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, Empty)
}

type AdminListArgs struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

type AdminListResp struct {
	Page  int          `json:"page"`
	Size  int          `json:"size"`
	Total int64        `json:"total"`
	List  []*AdminItem `json:"list"`
}
type AdminItem struct {
	ID         int           `json:"id"`
	Username   string        `json:"username"`
	Authority  acl.Authority `json:"authority"`
	Status     int8          `json:"status"`
	CreateTime string        `json:"create_time"`
	UpdateTime string        `json:"update_time"`
}

func (h *Handler) AdminList(c *gin.Context) {
	var r AdminListArgs
	_ = c.ShouldBindQuery(&r)
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Size < 10 || r.Size > 100 {
		r.Size = 10
	}
	total, list, err := h.service.PaginateAdmin(c, &proto.Pagination{
		Limit:  r.Size,
		Offset: (r.Page - 1) * r.Size,
	})
	if err != nil {
		logger.FromContext(c).Error("service.PaginateAdmin error", &r, err)
		c.JSON(RespWrapErr(err))
		return
	}
	items := make([]*AdminItem, 0, len(list))
	for _, v := range list {
		items = append(items, &AdminItem{
			ID:         v.ID,
			Username:   v.Username,
			Authority:  v.Authority,
			Status:     v.Status,
			CreateTime: v.CreateTime.Format(TimeFormat),
			UpdateTime: v.UpdateTime.Format(TimeFormat),
		})
	}
	c.JSON(OK, &AdminListResp{
		Page:  r.Page,
		Size:  r.Size,
		Total: total,
		List:  items,
	})
}

type CreateAdminArgs struct {
	Username string `json:"username" binding:"required,min=2,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

func (h *Handler) CreateAdmin(c *gin.Context) {
	var r CreateAdminArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWrapErr(err))
		return
	}
	if r.Username == acl.SuperUser {
		c.JSON(RespWithMsg(InvalidParam, "用户名不可用"))
		return
	}
	one, err := h.service.TakeAdminByUsername(c, r.Username)
	if err != nil {
		logger.FromContext(c).Error("service.TakeAdminByUsername error", r.Username, err)
		c.JSON(RespWrapErr(err))
		return
	}
	if one.ID > 0 {
		c.JSON(RespWithMsg(Conflict, "用户名已存在"))
		return
	}
	err = h.service.CreateAdmin(c, &acl.Admin{
		Username: r.Username,
		Password: r.Password,
		Status:   1,
	})
	if err != nil {
		logger.FromContext(c).Error("service.CreateAdmin error", &r, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, Empty)
}

type ResetAdminPasswordArgs struct {
	ID       int    `json:"id" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

func (h *Handler) ResetAdminPassword(c *gin.Context) {
	var r ResetAdminPasswordArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWrapErr(err))
		return
	}
	admin, err := h.service.FindAdminByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.FindAdminByID error", r.ID, err)
		c.JSON(RespWrapErr(err))
		return
	}
	if admin.ID == 0 || admin.Username == acl.SuperUser {
		c.JSON(NotFound, Empty)
		return
	}
	err = h.service.UpdateAdmin(c, &acl.Admin{
		ID:       r.ID,
		Password: r.Password,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UpdateAdmin error", r.ID, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, Empty)
}

type AssignAdminAuthorityArgs struct {
	ID        int          `json:"id" binding:"required,min=1"`
	Authority []*Authority `json:"authority" binding:"required,dive"`
}
type Authority struct {
	Module string `json:"module" binding:"required"`
	RW     int8   `json:"rw" binding:"oneof=0 1 2"`
}

func (h *Handler) AssignAdminAuthority(c *gin.Context) {
	var r AssignAdminAuthorityArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWrapErr(err))
		return
	}
	admin, err := h.service.FindAdminByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.FindAdminByID error", r.ID, err)
		c.JSON(RespWrapErr(err))
		return
	}
	if admin.ID == 0 || admin.Username == acl.SuperUser {
		c.JSON(NotFound, Empty)
		return
	}
	auth := make(acl.Authority)
	for _, v := range r.Authority {
		if _, ok := acl.AllAuthority[v.Module]; ok && v.RW > 0 {
			auth[v.Module] = v.RW
		}
	}
	err = h.service.UpdateAdmin(c, &acl.Admin{
		ID:        r.ID,
		Authority: auth,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UpdateAdmin error", &r, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, Empty)
}

type SwitchAdminStatusArgs struct {
	ID     int  `json:"id" binding:"required,min=1"`
	Status int8 `json:"status" binding:"required,eq=-1|eq=1"`
}

func (h *Handler) SwitchAdminStatus(c *gin.Context) {
	var r SwitchAdminStatusArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWrapErr(err))
		return
	}
	admin, err := h.service.FindAdminByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.FindAdminByID error", r.ID, err)
		c.JSON(RespWrapErr(err))
		return
	}
	if admin.ID == 0 || admin.Username == acl.SuperUser {
		c.JSON(NotFound, Empty)
		return
	}
	if r.Status == admin.Status {
		c.JSON(OK, Empty)
		return
	}
	err = h.service.UpdateAdmin(c, &acl.Admin{
		ID:     r.ID,
		Status: r.Status,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UpdateAdmin error", &r, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, Empty)
}

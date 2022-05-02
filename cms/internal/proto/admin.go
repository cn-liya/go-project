package proto

type CaptchaResp struct {
	SessionKey  string `json:"session_key"`
	Base64Image []byte `json:"base64_image"`
}

type LoginArgs struct {
	SessionKey string `json:"session_key" binding:"required"`
	Captcha    string `json:"captcha" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token    string   `json:"token"`
	Username string   `json:"username"`
	PwdTip   string   `json:"pwd_tip,omitempty"`
	IsSuper  bool     `json:"is_super,omitempty"`
	Actions  []string `json:"actions,omitempty"`
	Menus    []string `json:"menus,omitempty"`
}

type UserPasswordArgs struct {
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type AdminRoleListResp struct {
	Total int64            `json:"total"`
	List  []*AdminRoleItem `json:"list"`
}

type AdminRoleItem struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Actions []string `json:"actions"`
	Menus   []string `json:"menus"`
}

type AdminRoleCreateArgs struct {
	Name    string   `json:"name" binding:"min=2,max=32"`
	Actions []string `json:"actions"`
	Menus   []string `json:"menus"`
}

type AdminRoleUpdateArgs struct {
	ID int `json:"id" binding:"min=1"`
	AdminRoleCreateArgs
}

type AdminUserListArgs struct {
	ListArgs
	RoleID   int    `json:"role_id"`
	Username string `json:"username" binding:"max=32"`
	Status   int8   `json:"status" binding:"min=-1,max=1"`
}

type AdminUserListResp struct {
	Total int64            `json:"total"`
	List  []*AdminUserItem `json:"list"`
}

type AdminUserItem struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	RoleID   int    `json:"role_id"`
	IsSuper  bool   `json:"is_super,omitempty"`
	Status   int8   `json:"status"`
}

type AdminUserCreateArgs struct {
	Username string `json:"username" binding:"min=6,max=32"`
	Password string `json:"password" binding:"min=6,max=32"`
	RoleID   int    `json:"role_id" binding:"min=1"`
}

type AdminUserPasswordArgs struct {
	ID       int    `json:"id" binding:"min=1"`
	Password string `json:"password" binding:"min=6,max=32"`
}

type AdminUserRoleArgs struct {
	ID     int `json:"id" binding:"min=1"`
	RoleID int `json:"role_id" binding:"min=1"`
}

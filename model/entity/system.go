package entity

const (
	SystemActionGroup = 1
	SystemActionRoute = 2
)

type SystemAction struct {
	KeyName string `json:"key_name" gorm:"primaryKey"`
	Level   int8   `json:"-"`
	Title   string `json:"title"`
	Sort    int    `json:"sort"`
}

func (*SystemAction) TableName() string { return "system_action" }

type SystemConfig struct {
	KeyName string `json:"key_name" gorm:"primaryKey"`
	Content string `json:"content"`
}

func (*SystemConfig) TableName() string { return "system_config" }

const (
	SysCfgMenuKeys = "sys_menu_keys"
	SysCfgMenuTree = "sys_menu_trees"
)

type SysMenuTree struct {
	Name  string         `json:"name"`
	Title string         `json:"title"`
	Sub   []*SysMenuTree `json:"sub,omitempty"`
}

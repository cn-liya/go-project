package acl

const (
	ModuleAdmin    = "admin"
	ModuleAnalysis = "analysis"
)

type Module struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

var Modules = []*Module{
	{Key: ModuleAdmin, Name: "账号权限"},
	{Key: ModuleAnalysis, Name: "趋势分析"},
}

var AllAuthority = make(Authority)

func init() {
	for _, v := range Modules {
		AllAuthority[v.Key] = 2
	}
}

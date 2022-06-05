package model

const (
	StatusOff int8 = -1
	StatusOn  int8 = 1
)

const (
	DefaultCity = "100000"
)

var Cities = map[string]string{
	"110000": "北京",
	"310000": "上海",
	"330100": "杭州",
	"440100": "广州",
	"440300": "深圳",
}

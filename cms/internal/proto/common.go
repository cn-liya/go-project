package proto

type IdArgs struct {
	ID int `json:"id" binding:"min=1"`
}

type KeyArgs struct {
	KeyName string `json:"key_name" binding:"required"`
}

type ListArgs struct {
	Page int `json:"page" binding:"min=1"`
	Size int `json:"size" binding:"min=10,max=100"`
}

type SwitchStatusArgs struct {
	ID     int  `json:"id" binding:"min=1"`
	Status int8 `json:"status" binding:"eq=-1|eq=1"`
}

type OptionItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type OptionResp struct {
	List []*OptionItem `json:"list"`
}

type UploadResp struct {
	Host string `json:"host"`
	Path string `json:"path"`
}

type CryptoArgs struct {
	Text string `json:"text"`
}

type CryptoResp struct {
	Result string `json:"result"`
}

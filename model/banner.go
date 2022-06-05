package model

type Banner struct {
	ID        int    `json:"id"`
	City      string `json:"city"`
	Title     string `json:"title"`
	Img       string `json:"img"`
	Type      int8   `json:"type"`
	Link      string `json:"link"`
	Sort      int8   `json:"sort"`
	BeginTime int64  `json:"begin_time"`
	EndTime   int64  `json:"end_time"`
	Status    int8   `json:"status"`
}

func (*Banner) TableName() string {
	return "banner"
}

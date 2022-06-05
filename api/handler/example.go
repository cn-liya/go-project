package handler

import (
	"github.com/gin-gonic/gin"
	"project/pkg/logger"
	"strings"
	"time"
)

type BannerItem struct {
	Title string `json:"title"`
	Img   string `json:"img"`
	Type  int8   `json:"type"`
	Link  string `json:"link"`
}
type BannersResp struct {
	List []*BannerItem `json:"list"`
}

func (h *Handler) GetBanners(c *gin.Context) {
	city := c.Query("city")
	data, err := h.service.GetCityBanners(c, city)
	if err != nil {
		logger.FromContext(c).Error("service.GetCityBanners error", nil, err)
		c.JSON(RespWrapErr(err))
		return
	}
	list := make([]*BannerItem, 0, len(data))
	now := time.Now().Unix()
	for _, d := range data {
		if d.BeginTime <= now && now < d.EndTime {
			if !strings.HasPrefix(d.Img, "http") { //相对路径拼上cdn域名
				d.Img = h.cdn + d.Img
			}
			list = append(list, &BannerItem{
				Title: d.Title,
				Img:   d.Img,
				Type:  d.Type,
				Link:  d.Link,
			})
		}
	}
	c.JSON(OK, &BannersResp{List: list})
}

func (h *Handler) PushMessage(c *gin.Context) {
	//err := h.service.PushMessage(c, &model.MsgExample{
	//	UUID:   hex.EncodeToString(uuid.NewV4().Bytes()),
	//	Number: rand.Intn(100),
	//})
	//if err != nil {
	//	logger.FromContext(c).Error("service.PushMessage error", nil, err)
	//	c.JSON(RespWrapErr(err))
	//	return
	//}
	c.JSON(OK, Empty)
}

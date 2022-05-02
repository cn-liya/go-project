package wechatwork

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

/*
文档：https://open.work.weixin.qq.com/api/doc/90000/90136/91770
*/

type Text struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type Markdown struct {
	Content string `json:"content"`
}

type Image struct {
	Base64 string `json:"base64"`
	MD5    string `json:"md5"`
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url"`
	Picurl      string `json:"picurl,omitempty"`
}
type News struct {
	Articles []*Article `json:"articles"`
}

type File struct {
	MediaId string `json:"media_id"`
}

type Source struct {
	IconUrl string `json:"icon_url,omitempty"`
	Desc    string `json:"desc,omitempty"`
}
type TitleDesc struct {
	Title string `json:"title"`
	Desc  string `json:"desc,omitempty"`
}
type CardImage struct {
	Url         string  `json:"url"`
	AspectRatio float64 `json:"aspect_ratio,omitempty"`
}
type HorizontalContentList struct {
	Keyname string `json:"keyname"`
	Value   string `json:"value,omitempty"`
	Type    int    `json:"type,omitempty"` //0|1|2
	Url     string `json:"url,omitempty"`
	MediaId string `json:"media_id,omitempty"`
}
type JumpList struct {
	Title    string `json:"title"`
	Type     int    `json:"type,omitempty"` //0|1|2
	Url      string `json:"url,omitempty"`
	Appid    string `json:"appid,omitempty"`
	Pagepath string `json:"pagepath,omitempty"`
}
type CardAction struct {
	Type     int    `json:"type,omitempty"` //0|1|2
	Url      string `json:"url,omitempty"`
	Appid    string `json:"appid,omitempty"`
	Pagepath string `json:"pagepath,omitempty"`
}
type TemplateCard struct {
	CardType              string                   `json:"card_type"` //text_notice|template_card
	Source                *Source                  `json:"source,omitempty"`
	MainTitle             *TitleDesc               `json:"main_title"`
	EmphasisContent       *TitleDesc               `json:"emphasis_content,omitempty"`
	SubTitleText          string                   `json:"sub_title_text,omitempty"`
	CardImage             *CardImage               `json:"card_image,omitempty"`
	VerticalContentList   []*TitleDesc             `json:"vertical_content_list,omitempty"`
	HorizontalContentList []*HorizontalContentList `json:"horizontal_content_list,omitempty"`
	JumpList              []*JumpList              `json:"jump_list,omitempty"`
	CardAction            *CardAction              `json:"card_action"`
}

type Args struct {
	Msgtype      string        `json:"msgtype"`
	Text         *Text         `json:"text,omitempty"`
	Markdown     *Markdown     `json:"markdown,omitempty"`
	Image        *Image        `json:"image,omitempty"`
	News         *News         `json:"news,omitempty"`
	File         *File         `json:"file,omitempty"`
	TemplateCard *TemplateCard `json:"template_card,omitempty"`
}

type Resp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func send(webhookUrl string, args *Args) (*Resp, error) {
	b, _ := json.Marshal(args)
	resp, err := http.Post(webhookUrl, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var result Resp
	err = json.Unmarshal(data, &result)
	return &result, err
}

func SendText(webhookUrl string, text *Text) (*Resp, error) {
	return send(webhookUrl, &Args{
		Msgtype: "text",
		Text:    text,
	})
}

func SendMarkdown(webhookUrl string, content string) (*Resp, error) {
	return send(webhookUrl, &Args{
		Msgtype:  "markdown",
		Markdown: &Markdown{Content: content},
	})
}

func SendImage(webhookUrl string, bin []byte) (*Resp, error) {
	sum := md5.Sum(bin)
	return send(webhookUrl, &Args{
		Msgtype: "image",
		Image: &Image{
			Base64: base64.StdEncoding.EncodeToString(bin),
			MD5:    hex.EncodeToString(sum[:]),
		},
	})
}

func SendNews(webhookUrl string, articles []*Article) (*Resp, error) {
	return send(webhookUrl, &Args{
		Msgtype: "news",
		News:    &News{Articles: articles},
	})
}

func SendFile(webhookUrl string, mediaId string) (*Resp, error) {
	return send(webhookUrl, &Args{
		Msgtype: "file",
		File:    &File{MediaId: mediaId},
	})
}

func SendCard(webhookUrl string, card *TemplateCard) (*Resp, error) {
	return send(webhookUrl, &Args{
		Msgtype:      "template_card",
		TemplateCard: card,
	})
}

type UploadFileResp struct {
	Resp
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

func UploadFile(webhookUrl string, bin []byte, filename string) (*UploadFileResp, error) {
	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)
	form, _ := w.CreateFormFile("media", filename)
	form.Write(bin) // nolint
	w.Close()       // nolint
	webhook, _ := url.Parse(webhookUrl)
	resp, err := http.Post(
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?type=file&key="+webhook.Query().Get("key"),
		w.FormDataContentType(), buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var result UploadFileResp
	err = json.Unmarshal(data, &result)
	return &result, err
}

package dingtalk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

/*
文档：https://developers.dingtalk.com/document/robots/custom-robot-access
*/

type At struct {
	Mobiles []string `json:"atMobiles,omitempty"`
	UserIds []string `json:"atUserIds,omitempty"`
	All     bool     `json:"isAtAll"`
}
type Text struct {
	Content string `json:"content"`
}
type Link struct {
	Text       string `json:"text"`
	Title      string `json:"title"`
	PicURL     string `json:"picUrl"`
	MessageURL string `json:"messageUrl"`
}
type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
type ActionCardBtn struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}
type ActionCard struct {
	Title          string           `json:"title"`
	Text           string           `json:"text"`
	BtnOrientation string           `json:"btnOrientation,omitempty"` // 0|1
	SingleTitle    string           `json:"singleTitle,omitempty"`
	SingleURL      string           `json:"singleURL,omitempty"`
	Btns           []*ActionCardBtn `json:"btns,omitempty"`
}
type FeedCardLink struct {
	Title      string `json:"title"`
	MessageURL string `json:"messageURL"`
	PicURL     string `json:"picURL"`
}
type FeedCard struct {
	Links []*FeedCardLink `json:"links"`
}
type Args struct {
	Msgtype    string      `json:"msgtype"`
	At         *At         `json:"at,omitempty"`
	Text       *Text       `json:"text,omitempty"`
	Link       *Link       `json:"link,omitempty"`
	Markdown   *Markdown   `json:"markdown,omitempty"`
	ActionCard *ActionCard `json:"actionCard,omitempty"`
	FeedCard   *FeedCard   `json:"feedCard,omitempty"`
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

func SendText(webhookUrl string, text *Text, at *At) (*Resp, error) {
	return send(webhookUrl, &Args{
		Msgtype: "text",
		At:      at,
		Text:    text,
	})
}

func SendMarkdown(webhookUrl string, markdown *Markdown, at *At) (*Resp, error) {
	return send(webhookUrl, &Args{
		Msgtype:  "markdown",
		At:       at,
		Markdown: markdown,
	})
}

func SendActionCard(webhookUrl string, card *ActionCard) (*Resp, error) {
	return send(webhookUrl, &Args{
		Msgtype:    "actionCard",
		ActionCard: card,
	})
}

func SendFeedCard(webhookUrl string, card *FeedCard) (*Resp, error) {
	return send(webhookUrl, &Args{
		Msgtype:  "feedCard",
		FeedCard: card,
	})
}

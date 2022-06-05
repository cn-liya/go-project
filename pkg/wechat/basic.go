package wechat

import (
	"context"
	"net/url"
)

type JsCode2SessionResp struct {
	respErr
	Openid     string `json:"openid,omitempty"`
	Unionid    string `json:"unionid,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
}

func (api *basic) JsCode2Session(ctx context.Context, code string) (*JsCode2SessionResp, error) {
	data := make(url.Values)
	data.Set("appid", api.appid)
	data.Set("secret", api.secret)
	data.Set("js_code", code)
	data.Set("grant_type", "authorization_code")
	var resp JsCode2SessionResp
	err := api.get(ctx, "/sns/jscode2session", data, &resp)
	return &resp, err
}

type GetAccessTokenResp struct {
	respErr
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}

func (api *basic) GetAccessToken(ctx context.Context) (*GetAccessTokenResp, error) {
	data := make(url.Values)
	data.Set("appid", api.appid)
	data.Set("secret", api.secret)
	data.Set("grant_type", "client_credential")
	var resp GetAccessTokenResp
	err := api.get(ctx, "/cgi-bin/token", data, &resp)
	return &resp, err
}

package wechat

import (
	"context"
)

type GetUserPhoneNumberResp struct {
	respErr
	PhoneInfo *PhoneInfo `json:"phone_info,omitempty"`
}
type PhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}

func (api *server) GetUserPhoneNumber(ctx context.Context, code string) (*GetUserPhoneNumberResp, error) {
	var resp GetUserPhoneNumberResp
	err := api.post(ctx, "/wxa/business/getuserphonenumber",
		map[string]string{"code": code}, &resp)
	return &resp, err
}

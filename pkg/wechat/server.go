package wechat

import (
	"context"
)

type UserPhoneNumberResp struct {
	respErr
	PhoneInfo *PhoneInfo `json:"phone_info,omitempty"`
}
type PhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}

func (api *server) GetUserPhoneNumber(ctx context.Context, code string) (*UserPhoneNumberResp, error) {
	var resp UserPhoneNumberResp
	err := api.post(ctx, "/wxa/business/getuserphonenumber",
		map[string]string{"code": code}, &resp)
	return &resp, err
}

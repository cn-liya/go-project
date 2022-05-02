package coss

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// COS (Tencent Cloud Object Storage)
type COS interface {
	PutObject(ctx context.Context, path string, reader io.Reader) error
	GetSignURL(ctx context.Context, path string, expired time.Duration) (string, error)
}

type tcos struct {
	client *cos.Client
}

type CosConfig struct {
	BucketURL  string
	ServiceURL string
	SecretID   string
	SecretKey  string
}

func NewTCOS(cfg *CosConfig) COS {
	bu, _ := url.Parse(cfg.BucketURL)
	su, _ := url.Parse(cfg.ServiceURL)
	baseURL := &cos.BaseURL{
		BucketURL:  bu,
		ServiceURL: su,
	}
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})
	return &tcos{client}
}

func (s *tcos) PutObject(ctx context.Context, path string, reader io.Reader) error {
	_, err := s.client.Object.Put(ctx, strings.TrimLeft(path, "/"), reader, nil)
	return err
}

func (s *tcos) GetSignURL(ctx context.Context, path string, expired time.Duration) (string, error) {
	auth := s.client.GetCredential()
	u, err := s.client.Object.GetPresignedURL(ctx, http.MethodGet, strings.TrimLeft(path, "/"),
		auth.SecretID, auth.SecretKey, expired, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

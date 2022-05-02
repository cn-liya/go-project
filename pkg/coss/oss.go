package coss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"strings"
)

// OSS (Aliyun Object Storage Service)
type OSS interface {
	PutObject(path string, reader io.Reader) error
	GetSignURL(path string, expireSeconds int64) (string, error)
}

type alioss struct {
	bucket *oss.Bucket
}

type OssConfig struct {
	Endpoint   string
	KeyID      string
	KeySecret  string
	BucketName string
}

func NewAliOSS(cfg *OssConfig) OSS {
	client, err := oss.New(cfg.Endpoint, cfg.KeyID, cfg.KeySecret)
	if err != nil {
		log.Fatal(err)
	}
	bucket, err := client.Bucket(cfg.BucketName)
	if err != nil {
		log.Fatal(err)
	}
	return &alioss{bucket}
}

func (s *alioss) PutObject(path string, reader io.Reader) error {
	return s.bucket.PutObject(strings.TrimLeft(path, "/"), reader)
}

func (s *alioss) GetSignURL(path string, expireSeconds int64) (string, error) {
	return s.bucket.SignURL(strings.TrimLeft(path, "/"), oss.HTTPGet, expireSeconds)
}

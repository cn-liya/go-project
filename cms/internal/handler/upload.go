package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"project/cms/internal/proto"
	"project/pkg/files"
	"project/pkg/logger"
	"project/pkg/web/errcode"
)

func (h *Handler) UploadImage(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(errcode.InvalidParam.Response())
		return
	}
	if f.Size > 500<<10 {
		c.JSON(errcode.EntityTooLarge.WithMsg("单张图片限制500KB以内"))
		return
	}
	file, _ := f.Open()
	defer file.Close()
	b, _ := io.ReadAll(file)
	ext := files.CheckImage(b)
	if ext == "" {
		c.JSON(errcode.UnsupportedMediaType.WithMsg("无效的图片类型，仅支持jpg/png/gif格式"))
		return
	}
	remotePath := "/img/" + files.GenFileName(b) + "." + ext
	err = h.cos.PutObject(c, remotePath, bytes.NewReader(b))
	//err = h.oss.PutObject(remotePath, bytes.NewReader(b))
	if err != nil {
		logger.FromContext(c).Error("cos.PutObject error", nil, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, &proto.UploadResp{
		Host: h.cdn,
		Path: remotePath,
	})
}

func (h *Handler) UploadAudio(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(errcode.InvalidParam.Response())
		return
	}
	if f.Size > 2<<20 {
		c.JSON(errcode.EntityTooLarge.WithMsg("单个文件限制2M以内"))
		return
	}
	file, _ := f.Open()
	defer file.Close()
	b, _ := io.ReadAll(file)
	ext := files.CheckAudio(b)
	if ext == "" {
		c.JSON(errcode.UnsupportedMediaType.WithMsg("无效的音频类型，仅支持mp3/wav/oga格式"))
		return
	}
	remotePath := "/audio/" + files.GenFileName(b) + "." + ext
	err = h.cos.PutObject(c, remotePath, bytes.NewReader(b))
	//err = h.oss.PutObject(remotePath, bytes.NewReader(b))
	if err != nil {
		logger.FromContext(c).Error("cos.PutObject error", nil, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, &proto.UploadResp{
		Host: h.cdn,
		Path: remotePath,
	})
}

func (h *Handler) UploadVideo(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(errcode.InvalidParam.Response())
		return
	}
	if f.Size > 5<<20 {
		c.JSON(errcode.EntityTooLarge.WithMsg("单个文件限制5M以内"))
		return
	}
	file, _ := f.Open()
	defer file.Close()
	b, _ := io.ReadAll(file)
	ext := files.CheckVideo(b)
	if ext == "" {
		c.JSON(errcode.UnsupportedMediaType.WithMsg("无效的视频类型，仅支持mp4/webm/ogv格式"))
		return
	}
	remotePath := "/video/" + files.GenFileName(b) + "." + ext
	err = h.cos.PutObject(c, remotePath, bytes.NewReader(b))
	//err = h.oss.PutObject(remotePath, bytes.NewReader(b))
	if err != nil {
		logger.FromContext(c).Error("cos.PutObject error", nil, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, &proto.UploadResp{
		Host: h.cdn,
		Path: remotePath,
	})
}

func (h *Handler) UploadPDF(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(errcode.InvalidParam.Response())
		return
	}
	if f.Size > 1<<20 {
		c.JSON(errcode.EntityTooLarge.WithMsg("单个文件限制1M以内"))
		return
	}
	file, _ := f.Open()
	defer file.Close()
	b, _ := io.ReadAll(file)
	if http.DetectContentType(b) != "application/pdf" {
		c.JSON(errcode.UnsupportedMediaType.WithMsg("无效的文件类型，仅支持pdf格式"))
		return
	}
	remotePath := "/doc/" + files.GenFileName(b) + ".pdf"
	err = h.cos.PutObject(c, remotePath, bytes.NewReader(b))
	//err = h.oss.PutObject(remotePath, bytes.NewReader(b))
	if err != nil {
		logger.FromContext(c).Error("cos.PutObject error", nil, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, &proto.UploadResp{
		Host: h.cdn,
		Path: remotePath,
	})
}

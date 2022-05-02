package ades

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Cipher struct {
	block cipher.Block
	iv    []byte
}

//func NewDesCipher(key, iv []byte) (*Cipher, error) {
//	return newCipher(des.NewCipher, key, iv)
//}
//
//func NewTripleDesCipher(key, iv []byte) (*Cipher, error) {
//	return newCipher(des.NewTripleDESCipher, key, iv)
//}

func NewAesCipher(key, iv []byte) (*Cipher, error) {
	return newCipher(aes.NewCipher, key, iv)
}

func newCipher(f func([]byte) (cipher.Block, error), key, iv []byte) (*Cipher, error) {
	block, err := f(key)
	if err != nil {
		return nil, err
	}
	ivb := make([]byte, block.BlockSize())
	copy(ivb, iv)
	return &Cipher{
		block: block,
		iv:    ivb,
	}, nil
}

//func pkcs7Padding(src []byte, blockSize int) []byte {
//	length := len(src)
//	pad := blockSize - length%blockSize
//	padText := bytes.Repeat([]byte{byte(pad)}, pad)
//	res := make([]byte, length, length+pad)
//	copy(res, src)
//	return append(res, padText...)
//}
//
//var ErrInvalidPadLen = errors.New("invalid padding length")
//
//// 三种填充方式去除填充的方法都相同，根据最后一位来精确定位截取长度
//func trimPadding(dst []byte) ([]byte, error) {
//	length := len(dst)
//	if length == 0 {
//		return dst, nil
//	}
//	pad := int(dst[length-1])
//	if length < pad {
//		return nil, ErrInvalidPadLen
//	}
//	return dst[:length-pad], nil
//}
//
//func (c *Cipher) EncryptECB(src []byte) string {
//	bs := c.block.BlockSize()
//	src = pkcs7Padding(src, bs)
//	length := len(src)
//	dst := make([]byte, length)
//	for i := 0; i < length; i += bs {
//		c.block.Encrypt(dst[i:i+bs], src[i:i+bs])
//	}
//	return base64.StdEncoding.EncodeToString(dst)
//}
//func (c *Cipher) DecryptECB(s string) ([]byte, error) {
//	src, err := base64.StdEncoding.DecodeString(s)
//	if err != nil {
//		return nil, err
//	}
//	bs := c.block.BlockSize()
//	length := len(src)
//	if length%bs != 0 {
//		return nil, ErrInvalidPadLen
//	}
//	dst := make([]byte, length)
//	for i := 0; i < length; i += bs {
//		c.block.Decrypt(dst[i:i+bs], src[i:i+bs])
//	}
//	return trimPadding(dst)
//}
//
//func (c *Cipher) EncryptCBC(src []byte) string {
//	src = pkcs7Padding(src, c.block.BlockSize())
//	dst := make([]byte, len(src))
//	bm := cipher.NewCBCEncrypter(c.block, c.iv)
//	bm.CryptBlocks(dst, src)
//	return base64.StdEncoding.EncodeToString(dst)
//}
//func (c *Cipher) DecryptCBC(s string) ([]byte, error) {
//	src, err := base64.StdEncoding.DecodeString(s)
//	if err != nil {
//		return nil, err
//	}
//	dst := make([]byte, len(src))
//	bm := cipher.NewCBCDecrypter(c.block, c.iv)
//	bm.CryptBlocks(dst, src)
//	return trimPadding(dst)
//}
//
//func (c *Cipher) EncryptCFB(src []byte) string {
//	dst := make([]byte, len(src))
//	stream := cipher.NewCFBEncrypter(c.block, c.iv)
//	stream.XORKeyStream(dst, src)
//	return base64.StdEncoding.EncodeToString(dst)
//}
//func (c *Cipher) DecryptCFB(s string) ([]byte, error) {
//	src, err := base64.StdEncoding.DecodeString(s)
//	if err != nil {
//		return nil, err
//	}
//	dst := make([]byte, len(src))
//	stream := cipher.NewCFBDecrypter(c.block, c.iv)
//	stream.XORKeyStream(dst, src)
//	return dst, nil
//}
//
//func (c *Cipher) EncryptOFB(src []byte) string {
//	dst := make([]byte, len(src))
//	stream := cipher.NewOFB(c.block, c.iv)
//	stream.XORKeyStream(dst, src)
//	return base64.StdEncoding.EncodeToString(dst)
//}
//func (c *Cipher) DecryptOFB(s string) ([]byte, error) {
//	src, err := base64.StdEncoding.DecodeString(s)
//	if err != nil {
//		return nil, err
//	}
//	dst := make([]byte, len(src))
//	stream := cipher.NewOFB(c.block, c.iv)
//	stream.XORKeyStream(dst, src)
//	return dst, nil
//}

func (c *Cipher) EncryptCTR(src []byte) string {
	dst := make([]byte, len(src))
	stream := cipher.NewCTR(c.block, c.iv)
	stream.XORKeyStream(dst, src)
	return base64.StdEncoding.EncodeToString(dst)
}
func (c *Cipher) DecryptCTR(s string) ([]byte, error) {
	src, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(src))
	stream := cipher.NewCTR(c.block, c.iv)
	stream.XORKeyStream(dst, src)
	return dst, nil
}

package logger

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"
)

var (
	once   = &sync.Once{}
	sep    = regexp.MustCompile(`\\*"`)
	secret []*regexp.Regexp
	block  cipher.Block
)

type Config struct {
	Topic        string
	Output       string
	CipherKey    string
	CipherFields []string
}

func Setup(cfg *Config) {
	once.Do(func() {
		if len(cfg.CipherFields) > 0 && cfg.CipherKey != "" {
			b, err := aes.NewCipher([]byte(cfg.CipherKey))
			if err != nil {
				log.Fatal(err)
			}
			block = b
			valid := regexp.MustCompile(`^\w+$`) //字段名仅可包含字母数字下划线
			for _, f := range cfg.CipherFields {
				if !valid.MatchString(f) {
					log.Fatal("invalid log.cipher.field: ", f)
				}
				secret = append(secret, regexp.MustCompile(`(?i)\\*"\w*`+f+`\w*\\*"\s*:\s*\\*"(.*?)\\*"`))
			}
		}
		switch cfg.Output {
		case "file":
			setLogToFile(cfg.Topic)
		case "std":
			setLogToStdout()
		case "fmt":
			setLogToFormat()
		}
	})
}

type field struct{ left, right int }

func matchReplace(src []byte) []byte {
	match := make(map[field]struct{})
	for _, reg := range secret {
		allIndex := reg.FindAllIndex(src, -1)
		count := len(allIndex)
		if count == 0 {
			continue
		}
		for _, i := range allIndex {
			part := string(src[i[0]:i[1]])
			arr := sep.Split(part, -1)
			txtLen := len(arr[3])
			if txtLen == 0 {
				continue
			}
			l := strings.LastIndex(part, arr[3]) + i[0]
			r := l + txtLen
			match[field{l, r}] = struct{}{}
		}
	}
	count := len(match)
	if count == 0 {
		return src
	}
	items := make([]field, 0, count)
	for f := range match {
		items = append(items, f)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].left < items[j].left
	})
	buf := bytes.NewBuffer(nil)
	pre := 0
	for _, v := range items {
		buf.Write(src[pre:v.left])
		buf.WriteString(encrypt(src[v.left:v.right]))
		pre = v.right
	}
	buf.Write(src[pre:])
	return buf.Bytes()
}

func encrypt(src []byte) string {
	bs := block.BlockSize()
	length := len(src)
	pad := bs - length%bs
	padText := bytes.Repeat([]byte{byte(pad)}, pad)
	plain := make([]byte, length, length+pad)
	copy(plain, src)
	plain = append(plain, padText...)
	length = len(plain)
	dst := make([]byte, length)
	for i := 0; i < length; i += bs {
		block.Encrypt(dst[i:i+bs], plain[i:i+bs])
	}
	return base64.StdEncoding.EncodeToString(dst)
}

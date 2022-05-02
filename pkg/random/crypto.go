package random

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

// 随机字节数组
func Bytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

// V4版本不带-号的UUID
func UUID() string {
	return hex.EncodeToString(Bytes(16))
}

// 安全的六位数字码
func NumCode() string {
	num, _ := rand.Int(rand.Reader, big.NewInt(1e6))
	return fmt.Sprintf("%06d", num.Int64())
}

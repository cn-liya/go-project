package ades

import (
	"project/pkg/random"
	"testing"
)

func FuzzNewAesCipher(f *testing.F) {
	key := random.Bytes(16)
	iv := random.Bytes(16)
	c, err := NewAesCipher(key, iv)
	if err != nil {
		f.Fatal(err)
	}
	for i := 5; i < 15; i++ {
		f.Add(random.Strings(i))
	}
	f.Fuzz(func(t *testing.T, src string) {
		encrypted := c.EncryptCTR([]byte(src))
		decrypted, _ := c.DecryptCTR(encrypted)
		if string(decrypted) != src {
			t.Errorf("src:%s, encrypted:%s, decrypted:%s\n", src, encrypted, decrypted)
		}
	})
}

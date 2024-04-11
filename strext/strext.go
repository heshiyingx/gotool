package strext

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandStr 生成随机的字符串作为appKey
func GenerateRandStr() (string, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

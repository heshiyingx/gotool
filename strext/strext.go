package strext

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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

func ToJsonStr(obj any) string {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return "Err:" + err.Error()
	}
	return string(objBytes)

}

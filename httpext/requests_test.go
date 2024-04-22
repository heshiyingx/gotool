package httpext

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {
	url := ""
	header := map[string]string{"content-type": "application/json",
		"Authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJwaHBlcjY2Ni9qd3QiLCJ1aWQiOjYzNTUyNzY4MTI2OTYwODQ0OSwicGhvbmUiOiIxNzM3ODUxNjMyNSIsInN0YXR1cyI6MSwidGltZSI6MTcxMzc1MjY4OCwiand0X3NjZW5lIjoiYXBwbGljYXRpb24iLCJqdGkiOiJhcHBsaWNhdGlvbl82MzU1Mjc2ODEyNjk2MDg0NDkiLCJpYXQiOjE3MTM3NTI2ODgsIm5iZiI6MTcxMzc1MjY4OCwiZXhwIjoxNzEzODM5MDg4fQ.YY4-SXRCVt8-Qe82KI7o8p4-OmgMN4DtAJzRttLJOsYnedgUqAd8g2gopBz89mM1s6WQFewG6Ek9VTE8UprbGzlwB3W31KGd5SFP5Pfj78rueM8u4BV3oW4QovUwcvHTUzBqoLLkl-qqkhxSs4dcHbj_RomseTrajnOS1Eb2wRg",
	}
	reponse, err := Get[map[string]interface{}](context.Background(), url, nil, header)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(reponse)
}
func TestPOST(t *testing.T) {
	url := ""
	header := map[string]string{"content-type": "application/json",
		"operationID":   "xxx",
		"Authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJwaHBlcjY2Ni9qd3QiLCJ1aWQiOjYzNTUyNzY4MTI2OTYwODQ0OSwicGhvbmUiOiIxNzM3ODUxNjMyNSIsInN0YXR1cyI6MSwidGltZSI6MTcxMzc1MjY4OCwiand0X3NjZW5lIjoiYXBwbGljYXRpb24iLCJqdGkiOiJhcHBsaWNhdGlvbl82MzU1Mjc2ODEyNjk2MDg0NDkiLCJpYXQiOjE3MTM3NTI2ODgsIm5iZiI6MTcxMzc1MjY4OCwiZXhwIjoxNzEzODM5MDg4fQ.YY4-SXRCVt8-Qe82KI7o8p4-OmgMN4DtAJzRttLJOsYnedgUqAd8g2gopBz89mM1s6WQFewG6Ek9VTE8UprbGzlwB3W31KGd5SFP5Pfj78rueM8u4BV3oW4QovUwcvHTUzBqoLLkl-qqkhxSs4dcHbj_RomseTrajnOS1Eb2wRg",
	}
	param := map[string]interface{}{
		"areaCode":    "+86",
		"password":    "e9bc0e13a8a16cbb07b175d92a113126",
		"phoneNumber": "",
		"platform":    5,
	}
	reponse, err := POST[map[string]interface{}](context.Background(), url, param, header)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(reponse)
}

func BenchmarkGet(b *testing.B) {
	//var response = make(map[string]interface{},0)
	url := ""
	header := map[string]string{"content-type": "application/json",
		"Authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJwaHBlcjY2Ni9qd3QiLCJ1aWQiOjYzNTUyNzY4MTI2OTYwODQ0OSwicGhvbmUiOiIxNzM3ODUxNjMyNSIsInN0YXR1cyI6MSwidGltZSI6MTcxMzc1MjY4OCwiand0X3NjZW5lIjoiYXBwbGljYXRpb24iLCJqdGkiOiJhcHBsaWNhdGlvbl82MzU1Mjc2ODEyNjk2MDg0NDkiLCJpYXQiOjE3MTM3NTI2ODgsIm5iZiI6MTcxMzc1MjY4OCwiZXhwIjoxNzEzODM5MDg4fQ.YY4-SXRCVt8-Qe82KI7o8p4-OmgMN4DtAJzRttLJOsYnedgUqAd8g2gopBz89mM1s6WQFewG6Ek9VTE8UprbGzlwB3W31KGd5SFP5Pfj78rueM8u4BV3oW4QovUwcvHTUzBqoLLkl-qqkhxSs4dcHbj_RomseTrajnOS1Eb2wRg",
	}
	for i := 0; i < b.N; i++ {

		reponse, err := Get[map[string]interface{}](context.Background(), url, nil, header)
		if err != nil {
			b.Error(err)
			return
		}
		b.Log(reponse)
	}
}

type Response struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    struct {
		ImToken   string `json:"imToken"`
		ChatToken string `json:"chatToken"`
		UserID    string `json:"userID"`
	} `json:"data"`
}

func BenchmarkPOST(b *testing.B) {
	url := ""
	header := map[string]string{"content-type": "application/json",
		"operationID":   "xxx",
		"Authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJwaHBlcjY2Ni9qd3QiLCJ1aWQiOjYzNTUyNzY4MTI2OTYwODQ0OSwicGhvbmUiOiIxNzM3ODUxNjMyNSIsInN0YXR1cyI6MSwidGltZSI6MTcxMzc1MjY4OCwiand0X3NjZW5lIjoiYXBwbGljYXRpb24iLCJqdGkiOiJhcHBsaWNhdGlvbl82MzU1Mjc2ODEyNjk2MDg0NDkiLCJpYXQiOjE3MTM3NTI2ODgsIm5iZiI6MTcxMzc1MjY4OCwiZXhwIjoxNzEzODM5MDg4fQ.YY4-SXRCVt8-Qe82KI7o8p4-OmgMN4DtAJzRttLJOsYnedgUqAd8g2gopBz89mM1s6WQFewG6Ek9VTE8UprbGzlwB3W31KGd5SFP5Pfj78rueM8u4BV3oW4QovUwcvHTUzBqoLLkl-qqkhxSs4dcHbj_RomseTrajnOS1Eb2wRg",
	}
	param := map[string]interface{}{
		"areaCode":    "+86",
		"password":    "e9bc0e13a8a16cbb07b175d92a113126",
		"phoneNumber": "",
		"platform":    5,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reponse, err := POST[Response](context.Background(), url, param, header)
		if err != nil {
			b.Error(err)
			return
		}
		b.Log(reponse)
	}
}

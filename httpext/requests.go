package httpext

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"sync"
)

var clientPool = sync.Pool{
	New: func() interface{} {
		return resty.New()
	},
}

func Get[T any](url string, param map[string]string, header map[string]string) (T, error) {
	client := clientPool.Get().(*resty.Client)
	defer clientPool.Put(client)
	var result T
	response, err := client.R().EnableTrace().SetQueryParams(param).SetHeaders(header).Get(url)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(response.Body(), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func POST[T any](url string, param map[string]interface{}, header map[string]string) (T, error) {
	client := clientPool.Get().(*resty.Client)
	defer clientPool.Put(client)
	var result T
	response, err := client.R().SetBody(param).SetHeaders(header).Post(url)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(response.Body(), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

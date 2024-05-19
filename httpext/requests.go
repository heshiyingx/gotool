package httpext

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/go-resty/resty/v2"
	"sync"
)

var clientPool = sync.Pool{
	New: func() interface{} {
		return resty.New()
	},
}

func Get[T any](ctx context.Context, url string, param map[string]string, header map[string]string) (T, error) {
	client := clientPool.Get().(*resty.Client)
	defer clientPool.Put(client)
	var result T
	response, err := client.R().SetContext(ctx).EnableTrace().SetQueryParams(param).SetHeaders(header).Get(url)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(response.Body(), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func POST[T any](ctx context.Context, url string, param map[string]interface{}, header map[string]string) (T, error) {
	client := clientPool.Get().(*resty.Client)
	defer clientPool.Put(client)
	var result T
	response, err := client.SetRetryCount(3).R().SetContext(ctx).SetBody(param).SetHeaders(header).Post(url)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(response.Body(), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
func POSTXml[T any](ctx context.Context, url string, param interface{}, header map[string]string) (T, error) {
	client := clientPool.Get().(*resty.Client)
	defer clientPool.Put(client)
	var result T
	response, err := client.R().SetContext(ctx).SetBody(param).SetHeaders(header).Post(url)
	if err != nil {
		return result, err
	}
	err = xml.Unmarshal(response.Body(), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

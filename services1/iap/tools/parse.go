package tools

import (
	"bytes"
	"net/http"
)

// 处理Get请求
func Get(url string) (*http.Response, error) {
	var req *http.Request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(req)
}

// 处理POST请求
func Post(url string, body []byte, params map[string]string, headers map[string]string) (*http.Response, error) {
	if url == "" {
		return nil, ErrUrlError
	}
	var req *http.Request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	client := &http.Client{}
	return client.Do(req)
}

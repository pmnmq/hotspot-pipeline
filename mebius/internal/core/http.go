package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html/charset"
)

type HttpExecuter struct {
	client *http.Client
	Ctx    context.Context
	Proxy  string
}

func NewHttpExecuter(ctx context.Context) *HttpExecuter {
	return &HttpExecuter{
		client: &http.Client{},
		Ctx:    ctx,
	}
}

func (e *HttpExecuter) WithProxy(proxy string) error {
	e.Proxy = proxy
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	e.client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	return nil
}

func (e *HttpExecuter) Fetch(request *Request) (*Response, error) {
	req, err := request.toHttpRequest()
	if e.Ctx != nil {
		req = req.WithContext(e.Ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpResp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer httpResp.Body.Close()
	b, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	encoding, _, _ := charset.DetermineEncoding(b, httpResp.Header.Get("Content-Type"))
	decoder := encoding.NewDecoder()
	body, err := decoder.String(string(b))
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Protocol:   httpResp.Proto,
		Body:       body,
		Header:     httpResp.Header,
		Cookies:    httpResp.Cookies(),
	}
	return resp, nil
}

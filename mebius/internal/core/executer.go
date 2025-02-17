package core

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Executer interface {
	Fetch(*Request) (*Response, error)
	WithProxy(string) error
}

type HttpMethod string

const (
	GET  HttpMethod = "GET"
	POST HttpMethod = "POST"
)

const (
	userAgentText   = "User-Agent"
	refererText     = "Referer"
	contentTypeText = "Content-Type"
	JSONContentType = "application/json"
	FormContentType = "application/x-www-form-urlencoded"
)

type Request struct {
	URL     string
	Method  HttpMethod
	Header  http.Header
	Cookies []*http.Cookie
	Body    string
}

func NewRequest(url string, method HttpMethod) *Request {
	return &Request{
		URL:    url,
		Method: method,
	}
}

func (r *Request) AddHeader(key, value string) *Request {
	r.Header.Add(key, value)
	return r
}

func (r *Request) SetUserAgent(ua string) *Request {
	r.AddHeader(userAgentText, ua)
	return r
}

func (r *Request) SetReferer(referer string) *Request {
	r.AddHeader(referer, referer)
	return r
}

func (r *Request) SetCookie(cookie *http.Cookie) *Request {
	r.Cookies = append(r.Cookies, cookie)
	return r
}

func (r *Request) SetJsonBody(body map[string]any) *Request {
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	r.Body = string(b)
	r.Header.Add(contentTypeText, JSONContentType)
	return r
}

func (r *Request) SetJsonBodyString(body string) *Request {
	r.Body = body
	r.Header.Add(contentTypeText, JSONContentType)
	return r
}

func (r *Request) SetFormBodyString(body string) *Request {
	r.Body = body
	r.Header.Add(contentTypeText, FormContentType)
	return r
}

func (r *Request) SetFormBody(body map[string]any) *Request {
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	r.Body = string(b)
	r.Header.Add(contentTypeText, FormContentType)
	return r
}

func (r *Request) toHttpRequest() (*http.Request, error) {
	req, err := http.NewRequest(string(r.Method), r.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header
	if r.Body != "" && r.Method == POST {
		req.Body = io.NopCloser(strings.NewReader(r.Body))
	}

	for _, cookie := range r.Cookies {
		req.AddCookie(cookie)
	}
	return req, nil
}

type Response struct {
	Protocol   string
	StatusCode int
	Body       string
	Header     http.Header
	Cookies    []*http.Cookie
}

func (r *Response) GetCookieString() string {
	var cookieStr string
	for _, cookie := range r.Cookies {
		cookieStr += cookie.String()
	}
	return cookieStr
}

func (r *Response) JSON(v any) error {
	return json.Unmarshal([]byte(r.Body), v)
}

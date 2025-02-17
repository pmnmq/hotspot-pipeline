package core

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type BrowserExecuter struct {
	Ctx   context.Context
	Proxy string
}

func NewBrowserExecuter(ctx context.Context) *BrowserExecuter {
	return &BrowserExecuter{
		Ctx: ctx,
	}
}

func (b *BrowserExecuter) WithProxy(proxy string) error {
	_, err := url.Parse(proxy)
	b.Proxy = proxy
	return err
}

func (b *BrowserExecuter) Fetch(request *Request) (*http.Response, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	if b.Proxy != "" {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("proxy-server", b.Proxy),
		)
		allocCtx, allocCancel := chromedp.NewExecAllocator(b.Ctx, opts...)
		defer allocCancel()
		ctx, cancel = chromedp.NewContext(allocCtx)
	} else {
		ctx, cancel = chromedp.NewContext(b.Ctx)
	}

	defer cancel()
	var html string
	var respCookie string
	var cookieParams []*network.CookieParam
	for _, cookie := range request.Cookies {
		cookieParams = append(cookieParams, &network.CookieParam{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			HTTPOnly: cookie.HttpOnly,
			Secure:   cookie.Secure,
		})
	}
	chromedp.Run(ctx, network.Enable(), network.SetCookies(cookieParams))
	response, err := chromedp.RunResponse(
		ctx, chromedp.Navigate(request.URL),
		chromedp.OuterHTML("html", &html),
		chromedp.EvaluateAsDevTools(`document.cookie`, &respCookie),
	)
	if err != nil {
		return nil, err
	}
	resp := &http.Response{
		StatusCode: int(response.Status),
		Body:       io.NopCloser(strings.NewReader(html)),
	}
	return resp, nil
}

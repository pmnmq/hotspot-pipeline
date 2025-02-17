package core

import (
	"context"
	"io"
	"log"
	"testing"
	"time"
)

func TestBrowserExecuter(t *testing.T) {
	cases := []struct {
		url            string
		method         HttpMethod
		timeout        int
		wantStatusCode int
	}{
		{
			url:            "https://www.baidu.com",
			method:         "GET",
			timeout:        100,
			wantStatusCode: 200,
		},
	}

	for _, c := range cases {
		t.Run("Execute", func(t *testing.T) {
			t.Run("Should return response", func(t *testing.T) {
				// Given
				request := NewRequest(c.url, c.method)
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout)*time.Second)
				defer cancel()
				executer := &BrowserExecuter{
					Ctx: ctx,
				}
				resp, err := executer.Fetch(request)
				if err != nil {
					t.Fatalf("got error: %v", err)
				}
				if resp.StatusCode != c.wantStatusCode {
					t.Fatalf("got %d, want %d", resp.StatusCode, c.wantStatusCode)
				}
				body := resp.Body
				defer body.Close()
				b, err := io.ReadAll(body)
				if err != nil {
					t.Fail()
				}
				html := string(b)
				log.Println(html)
				// When
			})
		})
	}
}

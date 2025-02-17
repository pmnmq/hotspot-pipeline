package core

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestHttpExecuter(t *testing.T) {
	cases := []struct {
		url            string
		method         HttpMethod
		timeout        int
		wantStatusCode int
		header         http.Header
	}{
		{
			url:            "https://www.baidu.com",
			method:         GET,
			timeout:        100,
			wantStatusCode: 200,
			header: http.Header{
				"user-agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"},
			},
		},
	}

	for _, c := range cases {
		t.Run("Execute", func(t *testing.T) {
			t.Run("Should return response", func(t *testing.T) {
				request := NewRequest(c.url, c.method)
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout)*time.Second)
				defer cancel()
				executer := NewHttpExecuter(ctx)
				resp, err := executer.Fetch(request)
				if err != nil {
					t.Fatalf("got error: %v", err)
				}
				if resp.StatusCode != c.wantStatusCode {
					t.Fatalf("got %d, want %d", resp.StatusCode, c.wantStatusCode)
				}
				log.Println(resp.Body)
				// When
			})
		})
	}
}

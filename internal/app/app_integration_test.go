//go:build integration

package app_test

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/ianflanagan1/go-adapters/httpsend"

	"github.com/ianflanagan1/go-api-caller/internal/app"
	"github.com/ianflanagan1/go-api-caller/internal/core/posts"
)

func TestTheRealAPICreatesAPost(t *testing.T) {
	sender := httpsend.New(&http.Client{Timeout: 10 * time.Second}, 1<<20)
	var out bytes.Buffer
	if err := app.Publish(t.Context(), sender, &out, posts.BaseURL, []string{"Hello", "World"}); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); got != "Created post 101.\n" {
		t.Fatalf("got output %q", got)
	}
}

package app_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ianflanagan1/go-ianf-strict-ports/httpmsg"

	"github.com/ianflanagan1/go-api-caller/internal/app"
	"github.com/ianflanagan1/go-api-caller/internal/core/posts"
)

type fakeSender struct {
	resp httpmsg.Response
	err  error
	sent []httpmsg.Request
}

func (f *fakeSender) Send(_ context.Context, req httpmsg.Request) (httpmsg.Response, error) {
	f.sent = append(f.sent, req)
	return f.resp, f.err
}

func TestACreatedPostIsReportedWithItsID(t *testing.T) {
	sender := &fakeSender{resp: httpmsg.Response{Status: 201, Body: []byte(`{"id":101}`)}}
	var out bytes.Buffer
	if err := app.Publish(t.Context(), sender, &out, "https://x", []string{"Hello", "World"}); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); got != "Created post 101.\n" {
		t.Fatalf("got output %q", got)
	}
	if len(sender.sent) != 1 || sender.sent[0].URL != "https://x/posts" {
		t.Fatalf("sent %+v, want one request to https://x/posts", sender.sent)
	}
}

func TestInvalidArgumentsAreAUsageErrorAndSendNothing(t *testing.T) {
	sender := &fakeSender{}
	err := app.Publish(t.Context(), sender, &bytes.Buffer{}, "https://x", []string{"Hello"})
	if !errors.Is(err, posts.ErrUsage) {
		t.Fatalf("got %v, want ErrUsage", err)
	}
	if len(sender.sent) != 0 {
		t.Fatalf("sent %+v, want nothing", sender.sent)
	}
}

func TestAFailureToSendIsAnError(t *testing.T) {
	unreachable := errors.New("connection refused")
	err := app.Publish(t.Context(), &fakeSender{err: unreachable}, &bytes.Buffer{}, "https://x", []string{"Hello", "World"})
	if !errors.Is(err, unreachable) {
		t.Fatalf("got %v, want the send error", err)
	}
}

func TestAPostTheAPIDoesntCreateIsAnErrorWithTheReason(t *testing.T) {
	sender := &fakeSender{resp: httpmsg.Response{Status: 500}}
	var out bytes.Buffer
	err := app.Publish(t.Context(), sender, &out, "https://x", []string{"Hello", "World"})
	if err == nil || !strings.Contains(err.Error(), "the API didn't create the post: status 500") {
		t.Fatalf("got %v, want the API's refusal", err)
	}
	if out.Len() != 0 {
		t.Fatalf("got output %q, want none", out.String())
	}
}

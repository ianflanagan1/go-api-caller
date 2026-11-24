package posts_test

import (
	"errors"
	"testing"

	"github.com/ianflanagan1/go-ianf-strict-ports/httpmsg"

	"github.com/ianflanagan1/go-api-caller/internal/core/posts"
)

func TestATitleAndBodyAreParsedTrimmed(t *testing.T) {
	got, err := posts.ParseArgs([]string{" Hello ", "\tWorld\n"})
	if err != nil {
		t.Fatal(err)
	}
	if got != (posts.Input{Title: "Hello", Body: "World"}) {
		t.Fatalf("got %+v", got)
	}
}

func TestArgumentsThatArentATitleAndBodyAreAUsageError(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"none", nil},
		{"only a title", []string{"Hello"}},
		{"too many", []string{"Hello", "World", "again"}},
		{"a blank title", []string{"  ", "World"}},
		{"a blank body", []string{"Hello", ""}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := posts.ParseArgs(tc.args); !errors.Is(err, posts.ErrUsage) {
				t.Fatalf("got %v, want ErrUsage", err)
			}
		})
	}
}

func TestTheRequestPostsTheInputAsJSON(t *testing.T) {
	got := posts.Request("https://x/", posts.Input{Title: "Hello", Body: `say "hi"`})
	if got.Method != "POST" || got.URL != "https://x/posts" {
		t.Fatalf("got %s %s, want POST https://x/posts", got.Method, got.URL)
	}
	if ct := got.Header["Content-Type"]; len(ct) != 1 || ct[0] != "application/json; charset=UTF-8" {
		t.Fatalf("got Content-Type %v", ct)
	}
	if want := `{"title":"Hello","body":"say \"hi\""}`; string(got.Body) != want {
		t.Fatalf("got body %s, want %s", got.Body, want)
	}
}

func TestACreatedResponseWithAnIDIsCreated(t *testing.T) {
	got := posts.Assess(httpmsg.Response{Status: 201, Body: []byte(`{"title":"Hello","body":"World","id":101}`)})
	if got != (posts.Created{ID: 101}) {
		t.Fatalf("got %+v, want Created{101}", got)
	}
}

func TestAnyOtherResponseIsFailedWithTheReason(t *testing.T) {
	for _, tc := range []struct {
		name string
		resp httpmsg.Response
		want string
	}{
		{"another status", httpmsg.Response{Status: 500, Body: []byte(`{}`)}, "status 500"},
		{"a body that isn't JSON", httpmsg.Response{Status: 201, Body: []byte("<html>")}, "reading the response: invalid character '<' looking for beginning of value"},
		{"no ID", httpmsg.Response{Status: 201, Body: []byte(`{"title":"Hello"}`)}, "the response gave the post no ID"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := posts.Assess(tc.resp); got != (posts.Failed{Reason: tc.want}) {
				t.Fatalf("got %+v, want Failed{%q}", got, tc.want)
			}
		})
	}
}

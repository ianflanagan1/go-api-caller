// Package app publishes posts: it gathers the inputs, sends the request the
// core builds and acts on what the core makes of the response.
package app

import (
	"context"
	"fmt"
	"io"

	"github.com/ianflanagan1/go-ianf-strict-ports/httpmsg"

	"github.com/ianflanagan1/go-api-caller/internal/core/posts"
)

// Sender is the port for sending HTTP requests.
type Sender interface {
	Send(ctx context.Context, req httpmsg.Request) (httpmsg.Response, error)
}

// Publish publishes the post that args give to the API at baseURL, through
// sender, and writes the new post's ID to out. It returns an error when args
// aren't a post, the request can't be sent or the API doesn't create the post.
func Publish(ctx context.Context, sender Sender, out io.Writer, baseURL string, args []string) error {
	in, err := posts.ParseArgs(args)
	if err != nil {
		return fmt.Errorf("reading the arguments: %w", err)
	}
	resp, err := sender.Send(ctx, posts.Request(baseURL, in))
	if err != nil {
		return fmt.Errorf("sending the post: %w", err)
	}
	switch outcome := posts.Assess(resp).(type) {
	case posts.Failed:
		return fmt.Errorf("the API didn't create the post: %s", outcome.Reason)
	case posts.Created:
		if _, err := fmt.Fprintf(out, "Created post %d.\n", outcome.ID); err != nil {
			return fmt.Errorf("writing the result: %w", err)
		}
	}
	return nil
}

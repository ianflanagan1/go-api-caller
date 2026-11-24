// Command go-api-caller publishes a post to the JSONPlaceholder API and
// reports the ID the API gives it. It's the composition root: the only place
// that reads the configuration, builds adapters and passes them to the code
// that uses them.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ianflanagan1/go-adapters/httpsend"

	"github.com/ianflanagan1/go-api-caller/internal/app"
	"github.com/ianflanagan1/go-api-caller/internal/core/posts"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	stop()
	if err != nil {
		fmt.Fprintln(os.Stderr, "go-api-caller:", err)
		os.Exit(1)
	}
}

// run reads the configuration from args and the environment, builds the HTTP
// adapter and publishes the post, stopping when ctx is cancelled by an
// interrupt or SIGTERM. It writes the result to stdout. Name each remaining
// parameter, (stderr io.Writer), as run starts to use it.
func run(ctx context.Context, args []string, stdout, _ io.Writer) error {
	baseURL, ok := os.LookupEnv("API_BASE_URL")
	if !ok {
		baseURL = posts.BaseURL
	}
	sender := httpsend.New(&http.Client{Timeout: 10 * time.Second}, 1<<20)
	return app.Publish(ctx, sender, stdout, baseURL, args)
}

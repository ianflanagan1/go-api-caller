# go-api-caller

Publishes a post to the [JSONPlaceholder](https://jsonplaceholder.typicode.com) test API and reports the ID the API gives it. It's an example of a command on go-ianf-strict-baseline that calls a third-party API: it parses its arguments, sends a request, assesses the response and exits 0 on success or 1 with the reason on failure.

## Running

```sh
devbox run -- task build
./bin/go-api-caller "Hello" "World"
# Created post 101.
```

## Configuration

Read only in `cmd/go-api-caller`.

| Environment variable or flag | Meaning | Default |
|---|---|---|
| `API_BASE_URL` | The API's address | `https://jsonplaceholder.typicode.com` |

## Architecture

| Package | Role | Does |
|---|---|---|
| `internal/core/posts` | Core | Parses the arguments, builds the `httpmsg.Request`, assesses the `httpmsg.Response` into `Created` or `Failed` |
| `internal/app` | Shell | Declares the `Sender` port, calls the core, sends the request and acts on the outcome |
| `cmd/go-api-caller` | Composition root | Reads `API_BASE_URL`, builds go-adapters' `httpsend` and calls the shell |

`go tool go-ianf-strict-baseline map` lists each package's role, the ports the shell declares and the adapters that implement them. The rules are in `BASELINE.md`.

`devbox run -- task integration` also runs a test that publishes to the real API.

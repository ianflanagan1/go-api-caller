// Package posts decides how to publish a post to the JSONPlaceholder API
// (https://jsonplaceholder.typicode.com): it parses the caller's input,
// builds the request and assesses the response, as httpmsg values. It does
// no I/O.
package posts

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ianflanagan1/go-ianf-strict-ports/httpmsg"
)

// BaseURL is the API's address.
const BaseURL = "https://jsonplaceholder.typicode.com"

// ErrUsage is returned when the arguments aren't a title and a body.
var ErrUsage = errors.New("usage: go-api-caller <title> <body>")

// Input is a post to publish.
type Input struct {
	Title string
	Body  string
}

// ParseArgs reads a post from the command's arguments: a title and a body,
// neither of them blank.
func ParseArgs(args []string) (Input, error) {
	if len(args) != 2 {
		return Input{}, fmt.Errorf("%w: got %d arguments, want 2", ErrUsage, len(args))
	}
	title, body := strings.TrimSpace(args[0]), strings.TrimSpace(args[1])
	if title == "" || body == "" {
		return Input{}, fmt.Errorf("%w: the title and body can't be blank", ErrUsage)
	}
	return Input{Title: title, Body: body}, nil
}

// Request returns the request that publishes in to the API at baseURL.
func Request(baseURL string, in Input) httpmsg.Request {
	body, err := json.Marshal(struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}{in.Title, in.Body})
	if err != nil {
		panic(fmt.Sprintf("posts.Request: encoding a struct of strings: %v", err))
	}
	return httpmsg.Request{
		Method: "POST",
		URL:    strings.TrimRight(baseURL, "/") + "/posts",
		Header: map[string][]string{"Content-Type": {"application/json; charset=UTF-8"}},
		Body:   body,
	}
}

// Outcome is what the API's response to a Request means: Created or Failed.
//
//sumtype:decl
type Outcome interface{ isOutcome() }

// Created is a post the API accepted, with the ID it gave the post.
type Created struct{ ID int }

// Failed holds why the API didn't create the post.
type Failed struct{ Reason string }

func (Created) isOutcome() {}
func (Failed) isOutcome()  {}

// Assess decides what resp, the response to a Request, means.
func Assess(resp httpmsg.Response) Outcome {
	if resp.Status != 201 {
		return Failed{Reason: fmt.Sprintf("status %d", resp.Status)}
	}
	var created struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(resp.Body, &created); err != nil {
		return Failed{Reason: fmt.Sprintf("reading the response: %v", err)}
	}
	if created.ID < 1 {
		return Failed{Reason: "the response gave the post no ID"}
	}
	return Created{ID: created.ID}
}

package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/ides15/todoist/types"
)

var (
	client     = &Client{}
	testServer = &httptest.Server{}
)

func setup() {
	client, _ = NewClient("12345", nil)
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		params, _ := url.ParseQuery(buf.String())

		switch {
		case params.Get("resource_types") == "[\"projects\"]":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"projects":[{"id":1,"name":"Inbox"},{"id":2,"name":"Classes"}]}`))
			break
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"test":"hi"}`))
			break
		}
	}))
}

func TestNewClientOK(t *testing.T) {
	_, err := NewClient("12345", nil)
	if err != nil {
		t.Fatalf("expected nil error, received %v", err)
	}
}

func TestNewClientNilToken(t *testing.T) {
	_, err := NewClient("", nil)
	if err == nil {
		t.Fatalf("expected err, received %v", err)
	} else if err.Error() != types.ErrRequiredToken.Error() {
		t.Fatalf("expected %v, received %v", types.ErrRequiredToken.Error(), err)
	}
}

func TestNewRequestOKURL(t *testing.T) {
	setup()

	request, err := client.NewRequest("*", nil, nil)
	if err != nil {
		t.Fatalf("expected nil error, received %v", err)
	}

	if request.URL.String() != defaultBaseURL {
		t.Fatalf("expected %s, received %s", defaultBaseURL, request.URL.String())
	}
}

func TestNewRequestOKToken(t *testing.T) {
	setup()

	request, _ := client.NewRequest("*", nil, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`[^_]token=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) < 1 {
		t.Fatalf("expected a matching token in body, received %s", body)
	} else if matches[1] != "12345" {
		t.Log(body)
		t.Fatalf("expected token %s, received %s", "12345", matches[1])
	}
}

func TestBadNewRequest(t *testing.T) {
	setup()

	// ASCII control character will break `client.NewRequest`
	client.baseURL = "\t"

	_, err := client.NewRequest("*", nil, nil)
	if err == nil {
		t.Fatalf("expected err, received %v", err)
	}

	if err.Error() != types.ErrBuildRequest.Error() {
		t.Fatalf("expected %v, received %v", err.Error(), types.ErrBuildRequest.Error())
	}
}

func TestNewRequestSyncToken(t *testing.T) {
	setup()

	request, _ := client.NewRequest("*", nil, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`sync_token=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) < 1 {
		t.Fatalf("expected a matching sync_token in body, received %s", body)
	} else if matches[1] != "%2A" {
		t.Fatalf("expected synx_token '%s', received '%s'", "%2A", matches[1])
	}
}

func TestNewRequestNilSyncToken(t *testing.T) {
	setup()

	request, _ := client.NewRequest("", nil, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`sync_token=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) > 1 {
		t.Fatalf("expected no sync_token in body, received %s", body)
	}
}

func TestNewRequestContentType(t *testing.T) {
	setup()

	request, _ := client.NewRequest("*", nil, nil)

	expected := "application/x-www-form-urlencoded"
	if request.Header.Get("Content-Type") != expected {
		t.Fatalf("expected Content-Type of %s, received %s", expected, request.Header.Get("Content-Type"))
	}
}

func TestNewRequestUserAgent(t *testing.T) {
	setup()

	request, _ := client.NewRequest("*", nil, nil)

	expected := defaultUserAgent
	if request.Header.Get("User-Agent") != expected {
		t.Fatalf("expected User-Agent of %s, received %s", expected, request.Header.Get("User-Agent"))
	}
}

func TestNewRequestCommands(t *testing.T) {
	setup()

	commands := &[]types.Command{{
		Type: "project_add",
		Args: map[string]string{
			"arg": "test",
		},
		UUID:   "uuid",
		TempID: "tempID",
	}}

	request, _ := client.NewRequest("*", commands, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`commands=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) < 1 {
		t.Fatalf("expected matching commands in body, received %s", body)
	}
}

func TestNewRequestNilCommands(t *testing.T) {
	setup()

	request, _ := client.NewRequest("*", nil, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`commands=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) > 1 {
		t.Fatalf("expected no commands in body, received %s", body)
	}
}

func TestNewRequestResourceTypes(t *testing.T) {
	setup()

	resourceTypes := &[]string{"resource_type"}

	request, _ := client.NewRequest("*", nil, resourceTypes)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`resource_types=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) < 1 {
		t.Fatalf("expected matching resource_types in body, received %s", body)
	}
}

func TestNewRequestNilResourceTypes(t *testing.T) {
	setup()

	request, _ := client.NewRequest("*", nil, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`resource_types=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) > 1 {
		t.Fatalf("expected no resource_types in body, received %s", body)
	}
}

func TestDoRequestOK(t *testing.T) {
	setup()

	client.baseURL = testServer.URL

	request, _ := client.NewRequest("*", nil, nil)
	_, err := client.Do(context.Background(), request)
	if err != nil {
		t.Fatalf("expected no err, received %v", err)
	}
}

func TestDoRequestContextCancel(t *testing.T) {
	setup()

	client.baseURL = testServer.URL
	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	cancel()

	request, _ := client.NewRequest("*", nil, nil)
	_, err := client.Do(ctx, request)
	if err == nil {
		t.Fatalf("expected context cancelled error, received %v", err)
	}
}

func TestDoRequestError(t *testing.T) {
	setup()

	client.baseURL = testServer.URL

	request, _ := client.NewRequest("*", nil, nil)

	// Force error from `client.Do`
	request.URL = nil
	_, err := client.Do(context.Background(), request)
	if err == nil {
		t.Fatalf("expected Request.URL error, received %v", err)
	}
}

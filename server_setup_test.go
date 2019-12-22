package todoist_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/ides15/todoist"
)

var (
	TestClient = &todoist.Client{}
	TestServer = &httptest.Server{}
)

func Setup() {
	TestClient, _ = todoist.NewClient("12345", nil)

	TestServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/not-found":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
					"projects": [
						{
							"id": 1
						},
						{
							"id": 2
						}
					]
				}`))
			break
		case "/no-projects":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
					"projects": []
				}`))
			break
		case "/bad-json":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
					"projects": {
						"bad": true
					}
				}`))
			break
		case "/AUTH_CSRF_ERROR":
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{
					"error_tag": "AUTH_CSRF_ERROR",
					"error_code": 0,
					"http_code": 403,
					"error_extra": {
						"retry_after": 2,
						"access_type": "web_session"
					},
					"error": "AUTH_CSRF_ERROR"
				}`))
			break
		case "/AUTH_INVALID_TOKEN":
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{
					"error_tag": "AUTH_INVALID_TOKEN",
					"error_code": 401,
					"http_code": 403,
					"error_extra": {
						"retry_after": 2,
						"access_type": "access_token"
					},
					"error": "Invalid token"
				}`))
			break
		case "/invalid-error":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{`))
			break
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
					"projects": [
						{
							"id": 1,
							"name": "Inbox"
						},
						{
							"id": 2,
							"name": "Classes"
						}
					]
				}`))
			break
		}
	}))
}

package main

import (
  "fmt"

  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/aws/aws-lambda-go/events"
  "github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Run("Unable to get Go versions JSON", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))
		defer ts.Close()
    
		GetGoVersionsURL = ts.URL

		_, err := handler(events.APIGatewayProxyRequest{})
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

    expectedError, _ := err.(*json.SyntaxError)
    assert.Equal(t, expectedError, err)
	})

	t.Run("Non 200 Response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()

		GetGoVersionsURL = ts.URL

		_, err := handler(events.APIGatewayProxyRequest{})
		if err == nil {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}

    assert.Equal(t, ErrNon200Response, err)
	})

	t.Run("No Go versions error returned", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, `[]`)
		}))
		defer ts.Close()

		GetGoVersionsURL = ts.URL

		_, err := handler(events.APIGatewayProxyRequest{})
		if err == nil {
			t.Fatal("Everything should be ok; ", err)
		}

    assert.Equal(t, ErrNoVersions, err)
	})

	t.Run("Successful Request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, `[{"version": "1.21.4","stable": true,"release_url": "https://github.com/actions/go-versions/releases/tag/1.21.4-6807962903","files": [{"filename": "go-1.21.4-darwin-arm64.tar.gz","arch": "arm64","platform": "darwin","download_url": "https://github.com/actions/go-versions/releases/download/1.21.4-6807962903/go-1.21.4-darwin-arm64.tar.gz"},{"filename": "go-1.21.4-darwin-x64.tar.gz","arch": "x64","platform": "darwin","download_url": "https://github.com/actions/go-versions/releases/download/1.21.4-6807962903/go-1.21.4-darwin-x64.tar.gz"},{"filename": "go-1.21.4-linux-x64.tar.gz","arch": "x64","platform": "linux","download_url": "https://github.com/actions/go-versions/releases/download/1.21.4-6807962903/go-1.21.4-linux-x64.tar.gz"},{"filename": "go-1.21.4-win32-x64.zip","arch": "x64","platform": "win32","download_url": "https://github.com/actions/go-versions/releases/download/1.21.4-6807962903/go-1.21.4-win32-x64.zip"}]}]`)
		}))
		defer ts.Close()

		GetGoVersionsURL = ts.URL

		response, err := handler(events.APIGatewayProxyRequest{})
		if err != nil {
			t.Fatal("Everything should be ok; ", err)
		}

    assert.Equal(t, `{"version":"1.21.4","stable":true,"release_url":"https://github.com/actions/go-versions/releases/tag/1.21.4-6807962903","files":[{"filename":"go-1.21.4-darwin-arm64.tar.gz","arch":"arm64","platform":"darwin","download_url":"https://github.com/actions/go-versions/releases/download/1.21.4-6807962903/go-1.21.4-darwin-arm64.tar.gz"},{"filename":"go-1.21.4-darwin-x64.tar.gz","arch":"x64","platform":"darwin","download_url":"https://github.com/actions/go-versions/releases/download/1.21.4-6807962903/go-1.21.4-darwin-x64.tar.gz"},{"filename":"go-1.21.4-linux-x64.tar.gz","arch":"x64","platform":"linux","download_url":"https://github.com/actions/go-versions/releases/download/1.21.4-6807962903/go-1.21.4-linux-x64.tar.gz"},{"filename":"go-1.21.4-win32-x64.zip","arch":"x64","platform":"win32","download_url":"https://github.com/actions/go-versions/releases/download/1.21.4-6807962903/go-1.21.4-win32-x64.zip"}]}`, response.Body)
	})
}

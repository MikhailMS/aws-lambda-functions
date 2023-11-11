package main

import (
	"errors"
	"fmt"
	"io"

	"net/http"
  "encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// DefaultHTTPGetAddress Default Address
	GetGoVersionsURL = "https://raw.githubusercontent.com/actions/go-versions/main/versions-manifest.json"

	// ErrNoIP No IP found in response
	ErrNoVersions = errors.New("No Go versions detected")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)


// =====
// Unmarshall Go versions
type Versions struct {
	Versions []Version
}

type Version struct {
  Version    string `json:"version"`
  Stable     bool   `json: stable`
  ReleaseUrl string `json:"release_url"`
  Files      []File `json: files`
}

type File struct {
  Filename    string `json:"filename"`
  Arch        string `json:"arch"`
  Platform    string `json:"platform"`
  DownloadUrl string `json:"download_url"`
}
// =====

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  // Step 1. Send a request to get Go versions
	resp, err := http.Get(GetGoVersionsURL)
	if err != nil {
    return events.APIGatewayProxyResponse{
      Body: fmt.Sprintf("Error: %v", err)
    }, err
	}

	if resp.StatusCode != 200 {
    return events.APIGatewayProxyResponse{
      Body:       fmt.Sprintf("Error: %v", ErrNon200Response),
      StatusCode: resp.StatusCode,
    }, ErrNon200Response
	}

  // Step 2. Unmarshall reponse into Go struct
	rawJson, err := io.ReadAll(resp.Body)
	if err != nil {
    return events.APIGatewayProxyResponse{
      Body: fmt.Sprintf("Error: %v", err)
    }, err
	}

  var goVersions Versions

  err = json.Unmarshal(rawJson, &goVersions)
	if err != nil {
    return events.APIGatewayProxyResponse{
      Body: fmt.Sprintf("Error: %v", err)
    }, err
	}

	if len(goVersions) == 0 {
    return events.APIGatewayProxyResponse{
      Body:       fmt.Sprintf("Error: %v", ErrNoIP),
      StatusCode: 404,
    }, ErrNoIP
	}

  // Step 3. Marshall struct so it could be returned to the Lambda caller
  responseBody, err := json.Marshal(goVersions[0:5])
 	if err != nil {
    return events.APIGatewayProxyResponse{
      Body: fmt.Sprintf("Error: %v", err)
    }, err
	}

	return events.APIGatewayProxyResponse{
		Body:       responseBody,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}


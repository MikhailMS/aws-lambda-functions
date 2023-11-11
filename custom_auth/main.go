package main

import (
	"errors"
	"fmt"
	"net/http"

  "github.com/lestrrat-go/jwx/jwt"
  "github.com/lestrrat-go/jwx/jwa"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
  // secret to decode/parse JWT token would be set at build time
  TokenSecret string

  // algorithm used to encode/decode JWT token
  TokenAlg = jwa.RS256

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
  // Step 1. Get Authorization token from Request
  // bounds := len(event.AuthorizationToken)
	// token := event.AuthorizationToken[7:bounds]
	token := event.AuthorizationToken

  // Step 2. Validate token
	parsedToken, err := jwt.Parse([]byte(token), jwt.WithKey(TokenAlg, TokenSecret), jwt.WithValidate(true))

  // Step 3. If token invalid, return DENY reponse
  //         otherwises return ALLOW
  if err != nil {

		return events.APIGatewayCustomAuthorizerResponse{
			PrincipalID: parsedToken.Get("principalID"),
			PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
				Version: "2012-10-17",
				Statement: []events.IAMPolicyStatement{
					{
						Action:   []string{"execute-api:Invoke"},
						Effect:   "Deny",
						Resource: []string{"*"},
					},
				},
			},
			UsageIdentifierKey: "",
		}, err
	}

	return events.APIGatewayCustomAuthorizerResponse{
    PrincipalID: parsedToken.Get("principalID"),
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: []string{event.MethodArn},
				},
			},
		},
		UsageIdentifierKey: "",
	}, nil

}

func main() {
	lambda.Start(handler)
}


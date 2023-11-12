package main

import (
  "context"
  "errors"
  "fmt"

  "github.com/golang-jwt/jwt/v5"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
  // secret to decode/parse JWT token would be set at build time
  TokenSecret string
)


// Help function to generate an IAM policy
func generatePolicy(principalId, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId, UsageIdentifierKey: ""}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	return authResponse
}

func handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
  // Step 1. Get Authorization token from Request
  // bounds := len(event.AuthorizationToken)
	// token := event.AuthorizationToken[7:bounds]
	token := event.AuthorizationToken

  if len(token) == 0 {
    return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
  }

  // Step 2. Parse token
  parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
    _, ok := t.Method.(*jwt.SigningMethodHMAC)
    if !ok {
      return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"]))
    }
    
    return []byte(TokenSecret), nil
  })
  if err != nil {
    return events.APIGatewayCustomAuthorizerResponse{}, err
	}

  // Step 3. If token invalid, return DENY reponse
  if !parsedToken.Valid {
    return generatePolicy("", "Deny", "*"), err
  }

  // Step 4. At this point token is Valid, so we need to get principal out and return ALLOW response
  claims, ok := parsedToken.Claims.(jwt.MapClaims)
  if ok {
    principal := claims["principalID"].(string)
    return generatePolicy(fmt.Sprintf("%v", principal), "Allow", event.MethodArn), nil
  }

  // This is only called if everything is great, but principalID is not encoded into token
  return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Missing claims")
}

func main() {
	lambda.Start(handler)
}


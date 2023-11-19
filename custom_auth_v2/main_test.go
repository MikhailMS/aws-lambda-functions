package main

import (
  "errors"
  "context"
  "testing"
  "time"

  "github.com/aws/aws-lambda-go/events"
  "github.com/golang-jwt/jwt/v5"
  "github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
  TokenSecret = "secret"

	t.Run("Empty Token", func(t *testing.T) {

    headers := make(map[string]string)
    headers["Authorization"] = ""

    _, err := handler(context.Background(), events.APIGatewayV2CustomAuthorizerV2Request{
      Type:     "REQUEST",
      RouteArn: "arn:aws:execute-api:eu-west-2:123456789012:/test/POST/test",
      Headers:  headers,
    })

		if err == nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

    assert.Equal(t, errors.New("Error: Invalid auth token"), err)
	})

	t.Run("Incorrect Token", func(t *testing.T) {
    headers := make(map[string]string)
    headers["Authorization"] = "invalid_token"

    expectedResponse := events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{}
    expectedError    := "token is malformed: token contains an invalid number of segments"

    response, err := handler(context.Background(), events.APIGatewayV2CustomAuthorizerV2Request{
      Type:     "REQUEST",
      RouteArn: "arn:aws:execute-api:eu-west-2:123456789012:/test/POST/test",
      Headers:  headers,
    })
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

    assert.Equal(t, expectedResponse, response)
    assert.Equal(t, expectedError, err.Error())
	})

	t.Run("Valid Token", func(t *testing.T) {
    // Create valid token and sign it
    valid_token := jwt.New(jwt.SigningMethodHS256)

    claims := valid_token.Claims.(jwt.MapClaims)
    claims["principalID"] = "test_user"

    signed_token, err := valid_token.SignedString([]byte(TokenSecret))
    if err != nil {
      t.Fatal(err, TokenSecret, []byte(TokenSecret))
    }

    // Create expectedResponse
    expectedResponse := events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{
      PrincipalID: "test_user",
      PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
        Version: "2012-10-17",
        Statement: []events.IAMPolicyStatement{
          {
            Action:   []string{"execute-api:Invoke"},
            Effect:   "Allow",
            Resource: []string{"arn:aws:execute-api:eu-west-2:123456789012:/test/POST/test"},
          },
        },
      },
    }

    headers := make(map[string]string)
    headers["Authorization"] = signed_token

    response, err := handler(context.Background(), events.APIGatewayV2CustomAuthorizerV2Request{
      Type:     "REQUEST",
      RouteArn: "arn:aws:execute-api:eu-west-2:123456789012:/test/POST/test",
      Headers:  headers,
    })
		if err != nil {
      t.Fatal("Everything should be ok; ", err, "; Signed Token: ", signed_token)
		}

    assert.Equal(t, expectedResponse, response)
	})

	t.Run("Valid Token, Missing Claim", func(t *testing.T) {
    expectedError    := "Error: Missing claims"

    // Create valid token and sign it
    valid_token := jwt.New(jwt.SigningMethodHS256)

    signed_token, err := valid_token.SignedString([]byte(TokenSecret))
    if err != nil {
      t.Fatal(err, TokenSecret, []byte(TokenSecret))
    }

    headers := make(map[string]string)
    headers["Authorization"] = signed_token

    // Create expectedResponse
    expectedResponse := events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{}

    response, err := handler(context.Background(), events.APIGatewayV2CustomAuthorizerV2Request{
      Type:     "REQUEST",
      RouteArn: "arn:aws:execute-api:eu-west-2:123456789012:/test/POST/test",
      Headers:  headers,
    })
		if err == nil {
      t.Fatal("Everything should be ok; ", err, "; Signed Token: ", signed_token)
		}

    assert.Equal(t, expectedResponse, response)
    assert.Equal(t, expectedError, err.Error())
	})

  t.Run("Valid Token, Expired", func(t *testing.T) {
    // Create invalid token and sign it
    expired_token := jwt.New(jwt.SigningMethodHS256)

    claims := expired_token.Claims.(jwt.MapClaims)
    claims["exp"] = time.Now().Add(-1 * time.Hour).Unix()
    claims["iat"] = time.Now().Add(-2 * time.Hour).Unix()
    claims["principalID"] = "test_user"

    signed_expired_token, err := expired_token.SignedString([]byte(TokenSecret))
    if err != nil {
      t.Fatal(err, TokenSecret, []byte(TokenSecret))
    }

    headers := make(map[string]string)
    headers["Authorization"] = signed_expired_token

    expectedResponse := events.APIGatewayV2CustomAuthorizerIAMPolicyResponse{}
    expectedError    := "token has invalid claims: token is expired"

    response, err := handler(context.Background(), events.APIGatewayV2CustomAuthorizerV2Request{
      Type:     "REQUEST",
      RouteArn: "arn:aws:execute-api:eu-west-2:123456789012:/test/POST/test",
      Headers:  headers,
    })
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

    assert.Equal(t, expectedResponse, response)
    assert.Equal(t, expectedError, err.Error())
	})
}

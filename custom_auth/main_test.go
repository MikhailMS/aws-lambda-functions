package main

import (
  "errors"
  "context"
  "testing"

  "github.com/aws/aws-lambda-go/events"
  "github.com/golang-jwt/jwt/v5"
  "github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
  TokenSecret = "secret"

	t.Run("Empty Token", func(t *testing.T) {

    _, err := handler(context.Background(), events.APIGatewayCustomAuthorizerRequest{
      Type:               "TOKEN",
      AuthorizationToken: "",
      MethodArn:          "arn:aws:execute-api:eu-west-2:123456789012:/test/POST/test",
    })

		if err == nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

    assert.Equal(t, errors.New("Error: Invalid token"), err)
	})

	t.Run("Incorrect Token", func(t *testing.T) {
    expectedResponse := events.APIGatewayCustomAuthorizerResponse{}
    expectedError    := "token is malformed: token contains an invalid number of segments"

    response, err := handler(context.Background(), events.APIGatewayCustomAuthorizerRequest{
      Type:               "TOKEN",
      AuthorizationToken: "invalid_token",
      MethodArn:          "arn:aws:execute-api:eu-west-2:123456789012:/test/POST/test",
    })
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}

    assert.Equal(t, expectedResponse, response)
    assert.Equal(t, expectedError, err.Error())
	})

  // t.Run("Invalid Token", func(t *testing.T) {
  //   expectedResponse := events.APIGatewayCustomAuthorizerResponse{
  //     PrincipalID: "",
  //     PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
  //       Version: "2012-10-17",
  //       Statement: []events.IAMPolicyStatement{
  //         {
  //           Action:   []string{"execute-api:Invoke"},
  //           Effect:   "Deny",
  //           Resource: []string{"*"},
  //         },
  //       },
  //     },
  //     UsageIdentifierKey: "",
  //   }

  //   response, err := handler(context.Background(), events.APIGatewayCustomAuthorizerRequest{
  //     Type:               "TOKEN",
  //     AuthorizationToken: "invalid_token",
  //     MethodArn:          "arn:aws:execute-api:eu-west-2:123456789012:/test/POST/test",
  //   })
		// if err == nil {
			// t.Fatal("Error failed to trigger with an invalid request")
		// }

  //   t.Fatal(err)

  //   assert.Equal(t, expectedResponse, response)
	// })

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
    expectedResponse := events.APIGatewayCustomAuthorizerResponse{
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
      UsageIdentifierKey: "",
    }

    response, err := handler(context.Background(), events.APIGatewayCustomAuthorizerRequest{
      Type:               "TOKEN",
      AuthorizationToken: signed_token,
      MethodArn:          "arn:aws:execute-api:eu-west-2:123456789012:/test/POST/test",
    })
		if err != nil {
      t.Fatal("Everything should be ok; ", err, "; Signed Token: ", signed_token)
		}

    assert.Equal(t, expectedResponse, response)
	})
}

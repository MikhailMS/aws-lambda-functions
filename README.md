# AWS Lambda Functions

Simple AWS Lambda Functions written in Go to facilitate my exploration of building applications in AWS

1. `return_ip`         - returns IP address of the Lambda function
2. `fetch_go_versions` - returns JSON with recent 5 Go versions
3. `custom_auth`       - custom Lambda authenticator that controls access to above 2 functions when calling via API Gateway


## Notes
1. This repo uses Github Action to build & push compiled Golang into S3 and from there Functions are deployed with Terraform
2. Lambda function gets build by Golang version 1.20

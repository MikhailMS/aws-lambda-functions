# AWS Lambda Functions

Simple AWS Lambda Functions written in Go to facilitate my exploration of building applications in AWS

This project is part of one big project where I research how to build infrastructure in AWS for 2-tier application (and application as well, of course) :
1. [Terraform](https://github.com/MikhailMS/aws-2tier-lambda-api) contains all the terraform code to deploy required infra & application code
2. [Simple Web GUI](https://github.com/MikhailMS/aws-simple-web-gui) contains code for simple Web GUI to bridge the gap between user and Lambda functions
3. Lambda functions (this project) contain code for 3 Lambda functions that replicate simple backend functions
    1. `return_ip`         - returns IP address of the Lambda function
    2. `fetch_go_versions` - returns JSON with recent 5 Go versions
    3. `custom_auth`       - custom Lambda authorizer (only supports payload format `version 1.0`) that controls access to functions `1 & 2` when calling via API Gateway
    4. `custom_auth_v2`    - custom Lambda authorizer (only supports payload format `version 2.0`) that controls access to functions `1 & 2` when calling via API Gateway


## Notes
1. This repo uses Github Action to build & push compiled Golang into S3 and from there Functions are deployed with Terraform
2. Lambda function gets build by Golang version 1.20

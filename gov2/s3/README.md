# AWS SDK for Go V2 code examples for Amazon S3

## Purpose

These examples demonstrate how to perform several Amazon Simple Storage Service 
(Amazon S3) operations using version 2 of the AWS SDK for Go.

## Code examples

### Scenario examples
* [`common/`](common/) -- The most common S3 operations: Creating a bucket, managing objects within it, and deleting objects/buckets.

### API Examples
- [`GetBucketAcl`](GetBucketAcl/) and [`GetObjectAcl`](GetObjectAcl/) -- Two examples of working with ACLs on S3 objects.



## ⚠ Important

- As an AWS best practice, grant this code least privilege, or only the 
  permissions required to perform a task. For more information, see 
  [Grant Least Privilege](https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html#grant-least-privilege) 
  in the *AWS Identity and Access Management 
  User Guide*.
- This code has not been tested in all AWS Regions. Some AWS services are 
  available only in specific Regions. For more information, see the 
  [AWS Region Table](https://aws.amazon.com/about-aws/global-infrastructure/regional-product-services/)
  on the AWS website.
- Running this code might result in charges to your AWS account.


## Running the code

From the directory you wish to run the sample of, do the following:

```
go mod tidy
go run .
```

## Running the tests

From a directory containing `go.mod`, sue `go test` to run all unit tests:

```
go test ./...
```

This will test all modules in the current folder and any submodules.

### Prerequisites

You must have an AWS account, and have your default credentials and AWS Region
configured as described in
[Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)
in the AWS SDK for Go Developer Guide.

You must have at least Go 1.17 installed. 

## Additional information

- [AWS SDK for Go V3 Amazon S3 service reference](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3)
- [Amazon S3 documentation](https://docs.aws.amazon.com/s3)

---


Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0

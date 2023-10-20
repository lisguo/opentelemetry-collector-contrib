// Copyright The OpenTelemetry Authors
// Portions of this file Copyright 2018-2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aws // import "github.com/amazon-contributing/opentelemetry-collector-contrib/override/aws"

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"go.uber.org/zap"
)

const (
	DefaultIMDSRetries = 1
)

type IMDSRetryer struct {
	client.DefaultRetryer
	logger *zap.Logger
}

// NewIMDSRetryer allows us to retry imds errors
func NewIMDSRetryer(retryNumber int) IMDSRetryer {
	imdsRetryer := IMDSRetryer{
		DefaultRetryer: client.DefaultRetryer{
			NumMaxRetries: retryNumber,
		},
	}
	logger, err := zap.NewDevelopment()
	if err == nil {
		imdsRetryer.logger = logger
	}
	return imdsRetryer
}

func (r IMDSRetryer) ShouldRetry(req *request.Request) bool {
	// there is no enum of error codes
	// EC2MetadataError is not retryable by default
	// Fallback to SDK's built in retry rules
	shouldRetry := false
	if awsError, ok := req.Error.(awserr.Error); r.DefaultRetryer.ShouldRetry(req) || (ok && awsError != nil && awsError.Code() == "EC2MetadataError") {
		shouldRetry = true
	}
	if r.logger != nil {
		r.logger.Debug("imds error : ", zap.Bool("shouldRetry", shouldRetry), zap.Error(req.Error))
	}
	return shouldRetry
}

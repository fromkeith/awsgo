/*
 * Copyright (c) 2013, fromkeith
 * All rights reserved.
 * 
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted provided that the following conditions are met:
 * 
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 * 
 * * Redistributions in binary form must reproduce the above copyright notice, this
 *   list of conditions and the following disclaimer in the documentation and/or
 *   other materials provided with the distribution.
 * 
 * * Neither the name of the fromkeith nor the names of its
 *   contributors may be used to endorse or promote products derived from
 *   this software without specific prior written permission.
 * 
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
 * ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package cloudwatch


import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/fromkeith/awsgo"
)



type DescribeLogGroupsRequest struct {
    awsgo.RequestBuilder

    LogGroupNamePrefix      string          `json:"logGroupNamePrefix,omitempty"`
    Limit                   *int            `json:"limit,omitempty"`
    NextToken               string          `json:"nextToken,omitempty"`
}

type LogGroup struct {
    Arn                     string          `json:"arn"`
    CreationTime            int64           `json:"creationTime"`
    LogGroupName            string          `json:"logGroupName"`
    MetricFilterCount       int64           `json:"metricFilterCount"`
    RetentionInDays         int64           `json:"retentionInDays"`
    StoredBytes             int64           `json:"storedBytes"`
}


type DescribeLogGroupsResponse struct {
    NextToken               string          `json:"nextToken"`
    LogGroups               []LogGroup      `json:"logGroups"`
}


func NewDescribeLogGroupsRequest() *DescribeLogGroupsRequest {
    req := new(DescribeLogGroupsRequest)
    req.Host.Service = "logs"
    req.Host.Region = "us-east-1"
    req.Host.Domain = "amazonaws.com"
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = "Logs_20140328.DescribeLogGroups"
    req.Headers["Content-Type"] = "application/x-amz-json-1.1"
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"
    return req
}

func (req * DescribeLogGroupsRequest) SetLimit(limit int) {
    req.Limit = &limit
}

func (req * DescribeLogGroupsRequest) VerifyInput() (error) {
    return nil
}

func (req DescribeLogGroupsRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := CheckForErrorResponse(response, statusCode); err != nil {
        return err
    }
    if statusCode != 200 {
        return errors.New(fmt.Sprintf("Bad Status code: %d", statusCode))
    }
    resp := new(DescribeLogGroupsResponse)
    if err := json.Unmarshal(response, resp); err != nil {
        return err
    }
    return resp
}


func (req DescribeLogGroupsRequest) Request() (*DescribeLogGroupsResponse, error) {
    request, err := awsgo.NewAwsRequest(&req, req)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&req)
    if resp == nil {
        return nil, err
    }
    return resp.(*DescribeLogGroupsResponse), err
}

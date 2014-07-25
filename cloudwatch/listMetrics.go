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
    "github.com/fromkeith/awsgo"
    "fmt"
    "net/url"
    "encoding/xml"
)


type ListMetricsRequest struct {
    awsgo.RequestBuilder

    Dimensions                  []MetricDimensions
    MetricName                  string
    Namespace                   string
    NextToken                   string
}

type MetricInfo struct {
    Dimensions                  []MetricDimensions  `xml:"Dimensions>member"`
    MetricName                  string
    Namespace                   string
}

type ListMetricsResult struct {
    NextToken                string
    Metrics                  []MetricInfo      `xml:"Metrics>member"`
}


type ListMetricsResponse struct {
    ResponseMetadata                awsgo.ResponseMetaData
    ListMetricsResult               ListMetricsResult
}


func NewListMetricsRequest() *ListMetricsRequest {
    req := new(ListMetricsRequest)
    req.Host.Service = "monitoring"
    req.Host.Region = ""
    req.Host.Domain = "amazonaws.com"
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    req.Headers = make(map[string]string)
    req.RequestMethod = "GET"
    req.CanonicalUri = "/"
    return req
}


func (req * ListMetricsRequest) VerifyInput() (error) {
    req.Host.Service = "monitoring"
    if len(req.Host.Region) == 0 {
        req.Host.Region = "us-east-1"
    }

    vals := url.Values{}
    vals.Set("Action", "ListMetrics")
    vals.Set("Version", "2010-08-01")

    for i := range req.Dimensions {
        vals.Set(fmt.Sprintf("Dimensions.member.%d.Name", i + 1), req.Dimensions[i].Name)
        vals.Set(fmt.Sprintf("Dimensions.member.%d.Value", i + 1), req.Dimensions[i].Value)
    }
    if req.MetricName != "" {
        vals.Set("MetricName", req.MetricName)
    }
    if req.Namespace != "" {
        vals.Set("Namespace", req.Namespace)
    }
    if req.NextToken != "" {
        vals.Set("NextToken", req.NextToken)
    }
    req.CanonicalUri = "/?" + vals.Encode()
    return nil
}

func (req ListMetricsRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := awsgo.CheckForErrorXml(response); err != nil {
        return err
    }
    reqesponse := new(ListMetricsResponse)
    if err := xml.Unmarshal(response, reqesponse); err != nil {
        return err
    }
    return reqesponse
}

func (req ListMetricsRequest) Request() (*ListMetricsResponse, error) {
    request, err := awsgo.BuildEmptyContentRequest(&req)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&req)
    if resp == nil {
        return nil, err
    }
    return resp.(*ListMetricsResponse), err
}

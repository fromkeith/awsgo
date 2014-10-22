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

package ec2

import (
    "github.com/fromkeith/awsgo"
    "encoding/xml"
    "net/url"
    "fmt"
)


const (
    // tag key
    FILTER_key = "key"
    // the resource id
    FILTER_resource_id = "resource-id"
    // the resource type
    FILTER_resource_type = "resource-type"
    // the tag value
    FILTER_value = "value"
)



type DescribeTagsRequest struct {
    awsgo.RequestBuilder

    MaxResults          int
    NextToken           string
    Filters             []DescribeFilter

}


type TagSetItem struct {
    ResourceId          string                  `xml:"resourceId"`
    ResourceType        string                  `xml:"resourceType"`
    Key                 string                  `xml:"key"`
    Value               string                  `xml:"value"`
}

type DescribeTagsResponse struct {
    RequestId               string              `xml:"requestId"`
    TagSet                  []TagSetItem        `xml:"tagSet>item"`
    NextToken               string              `xml:"nextToken"`
    StatusCode              int
}


func NewDescribeTagsRequest() *DescribeTagsRequest {
    req := new(DescribeTagsRequest)
    req.Host.Service = "ec2"
    req.Host.Region = ""
    req.Host.Domain = "amazonaws.com"
    req.Headers = make(map[string]string)
    req.RequestMethod = "GET"
    req.CanonicalUri = "/"
    return req
}



func (req * DescribeTagsRequest) VerifyInput() (error) {
    val := make(url.Values)
    val.Set("Action", "DescribeTags")
    val.Set("Version", "2014-02-01")
    if req.NextToken != "" {
        val.Set("NextToken", req.NextToken)
    }
    if req.MaxResults > 0 {
        val.Set("MaxResults", fmt.Sprintf("%d", req.MaxResults))
    }

    for i := range req.Filters {
        val.Set(fmt.Sprintf("Filter.%d.Name", i + 1), req.Filters[i].Name)
        for j := range req.Filters[i].Value {
            val.Set(fmt.Sprintf("Filter.%d.Value.%d", i + 1, j + 1), req.Filters[i].Value[j])
        }
    }

    req.CanonicalUri = req.CanonicalUri + "?" + val.Encode()
    return nil

}


func (req DescribeTagsRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    giResponse := new(DescribeTagsResponse)
    xml.Unmarshal(response, giResponse)
    giResponse.StatusCode = statusCode
    return giResponse
}


func (gir DescribeTagsRequest) Request() (*DescribeTagsResponse, error) {
    request, err := awsgo.BuildEmptyContentRequest(&gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS2
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*DescribeTagsResponse), err
}

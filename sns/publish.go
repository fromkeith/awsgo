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

package sns

import (
    "github.com/fromkeith/awsgo"
    "encoding/xml"
    "net/url"
    "fmt"
)

type MessageAttributeValue struct {
    //BinaryValue             []byte // don't know how this is meant to be encoded
    DataType                string
    StringValue             string
}

type PublishRequest struct {
    awsgo.RequestBuilder

    Message                     string
    MessageAttributes           map[string]MessageAttributeValue
    MessageStructure            string
    Subject                     string
    TargetArn                   string
    TopicArn                    string
}

type PublishResult struct {
    MessageId               string
}
type ResponseMetadata struct {
    RequestId               string
}

type PublishResponse struct {
    PublishResult           PublishResult
    ResponseMetadata        ResponseMetadata
}

type PublishError struct {
    Response        *PublishResponse
    Code            int
    RawResponse     string
}

func (p PublishError) Error() string {
    return fmt.Sprintf("Response was %d", p.Code)
}


func NewPublishRequest() *PublishRequest {
    req := new(PublishRequest)
    req.Host.Service = "sns"
    req.Host.Region = ""
    req.Host.Domain = "amazonaws.com"
    req.Headers = make(map[string]string)
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"
    return req
}



func (req * PublishRequest) VerifyInput() (error) {
    val := make(url.Values)
    val.Set("Action", "Publish")
    val.Set("Version", "2010-03-31")
    val.Set("Message", req.Message)
    if req.MessageStructure != "" {
        val.Set("MessageStructure", req.MessageStructure)
    }
    if req.Subject != "" {
        val.Set("Subject", req.Subject)
    }
    if req.TargetArn != "" {
        val.Set("TargetArn", req.TargetArn)
    }
    if req.TopicArn != "" {
        val.Set("TopicArn", req.TopicArn)
    }
    // this is a guess... as i couldn't find an example of using message attributes
    ma := 1
    for k, v := range req.MessageAttributes {
        val.Set(fmt.Sprintf("MessageAttributes.%s.%d.DataType", k, ma), v.DataType)
        val.Set(fmt.Sprintf("MessageAttributes.%s.%d.StringValue", k, ma), v.StringValue)
        ma ++
    }

    req.CanonicalUri = req.CanonicalUri + "?" + val.Encode()
    return nil

}


func (req PublishRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    giResponse := new(PublishResponse)
    xml.Unmarshal(response, giResponse)
    if statusCode != 200 {
        return PublishError{
            Response: giResponse,
            Code: statusCode,
            RawResponse: string(response),
        }
    }
    return giResponse
}


func (gir PublishRequest) Request() (*PublishResponse, error) {
    request, err := awsgo.BuildEmptyContentRequest(&gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS2
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*PublishResponse), err
}

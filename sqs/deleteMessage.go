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

package sqs

import (
    "github.com/fromkeith/awsgo"
    "errors"
    "fmt"
    "net/url"
    "encoding/xml"
)

type DeleteMessageRequest struct {
    awsgo.RequestBuilder

    ReceiptHandle string
    TaskQueue string
}

type DeleteMessageResponse struct {
    ResponseMetadata awsgo.ResponseMetaData
}


func NewDeleteMessageRequest() *DeleteMessageRequest {
    req := new(DeleteMessageRequest)
    req.Host.Service = "sqs"
    req.Host.Region = ""
    req.Host.Domain = "amazonaws.com"
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    req.Headers = make(map[string]string)
    req.RequestMethod = "GET"
    req.CanonicalUri = "/"
    return req
}

func (gir * DeleteMessageRequest) VerifyInput() (error) {
    gir.Host.Service = "sqs"
    if len(gir.Host.Region) == 0 {
        return errors.New("Host.Region cannot be empty")
    }
    if len(gir.TaskQueue) == 0 {
        return errors.New("Task Queue cannot be empty")
    }

    gir.CanonicalUri = fmt.Sprintf("%s?Action=%s&Version=%s",
        gir.TaskQueue,
        url.QueryEscape("DeleteMessage"),
        url.QueryEscape("2012-11-05"),
    )
    if gir.ReceiptHandle == "" {
        return errors.New("ReceiptHandle cannot be empty")
    }
    gir.CanonicalUri = fmt.Sprintf("%s&ReceiptHandle=%s",
        gir.CanonicalUri,
        url.QueryEscape(gir.ReceiptHandle),
    )
    return nil
}

func (gir DeleteMessageRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := awsgo.CheckForErrorXml(response); err != nil {
        return err
    }
    giResponse := new(DeleteMessageResponse)
    //fmt.Println(string(response))
    err := xml.Unmarshal(response, giResponse)
    if err != nil {
        return err
    }
    return giResponse
}

func (gir DeleteMessageRequest) Request() (*DeleteMessageResponse, error) {    
    request, err := awsgo.BuildEmptyContentRequest(&gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*DeleteMessageResponse), err
}

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

const (
    SQS_RECEIVE_ATTR_NAME_All = "All"
    SQS_RECEIVE_ATTR_NAME_SenderId = "SenderId"
    SQS_RECEIVE_ATTR_NAME_SentTimestamp = "SentTimestamp"
    SQS_RECEIVE_ATTR_NAME_ReceiveCount = "ApproximateReceiveCount"
    SQS_RECEIVE_ATTR_NAME_FirstReceiveTime = "ApproximateFirstReceiveTimestamp"
)


type ReceiveMessageRequest struct {
    awsgo.RequestBuilder

    AttributeToGet              []string
    MessageAttributesToGet      []string
    MaxNumberOfMessages         int
    VisibilityTimeout           int
    WaitTimeSeconds             int
    TaskQueue                   string
}

type SqsAttributes struct {
    Name string
    Value string
}

type RecieveMessageResultMessage struct {
    MessageId string
    ReceiptHandle string
    MD5OfBody string
    Body string
    Attribute []SqsAttributes
    AttributeMap map[string]string
}
type ReceiveMessageResult struct {
    Message RecieveMessageResultMessage
}
type ReceiveMessageResponse struct {
    ReceiveMessageResult ReceiveMessageResult
    ResponseMetadata awsgo.ResponseMetaData
}


func NewReceiveMessageRequest() *ReceiveMessageRequest {
    req := new(ReceiveMessageRequest)
    req.MaxNumberOfMessages = 1
    req.VisibilityTimeout = -1
    req.WaitTimeSeconds = -1
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

func (gir * ReceiveMessageRequest) VerifyInput() (error) {
    gir.Host.Service = "sqs"
    if len(gir.Host.Region) == 0 {
        return errors.New("Host.Region cannot be empty")
    }
    if len(gir.TaskQueue) == 0 {
        return errors.New("Task Queue cannot be empty")
    }
    query := make(url.Values)
    query.Add("Action", "ReceiveMessage")
    query.Add("Version", "2012-11-05")

    if len(gir.AttributeToGet) == 1 {
        query.Add("AttributeName", gir.AttributeToGet[0])
    } else {
        for i := range gir.AttributeToGet {
            query.Add(fmt.Sprintf("AttributeName.member.%d", i + 1), gir.AttributeToGet[i])
        }
    }
    for i := range gir.MessageAttributesToGet {
        query.Add(fmt.Sprintf("MessageAttributeName.member.%d", i + 1), gir.MessageAttributesToGet[i])
    }

    if gir.MaxNumberOfMessages != -1 {
        query.Add("MaxNumberOfMessages", fmt.Sprintf("%d", gir.MaxNumberOfMessages))
    }
    if gir.VisibilityTimeout != -1 {
        query.Add("VisibilityTimeout", fmt.Sprintf("%d", gir.VisibilityTimeout))
    }
    if gir.WaitTimeSeconds != -1 {
        query.Add("WaitTimeSeconds", fmt.Sprintf("%d", gir.WaitTimeSeconds))
    }

    gir.CanonicalUri = fmt.Sprintf("%s?%s",
        gir.TaskQueue,
        query.Encode(),
    )
    return nil
}

func (gir ReceiveMessageRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := awsgo.CheckForErrorXml(response); err != nil {
        return err
    }
    giResponse := new(ReceiveMessageResponse)
    err := xml.Unmarshal(response, giResponse)
    if err != nil {
        return err
    }
    giResponse.ReceiveMessageResult.Message.AttributeMap = make(map[string]string)
    for i := range giResponse.ReceiveMessageResult.Message.Attribute {
        attr := giResponse.ReceiveMessageResult.Message.Attribute[i]
        giResponse.ReceiveMessageResult.Message.AttributeMap[attr.Name] = attr.Value
    }
    return giResponse
}

func (gir ReceiveMessageRequest) Request() (*ReceiveMessageResponse, error) {
    request, err := awsgo.BuildEmptyContentRequest(&gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*ReceiveMessageResponse), err
}

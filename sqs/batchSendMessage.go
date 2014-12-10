/*
 * Copyright (c) 2014, fromkeith
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



type SendMessageBatchRequestEntry struct {
    DelaySeconds            int64
    Id                      string
    MessageBody             string
}

type SendBatchMessageRequest struct {
    awsgo.RequestBuilder

    QueueUrl            string
    Entries        []SendMessageBatchRequestEntry
}

type SendMessageBatchResult struct {
    Id              string
    MessageId string
    MD5OfMessageBody string
}


type BatchResultError struct {
    Code            string
    Id              string
    Message         string
    SenderFault     bool
}

type SendBatchMessageResponse struct {
    SendMessageBatchResult      []SendMessageBatchResult            `xml:"SendMessageBatchResultEntry>member"`
    BatchResultError            []BatchResultError                  `xml:"BatchResultErrorEntry>member"`
    ResponseMetadata awsgo.ResponseMetaData
}


func NewSendBatchMessageRequest() *SendBatchMessageRequest {
    req := new(SendBatchMessageRequest)
    req.Host.Service = "sqs"
    req.Host.Region = ""
    req.Host.Domain = "amazonaws.com"
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    req.Headers = make(map[string]string)
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"
    return req
}

func (gir * SendBatchMessageRequest) VerifyInput() (error) {
    gir.Host.Service = "sqs"
    if len(gir.Host.Region) == 0 {
        return errors.New("Host.Region cannot be empty")
    }
    if len(gir.QueueUrl) == 0 {
        return errors.New("QueueUrl cannot be empty")
    }
    query := make(url.Values)
    query.Add("Action", "SendMessageBatch")
    query.Add("Version", "2012-11-05")

    for i := range gir.Entries {
        query.Add(fmt.Sprintf("SendMessageBatchRequestEntry.%d.Id", i + 1), gir.Entries[i].Id)
        query.Add(fmt.Sprintf("SendMessageBatchRequestEntry.%d.MessageBody", i + 1), gir.Entries[i].MessageBody)
        if gir.Entries[i].DelaySeconds > 0 {
            query.Add(fmt.Sprintf("SendMessageBatchRequestEntry.%d.DelaySeconds", i + 1), fmt.Sprintf("%d", gir.Entries[i].DelaySeconds))
        }
    }

    gir.CanonicalUri = fmt.Sprintf("%s?%s",
        gir.QueueUrl,
        query.Encode(),
    )
    return nil
}

func (gir SendBatchMessageRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := awsgo.CheckForErrorXml(response); err != nil {
        return err
    }
    giResponse := new(SendBatchMessageResponse)
    //fmt.Println(string(response))
    err := xml.Unmarshal(response, giResponse)
    if err != nil {
        return err
    }
    return giResponse
}

func (gir SendBatchMessageRequest) Request() (*SendBatchMessageResponse, error) {
    request, err := awsgo.BuildEmptyContentRequest(&gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*SendBatchMessageResponse), err
}
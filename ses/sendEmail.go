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
package ses

import (
    "github.com/fromkeith/awsgo"
    //"errors"
    "fmt"
    "net/url"
    "encoding/xml"
)


type EmailAddresses struct {
    BccAddresses    []string
    CcAddresses     []string 
    ToAddresses     []string
}

type MessageContent struct {
    Charset     string
    Data        string
}

type MessageBody struct {
    Html    MessageContent
    Text    MessageContent
}

type Message struct {
    Body    MessageBody
    Subject MessageContent
}

type SendEmailRequest struct {
    awsgo.RequestBuilder

    Destination EmailAddresses
    Message     Message
    ReturnPath  string
    Source      string
}

type SendEmailResult struct {
    MessageId string
}

type SendEmailResponse struct {
    SendEmailResult SendEmailResult
    ResponseMetadata awsgo.ResponseMetaData
}


func NewSendEmailRequest() *SendEmailRequest {
    req := new(SendEmailRequest)
    req.Host.Service = "ses"
    req.Host.Region = ""
    req.Host.Domain = ""
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    req.Headers = make(map[string]string)
    req.RequestMethod = "GET"
    req.CanonicalUri = "/"
    return req
}

func addToContentToUri(uri string, content MessageContent, prefix string) string {
    if content.Data != "" {
        uri = fmt.Sprintf("%s&%s.Data=%s",
            uri, prefix, url.QueryEscape(content.Data))
    }
    if content.Charset != "" {
        uri = fmt.Sprintf("%s&%s.Charset=%s",
            uri, prefix, url.QueryEscape(content.Charset))
    }
    return uri
}
func addStringList(uri string, item []string, prefix string) string {
    if item != nil && len(item) > 0 {
        for i := range item {
            uri = fmt.Sprintf("%s&%s.member.%d=%s",
                uri, prefix, i + 1, url.QueryEscape(item[i]))
        }
    }
    return uri
}

func (gir * SendEmailRequest) VerifyInput() (error) {
    gir.Host.Service = "email"
    gir.Host.Region = "us-east-1"

    gir.CanonicalUri = "/?Action=SendEmail&Version=2010-12-01"

    gir.CanonicalUri = addToContentToUri(gir.CanonicalUri, gir.Message.Body.Html, "Message.Body.Html")
    gir.CanonicalUri = addToContentToUri(gir.CanonicalUri, gir.Message.Body.Text, "Message.Body.Text")
    gir.CanonicalUri = addToContentToUri(gir.CanonicalUri, gir.Message.Subject, "Message.Subject")
    if gir.Source != "" {
        gir.CanonicalUri = fmt.Sprintf("%s&Source=%s", gir.CanonicalUri, url.QueryEscape(gir.Source))
    }
    if gir.ReturnPath != "" {
        gir.CanonicalUri = fmt.Sprintf("%s&ReturnPath=%s", gir.CanonicalUri, url.QueryEscape(gir.ReturnPath))
    }
    gir.CanonicalUri = addStringList(gir.CanonicalUri, gir.Destination.ToAddresses, "Destination.ToAddresses")
    gir.CanonicalUri = addStringList(gir.CanonicalUri, gir.Destination.BccAddresses, "Destination.BccAddresses")
    gir.CanonicalUri = addStringList(gir.CanonicalUri, gir.Destination.CcAddresses, "Destination.CcAddresses")

    return gir.RequestBuilder.VerifyInput()
}

func (gir SendEmailRequest) DeMarshalGetItemResponse(response []byte, headers map[string]string) (interface{}) {
    giResponse := new(SendEmailResponse)
    //fmt.Println(string(response))
    xml.Unmarshal(response, giResponse)
    //json.Unmarshal([]byte(response), giResponse)
    return giResponse
}

func (gir SendEmailRequest) Request() (*SendEmailResponse, error) {    
    request, err := awsgo.BuildEmptyContentRequest(&gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := awsgo.DoRequest(&gir, request)
    if resp == nil {
        return nil, err
    }
    return resp.(*SendEmailResponse), err
}

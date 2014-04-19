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

package cloudsearch


import (
    "github.com/fromkeith/awsgo"
    "errors"
    "encoding/json"
    "fmt"
)

var (
    Verification_Error_EndpointEmpty = errors.New("Endpoint can not be empty!")
    NoPermission_For_Endpoint = errors.New("403 Invalid permissions to hit endpoint.")
)


type DocumentInfo struct {
    // 'add' or 'delete'
    Type            string      `json:"type"`
    Id              string      `json:"id"`
    Fields          interface{} `json:"fields,omitempty"`
}


type BatchDocumentRequest struct {
    awsgo.RequestBuilder

    Items           []DocumentInfo
    // if you full endpoint is blah.us-west-2.cloudsearch.amazonaws.com
    // then set this to 'blah'
    Endpoint        string      `json:"-"`
}

type BatchDocumentMessage struct {
    Message         string  `json:"message"`
}

type BatchDocumentResponse struct {
    Status          string      `json:"status"`
    Adds            int         `json:"adds"`
    Deletes         int         `json:"deletes"`
    Errors          []BatchDocumentMessage `json:"errors"`
    Warnings        []BatchDocumentMessage `json:"warnings"`
    StatusCode      int
}


// Creates a new BatchDocumentRequest, populating in some defaults
func NewBatchDocumentRequest() *BatchDocumentRequest {
    req := new(BatchDocumentRequest)
    req.Host.Service = "cloudsearch"
    req.Host.Region = ""
    req.Host.Domain = "amazonaws.com"
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    req.Headers = make(map[string]string)
    req.Headers["Content-Type"] = "application/json"
    req.RequestMethod = "POST"
    req.CanonicalUri = "/2013-01-01/documents/batch"
    return req
}


func (gir * BatchDocumentRequest) VerifyInput() (error) {
    if len(gir.Endpoint) == 0 {
        return Verification_Error_EndpointEmpty
    }
    gir.Host.Override = fmt.Sprintf("%s.%s.%s.%s", gir.Endpoint, gir.Host.Region, gir.Host.Service, gir.Host.Domain)
    return nil
}

func (gir BatchDocumentRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if statusCode == 403 {
        return NoPermission_For_Endpoint
    }
    resp := new(BatchDocumentResponse)
    err := json.Unmarshal(response, resp)
    if err != nil {
        return err
    }
    resp.StatusCode = statusCode
    return resp
}


func (gir BatchDocumentRequest) Request() (*BatchDocumentResponse, error) {
    request, err := awsgo.NewAwsRequest(&gir, gir.Items)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*BatchDocumentResponse), err
}

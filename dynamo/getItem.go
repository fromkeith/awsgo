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
package dynamo

import (
    "github.com/fromkeith/awsgo"
    "errors"
    "encoding/json"
)


type GetItemRequest struct {
    awsgo.RequestBuilder

    AttributesToGet        []string  `json:",omitempty"`
    ConsistentRead         bool     `json:",string"`
    Search                 map[string]interface{}  `json:"Key"` 
    TableName              string
    ReturnConsumedCapacity string     
}

type GetItemResponse struct {
    ConsumedCapacity *CapacityResult             `json:",omitempty"`
    Item             map[string]interface{}      `json:"-"`
    RawItem          map[string]map[string]interface{}  `json:"Item"`
}

type GetItemResponseFuture struct {
    response chan *GetItemResponse
    errResponse chan error
}

func (f * GetItemResponseFuture) Get() (*GetItemResponse, error) {
    select {
    case err := <- f.errResponse:
        return nil, err
    case resp := <- f.response:
        return resp, nil
    }
}

func NewGetItemRequest() *GetItemRequest {
    req := new(GetItemRequest)
    req.AttributesToGet = nil
    req.ConsistentRead = false
    req.Search = make(map[string]interface{})
    req.TableName = ""
    req.ReturnConsumedCapacity = ConsumedCapacity_NONE
    req.Host.Service = "dynamodb"
    req.Host.Region = ""
    req.Host.Domain = ""
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = GetItemTarget
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"
    return req
}

func (gir * GetItemRequest) VerifyInput() (error) {
    gir.Host.Service = "dynamodb"
    if len(gir.TableName) == 0 {
        return errors.New("TableName cannot be empty")
    }
    if len(gir.Search) == 0 {
        return errors.New("Search parameters cannot be empty")
    }
    if len(gir.Host.Region) == 0 {
        return errors.New("Host.Region cannot be empty")
    }
    // repair any errors, like if you put a string, instead of an awsgo String item
    for k, v := range(gir.Search) {
        gir.Search[k] = awsgo.ConvertToAwsItem(v)
    }
    return gir.RequestBuilder.VerifyInput()
}

func (gir GetItemRequest) CoRequest() (*GetItemResponseFuture, error) {
    request, err := awsgo.BuildRequest(&gir, gir)
    if err != nil {
        return nil, err
    }
    future := new(GetItemResponseFuture)
    future.errResponse = make(chan error)
    future.response = make(chan * GetItemResponse)
    go gir.CoDoRequest(request, future)
    return future, nil
}

func (gir GetItemRequest) DeMarshalGetItemResponse(response []byte, headers map[string]string) (interface{}) {
    if err := CheckForErrorResponse(response); err != nil {
        return err
    }
    giResponse := new(GetItemResponse)
    err := json.Unmarshal([]byte(response), giResponse)
    if err != nil {
        return err
    }
    giResponse.Item = make(map[string]interface{})
    awsgo.FromRawMapToEasyTypedMap(giResponse.RawItem, giResponse.Item)
    return giResponse
}

func (gir GetItemRequest) CoDoRequest(request awsgo.AwsRequest, future * GetItemResponseFuture) {
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := awsgo.DoRequest(&gir, request)
    if err != nil {
        future.errResponse <- err
    } else {
        future.response <- resp.(*GetItemResponse)
    }
    close(future.errResponse)
    close(future.response)
}

func (gir GetItemRequest) Request() (*GetItemResponse, error) {    
    request, err := awsgo.BuildRequest(&gir, gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := awsgo.DoRequest(&gir, request)
    if resp == nil {
        return nil, err
    }
    return resp.(*GetItemResponse), err
}

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

type BatchGetItemRequestTable struct {
    AttributesToGet        []string  `json:",omitempty"`
    ConsistentRead         bool     `json:",string"`
    Search                 []map[string]interface{}  `json:"Keys"` 
}

type batchGetItemRequestTableDeserialized struct {
    AttributesToGet        []string  `json:",omitempty"`
    ConsistentRead         bool     `json:",string"`
    Search                 []map[string]map[string]interface{}  `json:"Keys"` 
}


type BatchGetItemRequest struct {
    awsgo.RequestBuilder

    RequestItems            map[string]BatchGetItemRequestTable
    ReturnConsumedCapacity  string
}

type BatchGetItemResponse struct {
    ConsumedCapacity *CapacityResult             `json:",omitempty"`
    Responses        map[string][]map[string]interface{}        `json:"-"`
    RawResponses     map[string][]map[string]map[string]interface{}  `json:"Responses"`
    UnprocessedKeys  map[string]BatchGetItemRequestTable        `json:"-"`
    RawUnprocessedKeys map[string]batchGetItemRequestTableDeserialized `json:"UnprocessedKeys"`
}


func NewBatchGetIteamRequestTable() (ret BatchGetItemRequestTable) {
    ret.AttributesToGet = nil
    ret.ConsistentRead = false
    ret.Search = nil
    return
}

func NewBatchGetItemRequest() *BatchGetItemRequest {
    req := new(BatchGetItemRequest)
    req.RequestItems = make(map[string]BatchGetItemRequestTable)
    req.ReturnConsumedCapacity = ConsumedCapacity_NONE
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = BatchGetItemTarget
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"
    req.Host.Service = ""
    req.Host.Region = ""
    req.Host.Domain = ""
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    return req
}

func (gir * BatchGetItemRequest) VerifyInput() (error) {
    gir.Host.Service = "dynamodb"
    if len(gir.RequestItems) == 0 {
        return errors.New("RequestItems cannot be empty")
    }
    if len(gir.Host.Region) == 0 {
        return errors.New("Host.Region cannot be empty")
    }
    return gir.RequestBuilder.VerifyInput()
}
/*
func (gir BatchGetItemRequest) CoRequest() (*GetItemResponseFuture, error) {
    request, err := awsgo.BuildRequest(&gir, GetItemTarget, gir)
    if err != nil {
        return nil, err
    }
    future := new(GetItemResponseFuture)
    future.errResponse = make(chan error)
    future.response = make(chan * GetItemResponse)
    go gir.CoDoRequest(request, future)
    return future, nil
}*/

func (gir BatchGetItemRequest) DeMarshalGetItemResponse(response []byte, headers map[string]string) (interface{}) {
    if err := CheckForErrorResponse(response); err != nil {
        return err
    }
    giResponse := new(BatchGetItemResponse)
    err := json.Unmarshal([]byte(response), giResponse)
    if err != nil {
        return err
    }
    giResponse.Responses = make(map[string][]map[string]interface{})
    giResponse.UnprocessedKeys = make(map[string]BatchGetItemRequestTable)
    //RawResponses: map[string][]map[string]map[string]string
    for key, val := range giResponse.RawResponses {
        giResponse.Responses[key] = make([]map[string]interface{}, len(val))
        for i := range val {
            giResponse.Responses[key][i] = make(map[string]interface{})
            awsgo.FromRawMapToAwsItemMap(val[i], giResponse.Responses[key][i])
        }
    }
    // RawUnprocessed, including Search: map[string]{Search}[]map[string]map[string]string
    for key, val := range giResponse.RawUnprocessedKeys {
        var c BatchGetItemRequestTable
        c.AttributesToGet = val.AttributesToGet
        c.ConsistentRead = val.ConsistentRead
        c.Search = make([]map[string]interface{}, len(val.Search))
        for i := range val.Search {
            c.Search[i] = make(map[string]interface{})
            awsgo.FromRawMapToAwsItemMap(val.Search[i], c.Search[i])
        }
        giResponse.UnprocessedKeys[key] = c
    }
    return giResponse
}
/*
func (gir BatchGetItemRequest) CoDoRequest(request awsgo.AwsRequest, future * GetItemResponseFuture) {
    resp, err := awsgo.DoRequest(&gir, request)
    if err != nil {
        future.errResponse <- err
    } else {
        future.response <- resp.(*GetItemResponse)
    }
    close(future.errResponse)
    close(future.response)
}*/

func (gir BatchGetItemRequest) Request() (*BatchGetItemResponse, error) {    
    request, err := awsgo.BuildRequest(&gir, gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := awsgo.DoRequest(&gir, request)
    if resp == nil {
        return nil, err
    }
    return resp.(*BatchGetItemResponse), err
}
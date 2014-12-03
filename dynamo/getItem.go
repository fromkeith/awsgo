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

var (
    Verification_Error_TableNameEmpty = errors.New("TableName cannot be empty")
    Verification_Error_SearchEmpty = errors.New("Search parameters cannot be empty")
    Verification_Error_RegionEmpty = errors.New("Host.Region cannot be empty")
    Verification_Error_ServiceEmpty = errors.New("Host.Service cannot be empty")
)


type GetItemRequest struct {
    awsgo.RequestBuilder

    AttributesToGet        []string  `json:",omitempty"`
    ConsistentRead         bool
    // The key to search for
    Search                 map[string]interface{}  `json:"Key"`
    TableName              string
    ReturnConsumedCapacity string
}

type GetItemResponse struct {
    // if return consumed capacity was set, this will populate with the result
    ConsumedCapacity *CapacityResult             `json:"ConsumedCapacity,omitempty"`
    // the item response, with easily castable values
    Item             map[string]interface{}      `json:"-"`
    // the raw item response from the wire
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

// Creates a new GetItemRequest, populating in some defaults
func NewGetItemRequest() *GetItemRequest {
    req := new(GetItemRequest)
    req.AttributesToGet = nil
    req.ConsistentRead = false
    req.Search = make(map[string]interface{})
    req.TableName = ""
    req.ReturnConsumedCapacity = ConsumedCapacity_NONE
    req.Host.Service = "dynamodb"
    req.Host.Region = ""
    req.Host.Domain = "amazonaws.com"
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = GetItemTarget
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"
    return req
}

func (gir * GetItemRequest) VerifyInput() (error) {
    if len(gir.Host.Service) == 0 {
        return Verification_Error_ServiceEmpty
    }
    if len(gir.TableName) == 0 {
        return Verification_Error_TableNameEmpty
    }
    if len(gir.Search) == 0 {
        return Verification_Error_SearchEmpty
    }
    if len(gir.Host.Region) == 0 {
        return Verification_Error_RegionEmpty
    }
    // repair any errors, like if you put a string, instead of an awsgo String item
    for k, v := range(gir.Search) {
        gir.Search[k] = awsgo.ConvertToAwsItem(v)
    }
    return nil
}

func (gir GetItemRequest) CoRequest() (*GetItemResponseFuture, error) {
    request, err := awsgo.NewAwsRequest(&gir, gir)
    if err != nil {
        return nil, err
    }
    future := new(GetItemResponseFuture)
    future.errResponse = make(chan error)
    future.response = make(chan * GetItemResponse)
    go gir.CoDoAndDemarshall(request, future)
    return future, nil
}

func (gir GetItemRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := CheckForErrorResponse(response, statusCode); err != nil {
        return err
    }
    giResponse := new(GetItemResponse)
    err := json.Unmarshal(response, giResponse)
    if err != nil {
        newErr := &awsgo.UnmarhsallingError {
            ActualContent : string(response),
            MarshallError : err,
        }
        return newErr
    }
    giResponse.Item = make(map[string]interface{})
    awsgo.FromRawMapToEasyTypedMap(giResponse.RawItem, giResponse.Item)
    return giResponse
}

func (gir GetItemRequest) CoDoAndDemarshall(request awsgo.AwsRequest, future * GetItemResponseFuture) {
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&gir)
    if err != nil {
        future.errResponse <- err
    } else {
        future.response <- resp.(*GetItemResponse)
    }
    close(future.errResponse)
    close(future.response)
}

func (gir GetItemRequest) Request() (*GetItemResponse, error) {
    request, err := awsgo.NewAwsRequest(&gir, gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*GetItemResponse), err
}

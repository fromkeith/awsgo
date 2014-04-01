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

// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_PutItem.html
type PutItemRequest struct {
    awsgo.RequestBuilder

    Expected                map[string]ExpectedItem  `json:",omitempty"`
    Item                    map[string]interface{}  `json:"Item"`
    TableName               string
    ReturnConsumedCapacity  string
    ReturnItemCollectionMetrics  string
    ReturnValues            string
}


type PutItemResponse struct {
    RawBeforeAttributes      map[string]map[string]interface{}   `json:"Attributes"`
    BeforeAttributes         map[string]interface{}         `json:"-"`
    ConsumedCapacity         *CapacityResult                `json:",omitempty"`
    ItemCollectionMetrics    *ItemCollectionMetricsStruct   `json:",omitempty"`
}

type PutItemResponseFuture struct {
    response chan *PutItemResponse
    errResponse chan error
}

func (f * PutItemResponseFuture) Get() (*PutItemResponse, error) {
    select {
    case err := <- f.errResponse:
        return nil, err
    case resp := <- f.response:
        return resp, nil
    }
}

func NewPutItemRequest() *PutItemRequest {
    req := new(PutItemRequest)
    req.Item = make(map[string]interface{})
    req.Expected = make(map[string]ExpectedItem)
    req.TableName = ""
    req.ReturnConsumedCapacity = ConsumedCapacity_NONE
    req.ReturnItemCollectionMetrics = ReturnItemCollection_NONE
    req.ReturnValues = ReturnValues_NONE
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = PutItemTarget
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"

    req.Host.Service = "dynamodb"
    req.Host.Region = ""
    req.Host.Domain = ""
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    return req
}

func (pir * PutItemRequest) VerifyInput() (error) {
    if len(pir.TableName) == 0 {
        return errors.New("TableName cannot be empty")
    }
    if len(pir.Item) == 0 {
        return errors.New("Item cannot be empty")
    }
    if len(pir.Host.Region) == 0 {
        return errors.New("Host.Region cannot be empty")
    }
    for k, v := range pir.Item {
        pir.Item[k] = awsgo.ConvertToAwsItem(v)
    }
    for k, v := range pir.Expected {
        if v.Exists {
            pir.Expected[k] = ExpectedItem{v.Exists, awsgo.ConvertToAwsItem(v.Value)}
        }
    }
    return nil
}
func (pir PutItemRequest) CoRequest() (*PutItemResponseFuture, error) {
    request, err := awsgo.NewAwsRequest(&pir, pir)
    if err != nil {
        return nil, err
    }
    future := new(PutItemResponseFuture)
    future.errResponse = make(chan error)
    future.response = make(chan * PutItemResponse)
    go pir.CoDoAndDemarshall(request, future)
    return future, nil
}

func (pir PutItemRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := CheckForErrorResponse(response, statusCode); err != nil {
        return err
    }
    piResponse := new(PutItemResponse)
    err := json.Unmarshal([]byte(response), piResponse)
    if err != nil {
        return err
    }
    if len(piResponse.RawBeforeAttributes) > 0 {
        piResponse.BeforeAttributes = make(map[string]interface{})
        awsgo.FromRawMapToEasyTypedMap(piResponse.RawBeforeAttributes, piResponse.BeforeAttributes)
    }
    if piResponse.ItemCollectionMetrics != nil {
        piResponse.ItemCollectionMetrics.ItemCollectionKey = make(map[string]interface{})
        awsgo.FromRawMapToEasyTypedMap(piResponse.ItemCollectionMetrics.RawItemCollectionKey, piResponse.ItemCollectionMetrics.ItemCollectionKey)
    }
    return piResponse
}


func (pir PutItemRequest) CoDoAndDemarshall(request awsgo.AwsRequest, future * PutItemResponseFuture) {
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&pir)
    if err != nil {
        future.errResponse <- err
    } else {
        future.response <- resp.(*PutItemResponse)
    }
    close(future.errResponse)
    close(future.response)
}

func (pir PutItemRequest) Request() (*PutItemResponse, error) {
    request, err := awsgo.NewAwsRequest(&pir, pir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&pir)
    if resp == nil {
        return nil, err
    }
    return resp.(*PutItemResponse), err
}

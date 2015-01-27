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

const (
    AttributeUpdate_Action_Add = "ADD"
    AttributeUpdate_Action_Put = "PUT"
    AttributeUpdate_Action_Delete = "DELETE"
)

type AttributeUpdates struct {
    Action                  string
    Value                   interface{}
}

// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_UpdateItem.html
type UpdateItemRequest struct {
    awsgo.RequestBuilder

    Expected                map[string]ExpectedItem  `json:",omitempty"`
    UpdateKey               map[string]interface{}  `json:"Key"`
    Update                  map[string]AttributeUpdates  `json:"AttributeUpdates"`
    TableName               string
    ReturnConsumedCapacity  string
    ReturnItemCollection    string
    ReturnValues            string
}


type UpdateItemResponse struct {
    RawBeforeAttributes      map[string]map[string]interface{}   `json:"Attributes"`
    BeforeAttributes         map[string]interface{}         `json:"-"`
    ConsumedCapacity         *CapacityResult                `json:",omitempty"`
    ItemCollectionMetrics    *ItemCollectionMetricsStruct   `json:",omitempty"`
}

type UpdateItemResponseFuture struct {
    response chan *UpdateItemResponse
    errResponse chan error
}

func (f * UpdateItemResponseFuture) Get() (*UpdateItemResponse, error) {
    select {
    case err := <- f.errResponse:
        return nil, err
    case resp := <- f.response:
        return resp, nil
    }
}

func NewUpdateItemRequest() *UpdateItemRequest {
    req := new(UpdateItemRequest)
    req.UpdateKey = make(map[string]interface{})
    req.Expected = make(map[string]ExpectedItem)
    req.Update = make(map[string]AttributeUpdates)
    req.TableName = ""
    req.ReturnConsumedCapacity = ConsumedCapacity_NONE
    req.ReturnItemCollection = ReturnItemCollection_NONE
    req.ReturnValues = ReturnValues_NONE
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = UpdateItemTarget
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"

    req.Host.Service = ""
    req.Host.Region = ""
    req.Host.Domain = "amazonaws.com"
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    return req
}

func (pir * UpdateItemRequest) VerifyInput() (error) {
    pir.Host.Service = "dynamodb"
    if len(pir.TableName) == 0 {
        return errors.New("TableName cannot be empty")
    }
    if len(pir.UpdateKey) == 0 {
        return errors.New("UpdateKey cannot be empty")
    }
    if len(pir.Update) == 0 {
        return errors.New("Update cannot be empty")
    }
    if len(pir.Host.Region) == 0 {
        return errors.New("Host.Region cannot be empty")
    }
    for k, v := range pir.UpdateKey {
        pir.UpdateKey[k] = awsgo.ConvertToAwsItem(v)
    }
    for k, v := range pir.Update {
        if v.Action != AttributeUpdate_Action_Delete {
            pir.Update[k] = AttributeUpdates{v.Action, awsgo.ConvertToAwsItem(v.Value)}
        } else if v.Value != nil {
            pir.Update[k] = AttributeUpdates{v.Action, awsgo.ConvertToAwsItem(v.Value)}
        }
    }
    for k, v := range pir.Expected {
        if v.Exists {
            pir.Expected[k] = ExpectedItem{v.Exists, awsgo.ConvertToAwsItem(v.Value)}
        }
    }
    return nil
}
func (pir UpdateItemRequest) CoRequest() (*UpdateItemResponseFuture, error) {
    request, err := awsgo.NewAwsRequest(&pir, pir)
    if err != nil {
        return nil, err
    }
    future := new(UpdateItemResponseFuture)
    future.errResponse = make(chan error)
    future.response = make(chan * UpdateItemResponse)
    go pir.CoDoAndDemarshall(request, future)
    return future, nil
}

func (pir UpdateItemRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := CheckForErrorResponse(response, statusCode); err != nil {
        return err
    }
    piResponse := new(UpdateItemResponse)
    err := json.Unmarshal(response, piResponse)
    if err != nil {
        return err
    }
    if len(piResponse.RawBeforeAttributes) > 0 {
        piResponse.BeforeAttributes = make(map[string]interface{})
        awsgo.FromRawMapToEasyTypedMap(piResponse.RawBeforeAttributes, piResponse.BeforeAttributes)
    }
    return piResponse
}


func (pir UpdateItemRequest) CoDoAndDemarshall(request awsgo.AwsRequest, future * UpdateItemResponseFuture) {
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&pir)
    if err != nil {
        future.errResponse <- err
    } else {
        future.response <- resp.(*UpdateItemResponse)
    }
    close(future.errResponse)
    close(future.response)
}

func (pir UpdateItemRequest) Request() (*UpdateItemResponse, error) {
    request, err := awsgo.NewAwsRequest(&pir, pir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&pir)
    if resp == nil {
        return nil, err
    }
    return resp.(*UpdateItemResponse), err
}

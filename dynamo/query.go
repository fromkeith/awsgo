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
    "fmt"
)

type KeyConditions struct {
    AttributeValueList      []interface{}
    ComparisonOperator      string
}

type QueryRequest struct {
    awsgo.RequestBuilder

    AttributesToGet         []string    `json:",omitempty"`
    ConsistentRead          bool        `json:",string"`
    ExclusiveStartKey       map[string]interface{} `json:",omitempty"`
    IndexName               string      `json:",omitempty"`
    KeyConditions           map[string]KeyConditions `json:",omitempty"`
    Limit                   float64     `json:",omitempty"`
    ReturnConsumedCapacity  string      `json:",omitempty"`
    ScanIndexForward        *bool       `json:",omitempty"`
    Select                  string      `json:",omitempty"`
    TableName               string
}

type QueryResponse struct {
    ConsumedCapacity *CapacityResult             `json:",omitempty"`
    Count                   float64
    Items                   []map[string]interface{}    `json:"-"`
    RawItems                []map[string]map[string]interface{} `json:"Items"`
    LastEvaluatedKey        map[string]interface{}      `json:"-"`
    RawLastEvaluatedKey     map[string]map[string]interface{}    `json:"LastEvaluatedKey"`
}

func (req *QueryRequest) AddKeyCondition(keyName string, values []interface{}, operator string) {
    if req.KeyConditions == nil {
        req.KeyConditions = make(map[string]KeyConditions)
    }
    var condition KeyConditions
    condition.AttributeValueList = values
    condition.ComparisonOperator = operator
    req.KeyConditions[keyName] = condition
}


func NewQueryRequest() *QueryRequest {
    req := new(QueryRequest)
    req.AttributesToGet = nil
    req.ConsistentRead = false
    req.ExclusiveStartKey = nil
    req.IndexName = ""
    req.KeyConditions = nil
    req.ReturnConsumedCapacity = ConsumedCapacity_NONE
    req.ScanIndexForward = nil
    req.Select = ""
    req.TableName = ""
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = QueryTarget
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"
    return req
}

func (gir * QueryRequest) VerifyInput() (error) {
    gir.Host.Service = "dynamodb"
    if len(gir.TableName) == 0 {
        return errors.New("TableName cannot be empty")
    }
    for _, condition := range gir.KeyConditions {
        for i := range condition.AttributeValueList {
            condition.AttributeValueList[i] = awsgo.ConvertToAwsItem(condition.AttributeValueList[i])
        }
    }
    return gir.RequestBuilder.VerifyInput()
}

func (gir QueryRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := CheckForErrorResponse(response, statusCode); err != nil {
        return err
    }
    giResponse := new(QueryResponse)
    err := json.Unmarshal([]byte(response), giResponse)
    if err != nil {
        fmt.Println("Error unmarshalling query response!", string(response))
        return err
    }

    giResponse.Items = make([]map[string]interface{}, len(giResponse.RawItems))
    for i := range giResponse.RawItems {
        giResponse.Items[i] = make(map[string]interface{})
        awsgo.FromRawMapToEasyTypedMap(giResponse.RawItems[i], giResponse.Items[i])
    }
    giResponse.LastEvaluatedKey = make(map[string]interface{})
    awsgo.FromRawMapToEasyTypedMap(giResponse.RawLastEvaluatedKey, giResponse.LastEvaluatedKey)
    return giResponse
}

func (gir QueryRequest) Request() (*QueryResponse, error) {
    request, err := awsgo.BuildRequest(&gir, gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := awsgo.DoRequest(&gir, request)
    if resp == nil {
        return nil, err
    }
    return resp.(*QueryResponse), err
}


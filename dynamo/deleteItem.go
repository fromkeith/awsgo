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
    Verification_Error_DeleteKeyEmpty = errors.New("DeleteKey cannot be empty")
)


type DeleteItemRequest struct {
    awsgo.RequestBuilder

    Expected                 map[string]ExpectedItem  `json:",omitempty"`
    DeleteKey                map[string]interface{}  `json:"Key,omitempty"`
    ReturnConsumedCapacity   string
    ReturnItemCollectionMetrics    string
    ReturnValues             string
    TableName                string
}

type DeleteItemResponse struct {
    ConsumedCapacity *CapacityResult             `json:"ConsumedCapacity,omitempty"`
    Attributes             map[string]interface{}      `json:"-"`
    RawAttributes          map[string]map[string]interface{}  `json:"Attributes"`
}

func NewDeleteItemRequest() *DeleteItemRequest {
    req := new(DeleteItemRequest)

    req.DeleteKey = make(map[string]interface{})
    req.Expected = make(map[string]ExpectedItem)
    req.TableName = ""
    req.ReturnConsumedCapacity = ConsumedCapacity_NONE
    req.ReturnItemCollectionMetrics = ReturnItemCollection_NONE
    req.ReturnValues = ReturnValues_NONE
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = DeleteItemTarget
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"

    req.Host.Service = "dynamodb"
    req.Host.Region = ""
    req.Host.Domain = ""
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    return req
}

func (gir * DeleteItemRequest) VerifyInput() (error) {
    if len(gir.Host.Service) == 0 {
        return Verification_Error_ServiceEmpty
    }
    if len(gir.TableName) == 0 {
        return Verification_Error_TableNameEmpty
    }
    if len(gir.DeleteKey) == 0 {
        return Verification_Error_DeleteKeyEmpty
    }
    // repair any errors, like if you put a string, instead of an awsgo String item
    for k, v := range(gir.DeleteKey) {
        gir.DeleteKey[k] = awsgo.ConvertToAwsItem(v)
    }
    for k, v := range gir.Expected {
        gir.Expected[k] = ExpectedItem{v.Exists, awsgo.ConvertToAwsItem(v.Value)}
    }
    return gir.RequestBuilder.VerifyInput()
}


func (gir DeleteItemRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := CheckForErrorResponse(response, statusCode); err != nil {
        return err
    }
    giResponse := new(DeleteItemResponse)
    err := json.Unmarshal(response, giResponse)
    if err != nil {
        newErr := &awsgo.UnmarhsallingError {
            ActualContent : string(response),
            MarshallError : err,
        }
        return newErr
    }
    giResponse.Attributes = make(map[string]interface{})
    awsgo.FromRawMapToEasyTypedMap(giResponse.RawAttributes, giResponse.Attributes)
    return giResponse
}

func (gir DeleteItemRequest) Request() (*DeleteItemResponse, error) {
    request, err := awsgo.BuildRequest(&gir, gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := awsgo.DoRequest(&gir, request)
    if resp == nil {
        return nil, err
    }
    return resp.(*DeleteItemResponse), err
}
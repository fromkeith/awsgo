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
    "encoding/json"
    "fmt"
)

const (
    ConsumedCapacity_TOTAL = "TOTAL"
    ConsumedCapacity_NONE = "NONE"
    ReturnItemCollection_SIZE = "SIZE"
    ReturnItemCollection_NONE = "NONE"
    ReturnValues_ALL_OLD = "ALL_OLD"
    ReturnValues_NONE = "NONE"
)


// Targets
const (
    GetItemTarget = "DynamoDB_20120810.GetItem"
    PutItemTarget = "DynamoDB_20120810.PutItem"
    BatchGetItemTarget = "DynamoDB_20120810.BatchGetItem"
    UpdateItemTarget = "DynamoDB_20120810.UpdateItem"
)
// Known Errors
const (
    ConditionalCheckFailed = "com.amazonaws.dynamodb.v20120810#ConditionalCheckFailedException"
    SerializationException = "com.amazon.coral.service#SerializationException"
)


type CapacityResult struct {
    CapacityUnits int           `json:",string"`
    TableName     string
}

type ExpectedItem struct {
    Exists  bool        `json:",string"`
    Value   interface{}
}


type ItemCollectionMetricsStruct struct {
    RawItemCollectionKey        map[string]map[string]string   `json:"ItemCollectionKey"`
    ItemCollectionKey           map[string]interface{}  `json:"-"`
    SizeEstimateRangeGB         []float64         `json:",string"`
}

type ErrorResult struct {
    Type        string  `json:"__type"`
    Message     string  `json:"message"`
}

func (e * ErrorResult) Error() string {
    return fmt.Sprintf("%s : %s", e.Type, e.Message)
}

func CheckForErrorResponse(response []byte) error {
    errorResult := new(ErrorResult)
    err2 := json.Unmarshal([]byte(response), errorResult)
    if err2 == nil {
        if errorResult.Type != "" {
            return errorResult
        }
    }
    return nil
}
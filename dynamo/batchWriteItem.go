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

type BatchWriteItemDeleteRequest struct {
    Key         map[string]interface{}
}

type BatchWriteItemPutRequest struct {
    Item         map[string]interface{}
}

type BatchWriteItem struct {
    DeleteRequest *BatchWriteItemDeleteRequest `json:",omitempty"`
    PutRequest *BatchWriteItemPutRequest `json:",omitempty"`
}

type BatchWriteItemRequest struct {
    awsgo.RequestBuilder

    RequestItems            map[string][]BatchWriteItem
    ReturnConsumedCapacity  string
    ReturnItemCollectionMetrics string
}

type BatchWriteItemResponse struct {
    ConsumedCapacity *CapacityResult             `json:",omitempty"`
    ItemCollectionMetrics * ItemCollectionMetricsStruct `json:",omitempty"`
    UnprocessedItems    map[string][]BatchWriteItem `json:"-"`
                    //     table  | operations| key/item | name     | item type/value
    RawUnprocessItems  map[string][]map[string]map[string]map[string]map[string]interface{} `json:"UnprocessedItems"`
}


func (req *BatchWriteItemRequest) AddTable(name string) {
    req.RequestItems[name] = make([]BatchWriteItem, 25)[0:0]
}

func (req * BatchWriteItemRequest) checkAndGrow(table string) ([]BatchWriteItem) {
    var genericItems []BatchWriteItem
    var ok bool
    if genericItems, ok = req.RequestItems[table]; !ok {
        req.AddTable(table)
        genericItems = req.RequestItems[table]
    }
    if len(genericItems) == cap(genericItems) {
        newItems := make([]BatchWriteItem, len(genericItems) * 2)
        for i := range genericItems {
            newItems[i] = genericItems[i]
        }
        genericItems = newItems[:len(genericItems)]
    }
    genericItems = genericItems[:len(genericItems) + 1]
    return genericItems
}

/** Adds the delete request items into the request.
 * @param table the name of the table to modify
 * @param items a map of awsgo.AwsStringItem and awsgo.AwsNumberItem
 */
func (req *BatchWriteItemRequest) AddDeleteRequest(table string, items map[string]interface{}) {
    genericItems := req.checkAndGrow(table)

    writeReq := new(BatchWriteItemDeleteRequest)
    writeReq.Key = make(map[string]interface{})
    for key, val := range items {
        writeReq.Key[key] = val
    }

    genericItems[len(genericItems) - 1] = BatchWriteItem{DeleteRequest: writeReq}
    req.RequestItems[table] = genericItems
}

/** Adds the put request items into the request.
 * @param table the name of the table to modify
 * @param items a map of awsgo.AwsStringItem and awsgo.AwsNumberItem
 */
func (req *BatchWriteItemRequest) AddPutRequest(table string, items map[string]interface{}) {
    genericItems := req.checkAndGrow(table)

    writeReq := new(BatchWriteItemPutRequest)
    writeReq.Item = make(map[string]interface{})
    for key, val := range items {
        writeReq.Item[key] = val
    }

    genericItems[len(genericItems) - 1] = BatchWriteItem{PutRequest: writeReq}
    req.RequestItems[table] = genericItems
}

func NewBatchWriteItemRequest() * BatchWriteItemRequest {
    req := new(BatchWriteItemRequest)
    req.RequestItems = make(map[string][]BatchWriteItem)
    req.ReturnConsumedCapacity = ConsumedCapacity_NONE
    req.ReturnItemCollectionMetrics = ItemCollectionMetrics_NONE
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = BatchWriteItemTarget
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"
    return req
}

func (gir BatchWriteItemRequest) DeMarshalGetItemResponse(response []byte, headers map[string]string) (interface{}) {
    if err := CheckForErrorResponse(response); err != nil {
        return err
    }
    giResponse := new(BatchWriteItemResponse)
    err := json.Unmarshal([]byte(response), giResponse)
    if err != nil {
        return err
    }
    giResponse.UnprocessedItems = make(map[string][]BatchWriteItem)

    //                        table  | operations| key/item | name     | item type/value
    // RawUnprocessItems  map[string][]map[string]map[string]map[string]map[string]string `json:"UnprocessedItems"`
    for table, operations := range giResponse.RawUnprocessItems {
        giResponse.UnprocessedItems[table] = make([]BatchWriteItem, len(operations))
        for op := range operations {
            if opName, ok := operations[op]["DeleteRequest"]; ok {
                var del BatchWriteItemDeleteRequest
                del.Key = make(map[string]interface{})
                awsgo.FromRawMapToAwsItemMap(opName["Key"], del.Key)
                giResponse.UnprocessedItems[table][op] = BatchWriteItem{DeleteRequest: &del}
            } else if opName, ok := operations[op]["PutRequest"]; ok {
                var put BatchWriteItemPutRequest
                put.Item = make(map[string]interface{})
                awsgo.FromRawMapToAwsItemMap(opName["Item"], put.Item)
                giResponse.UnprocessedItems[table][op] = BatchWriteItem{PutRequest: &put}
            }
        }
    }
    return giResponse
}

func (gir * BatchWriteItemRequest) VerifyInput() (error) {
    gir.Host.Service = "dynamodb"
    if len(gir.RequestItems) == 0 {
        return errors.New("RequestItems cannot be empty")
    }
    if len(gir.Host.Region) == 0 {
        return errors.New("Host.Region cannot be empty")
    }
    return gir.RequestBuilder.VerifyInput()
}

func (gir BatchWriteItemRequest) Request() (*BatchWriteItemResponse, error) {    
    request, err := awsgo.BuildRequest(&gir, gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := awsgo.DoRequest(&gir, request)
    if resp == nil {
        return nil, err
    }
    return resp.(*BatchWriteItemResponse), err
}

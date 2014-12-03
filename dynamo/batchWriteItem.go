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
    "time"
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
    req.Host.Domain = "amazonaws.com"
    return req
}

func (gir BatchWriteItemRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := CheckForErrorResponse(response, statusCode); err != nil {
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
                awsgo.FromRawMapToEasyTypedMap(opName["Key"], del.Key)
                giResponse.UnprocessedItems[table][op] = BatchWriteItem{DeleteRequest: &del}
            } else if opName, ok := operations[op]["PutRequest"]; ok {
                var put BatchWriteItemPutRequest
                put.Item = make(map[string]interface{})
                awsgo.FromRawMapToEasyTypedMap(opName["Item"], put.Item)
                giResponse.UnprocessedItems[table][op] = BatchWriteItem{PutRequest: &put}
            }
        }
    }
    if giResponse.ItemCollectionMetrics != nil {
        giResponse.ItemCollectionMetrics.ItemCollectionKey = make(map[string]interface{})
        awsgo.FromRawMapToEasyTypedMap(giResponse.ItemCollectionMetrics.RawItemCollectionKey, giResponse.ItemCollectionMetrics.ItemCollectionKey)
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
    for _, batchRequest := range gir.RequestItems {
        for s := range batchRequest {
            if batchRequest[s].DeleteRequest != nil {
                for k, v := range batchRequest[s].DeleteRequest.Key {
                    batchRequest[s].DeleteRequest.Key[k] = awsgo.ConvertToAwsItem(v)
                }
            }
            if batchRequest[s].PutRequest != nil {
                for k, v := range batchRequest[s].PutRequest.Item {
                    batchRequest[s].PutRequest.Item[k] = awsgo.ConvertToAwsItem(v)
                }
            }
        }
    }
    return nil
}

func (gir BatchWriteItemRequest) Request() (*BatchWriteItemResponse, error) {
    request, err := awsgo.NewAwsRequest(&gir, gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*BatchWriteItemResponse), err
}

// makes the request, and tries to include the unprocessed request items in subsequent requests
func (gir BatchWriteItemRequest) RequestIncludingUnprocessed() (*BatchWriteItemResponse, error) {
    headerCopy := make(map[string]string)
    for k, v := range gir.Headers {
        headerCopy[k] = v
    }

    resp, err := gir.Request()
    if err != nil {
        return resp, err
    }
    lastSize := len(resp.UnprocessedItems)
    for backOff := 0; len(resp.UnprocessedItems) > 0; {
        if lastSize == len(resp.UnprocessedItems) {
            backOff ++
        }
        lastSize = len(resp.UnprocessedItems)
        if backOff > 5 {
            return resp, errors.New("Backoff exceeded. Giving up.")
        }
        time.Sleep(time.Duration(backOff * 100))

        retryRequest := NewBatchWriteItemRequest()
        retryRequest.RequestBuilder = gir.deepCopyRequestBuilder(headerCopy)
        retryRequest.RequestItems = resp.UnprocessedItems
        resp, err = retryRequest.Request()
        if err != nil {
            return resp, err
        }
    }
    return resp, nil
}

func (gir BatchWriteItemRequest) deepCopyRequestBuilder(headerCopy map[string]string) awsgo.RequestBuilder {
    var theCopy awsgo.RequestBuilder
    theCopy.Host.Service = gir.Host.Service
    theCopy.Host.Region = gir.Host.Region
    theCopy.Host.Domain = gir.Host.Domain
    theCopy.Host.Override = gir.Host.Override
    // ignore the cert for now
    // theCopy.Host.CustomCertificates

    theCopy.Key = gir.Key
    theCopy.Headers = make(map[string]string)
    for k, v := range headerCopy {
        theCopy.Headers[k] = v
    }
    theCopy.RequestMethod = gir.RequestMethod
    theCopy.CanonicalUri = gir.CanonicalUri
    return theCopy
}

// Makes multiple requests, if too many actions were added.
// Also automatically retries unprocessed items.
// returns on the first error.
// @param sleep - the amount of time to sleep between requests. 0 implies no added sleeping
func (gir BatchWriteItemRequest) RequestSplit(sleep time.Duration) ([]*BatchWriteItemResponse, error) {

    responses := make([]*BatchWriteItemResponse, 0, 10)

    headerCopy := make(map[string]string)
    for k, v := range gir.Headers {
        headerCopy[k] = v
    }

    var curSubRequest *BatchWriteItemRequest
    itemsInSet := 0
    for table, reqs := range gir.RequestItems {
        for i := range reqs {
            if curSubRequest == nil {
                curSubRequest = NewBatchWriteItemRequest()
                curSubRequest.RequestBuilder = gir.deepCopyRequestBuilder(headerCopy)
                curSubRequest.ReturnConsumedCapacity = gir.ReturnConsumedCapacity
                curSubRequest.ReturnItemCollectionMetrics = gir.ReturnItemCollectionMetrics
                itemsInSet = 0
            }
            if reqs[i].DeleteRequest != nil {
                curSubRequest.AddDeleteRequest(table, reqs[i].DeleteRequest.Key)
            } else if reqs[i].PutRequest != nil {
                curSubRequest.AddPutRequest(table, reqs[i].PutRequest.Item)
            }
            itemsInSet ++
            if itemsInSet >= 25 {
                resp, err := curSubRequest.RequestIncludingUnprocessed()
                if err != nil {
                    return responses, err
                }
                if sleep != 0 {
                    time.Sleep(sleep)
                }
                if len(responses) == cap(responses) {
                    newItems := make([]*BatchWriteItemResponse, len(responses) * 2)
                    for i := range responses {
                        newItems[i] = responses[i]
                    }
                    responses = newItems[:len(responses)]
                }
                responses = responses[:len(responses) + 1]
                responses[len(responses) - 1] = resp
                itemsInSet = 0
                curSubRequest = nil
            }
        }
    }
    if curSubRequest != nil {
        resp, err := curSubRequest.RequestIncludingUnprocessed()
        if err != nil {
            return responses, err
        }
        if len(responses) == cap(responses) {
            newItems := make([]*BatchWriteItemResponse, len(responses) * 2)
            for i := range responses {
                newItems[i] = responses[i]
            }
            responses = newItems[:len(responses)]
        }
        responses = responses[:len(responses) + 1]
        responses[len(responses) - 1] = resp
    }
    return responses, nil
}

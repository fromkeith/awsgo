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
    "github.com/fromkeith/awsgo"
)


type SetProvisionedThroughput struct {
    ReadCapacityUnits           float64
    WriteCapacityUnits          float64
}

type globalSecondaryIndexUpdate struct {
    IndexName                   string
    ProvisionedThroughput       SetProvisionedThroughput
}

type setGlobalSecondaryIndex struct {
    Update              globalSecondaryIndexUpdate
}

type UpdateTableRequest struct {
    awsgo.RequestBuilder

    GlobalSecondaryIndexUpdates         []setGlobalSecondaryIndex       `json:",omitempty"`
    ProvisionedThroughput               *SetProvisionedThroughput       `json:",omitempty"`
    TableName                           string
}


type updateTableDescription struct {
    AttributeDefinitions        []attributeDefinitions
    CreationDateTime            float64
    GlobalSecondaryIndexes      []secondaryIndex
    ItemCount                   float64
    KeySchema                   []keySchema
    LocalSecondaryIndexes       []secondaryIndex
    ProvisionedThroughput       provisionedThroughput
    TableName                   string
    TableSizeBytes              float64
    TableStatus                 string
}

type UpdateTableResponse struct {
    TableDescription            updateTableDescription
}




// Creates a new UpdateTableRequest, populating in some defaults
func NewUpdateTableRequest() *UpdateTableRequest {
    req := new(UpdateTableRequest)
    req.Host.Service = "dynamodb"
    req.Host.Region = ""
    req.Host.Domain = "amazonaws.com"
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = UpdateTableTarget
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"
    return req
}


func (req * UpdateTableRequest) VerifyInput() (error) {
    if len(req.Host.Service) == 0 {
        return Verification_Error_ServiceEmpty
    }
    if len(req.TableName) == 0 {
        return Verification_Error_TableNameEmpty
    }
    if len(req.Host.Region) == 0 {
        return Verification_Error_RegionEmpty
    }
    return nil
}


func (req UpdateTableRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := CheckForErrorResponse(response, statusCode); err != nil {
        return err
    }
    resp := new(UpdateTableResponse)
    err := json.Unmarshal(response, resp)
    if err != nil {
        newErr := &awsgo.UnmarhsallingError {
            ActualContent : string(response),
            MarshallError : err,
        }
        return newErr
    }
    return resp
}

func (gir UpdateTableRequest) Request() (*UpdateTableResponse, error) {
    request, err := awsgo.NewAwsRequest(&gir, gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*UpdateTableResponse), err
}

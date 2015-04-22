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

package s3

import (
    "github.com/fromkeith/awsgo"
    "errors"
    "fmt"
    "io"
    "strings"
)


type PutObjectRequest struct {
    awsgo.RequestBuilder

    Source io.ReadCloser
    Length int64
    ContentType string
    Permissions string
    Path string
    ServerSideEncryption        bool
}

type PutObjectResponse struct {
    Hash string
    RequestId string
    RequestId2 string
    VersionId string
    StatusCode      int
}
type PutObjectResponseFuture struct {
    response chan *PutObjectResponse
    errResponse chan error
}

func NewPutObjectRequest() (*PutObjectRequest) {
    req := new(PutObjectRequest)
    req.Headers = make(map[string]string)
    req.ServerSideEncryption = false
    req.RequestMethod = "PUT"
    req.Host.Domain = "amazonaws.com"
    return req
}

type BadStatusCodeError struct {
    StatusCode              int
    Content                 string
}
func (b BadStatusCodeError) Error() string {
    return fmt.Sprintf("Code: %d", b.StatusCode)
}

func (por PutObjectRequest) DeMarshalResponse(a []byte, headers map[string]string, statusCode int) (interface{}) {
    if headers == nil {
        return nil
    }
    if statusCode < 200 || statusCode >= 300 {
        return BadStatusCodeError{
            StatusCode: statusCode,
            Content: string(a),
        }
    }
    response := new(PutObjectResponse)
    if v, ok := headers["etag"]; ok {
        response.Hash = strings.Trim(v, "\"")
    }
    if v, ok := headers["x-amz-id-2"]; ok {
        response.RequestId2 = v
    }
    if v, ok := headers["x-amz-request-id"]; ok {
        response.RequestId = v
    }
    if v, ok := headers["x-amz-version-id"]; ok {
        response.VersionId = v
    }
    response.StatusCode = statusCode
    return response
}

func (por * PutObjectRequest) VerifyInput() (error) {
    por.Host.Service = "s3"
    if len(por.ContentType) == 0 {
        return errors.New("ContentType be empty")
    }
    if len(por.Permissions) == 0 {
        return errors.New("Permissions be empty")
    }
    if len(por.Path) == 0 {
        return errors.New("Path be empty")
    }
    por.Headers["Content-Type"] = por.ContentType
    por.Headers["x-amz-acl"] = por.Permissions
    por.Headers["Content-Length"] = fmt.Sprintf("%d", por.Length)
    por.Headers["Expect"] = "100-continue"
    if por.ServerSideEncryption {
        por.Headers["x-amz-server-side-encryption"] = "AES256"
    }
    por.CanonicalUri = fmt.Sprintf("/%s", por.Path)
    return nil
}


func (por PutObjectRequest) CoRequest() (*PutObjectResponseFuture, error) {
    request, err := awsgo.NewAwsRequest(&por, por.Source)
    if err != nil {
        return nil, err
    }
    future := new(PutObjectResponseFuture)
    future.errResponse = make(chan error)
    future.response = make(chan * PutObjectResponse)
    go por.CoDoAndDemarshall(request, future)
    return future, nil
}

func (por PutObjectRequest) CoDoAndDemarshall(request awsgo.AwsRequest, future * PutObjectResponseFuture) {
    request.RequestSigningType = awsgo.RequestSigningType_REST
    resp, err := request.DoAndDemarshall(&por)
    if err != nil {
        future.errResponse <- err
    } else {
        future.response <- resp.(*PutObjectResponse)
    }
    close(future.errResponse)
    close(future.response)
}

func (por PutObjectRequest) Request() (*PutObjectResponse, error) {
    request, err := awsgo.NewAwsRequest(&por, por.Source)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_REST
    resp, err := request.DoAndDemarshall(&por)
    if resp == nil {
        return nil, err
    }
    return resp.(*PutObjectResponse), err
}
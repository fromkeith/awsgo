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
    "fmt"
)


type GetObjectRequest struct {
    awsgo.RequestBuilder

    // headers
    Range               string
    IfModifiedSince     string
    IfUnmodifiedSince   string
    IfMatch             string
    IfNoneMatch         string

    // the file
    Path                string
}

type GetObjectResponse struct {
    // headers
    DeleteMarker        bool
    Expiration          string
    Encyption           string
    Restore             string
    VersionId           string
    WebsiteRedirectLocation string
    // Actual Data
    Data                []byte
    StatusCode          int
}

func NewGetObjectRequest() (*GetObjectRequest) {
    req := new(GetObjectRequest)
    req.Headers = make(map[string]string)
    req.RequestMethod = "GET"
    return req
}

func (por GetObjectRequest) DeMarshalResponse(a []byte, headers map[string]string, statusCode int) (interface{}) {
    if headers == nil {
        return nil
    }
    response := new(GetObjectResponse)
    if _, ok := headers["x-amz-delete-marker"]; ok {
        response.DeleteMarker = true // if false, it won't appear according to docs
    }
    if v, ok := headers["x-amz-expiration"]; ok {
        response.Expiration = v
    }
    if v, ok := headers["x-amz-server-side​-encryption"]; ok {
        response.Encyption= v
    }
    if v, ok := headers["x-amz-restore"]; ok {
        response.Restore = v
    }
    if v, ok := headers["x-amz-version-id"]; ok {
        response.VersionId = v
    }
    if v, ok := headers["x-amz-website-redirect-location"]; ok {
        response.WebsiteRedirectLocation = v
    }
    response.Data = a
    response.StatusCode = statusCode
    return response
}

func (gi * GetObjectRequest) VerifyInput() (error) {
    gi.Host.Service = "s3"
    
    if gi.Range != "" {
        gi.Headers["Range"] = gi.Range
    }
    if gi.IfModifiedSince != "" {
        gi.Headers["If-Modified-Since"] = gi.IfModifiedSince
    }
    if gi.IfUnmodifiedSince != "" {
        gi.Headers["If-Unmodified-Since"] = gi.IfUnmodifiedSince
    }
    if gi.IfMatch != "" {
        gi.Headers["If-Match"] = gi.IfMatch
    }
    if gi.IfNoneMatch != "" {
        gi.Headers["If-None-Match"] = gi.IfNoneMatch
    }
    
    gi.CanonicalUri = fmt.Sprintf("/%s", gi.Path)
    return gi.RequestBuilder.VerifyInput()
}



func (gor GetObjectRequest) Request() (*GetObjectResponse, error) {
    request, err := awsgo.BuildReaderRequest(&gor, nil)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_REST
    resp, err := awsgo.DoRequest(&gor, request)
    if resp == nil {
        return nil, err
    }
    return resp.(*GetObjectResponse), err
}
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

package cloudsearch


import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

var (
    Verification_Error_EndpointEmpty = errors.New("Endpoint can not be empty!")
    NoPermission_For_Endpoint = errors.New("403 Invalid permissions to hit endpoint.")
)


type DocumentInfo struct {
    // 'add' or 'delete'
    Type            string      `json:"type"`
    Id              string      `json:"id"`
    Fields          interface{} `json:"fields,omitempty"`
}


type BatchDocumentRequest struct {
    Items           []DocumentInfo
    // if you full endpoint is blah.us-west-2.cloudsearch.amazonaws.com
    // then set this to 'blah'
    Endpoint        string      `json:"-"`
    Region          string
}

type BatchDocumentMessage struct {
    Message         string  `json:"message"`
}

type BatchDocumentResponse struct {
    Status          string      `json:"status"`
    Adds            int         `json:"adds"`
    Deletes         int         `json:"deletes"`
    Errors          []BatchDocumentMessage `json:"errors"`
    Warnings        []BatchDocumentMessage `json:"warnings"`
    StatusCode      int
}


// Creates a new BatchDocumentRequest, populating in some defaults
func NewBatchDocumentRequest() *BatchDocumentRequest {
    req := new(BatchDocumentRequest)
    req.Region = ""
    return req
}


func (gir BatchDocumentRequest) DeMarshalResponse(response []byte, headers http.Header, statusCode int) (*BatchDocumentResponse, error) {
    if statusCode == 403 {
        return nil, NoPermission_For_Endpoint
    }
    resp := new(BatchDocumentResponse)
    err := json.Unmarshal(response, resp)
    if err != nil {
        return nil, err
    }
    resp.StatusCode = statusCode
    return resp, err
}


func (gir BatchDocumentRequest) Request() (*BatchDocumentResponse, error) {
    urlStr := fmt.Sprintf("https://%s.%s.cloudsearch.amazonaws.com/2013-01-01/documents/batch", gir.Endpoint, gir.Region)
    u, err := url.Parse(urlStr)
    if err != nil {
        return nil, err
    }
    itemJson, err := json.Marshal(gir.Items)
    if err != nil {
        return nil, err
    }
    hreq := http.Request{
        URL: u,
        ContentLength: int64(len(itemJson)),
        Body: ioutil.NopCloser(strings.NewReader(string(itemJson))),
        Header: http.Header{
            "Content-Type": []string{"application/json"},
        },
        Method: "POST",
        Close: true,
    }
    resp, err := http.DefaultClient.Do(&hreq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    buf := bytes.Buffer{}
    io.Copy(&buf, resp.Body)
    return gir.DeMarshalResponse([]byte(buf.String()), resp.Header, resp.StatusCode)

}

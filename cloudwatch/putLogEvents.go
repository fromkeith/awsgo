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

package cloudwatch


import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/fromkeith/awsgo"
    "time"
)


type LogEvents struct {
    Message             string          `json:"message"`
    // in milliseconds
    Timestamp           int64           `json:"timestamp"`
}


type PutLogEventsRequest struct {
    awsgo.RequestBuilder

    LogGroupName            string          `json:"logGroupName"`
    LogStreamName           string          `json:"logStreamName"`
    SequenceToken           string          `json:"sequenceToken,omitempty"`
    LogEvents               []LogEvents     `json:"logEvents"`
}

type PutLogEventsResponse struct {
    NextSequenceToken       string          `json:"nextSequenceToken"`
}


func NewPutLogEventsRequest() *PutLogEventsRequest {
    req := new(PutLogEventsRequest)
    req.Host.Service = "logs"
    req.Host.Region = "us-east-1"
    req.Host.Domain = "amazonaws.com"
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = "Logs_20140328.PutLogEvents"
    req.Headers["Content-Type"] = "application/x-amz-json-1.1"
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"
    req.LogEvents = make([]LogEvents, 0, 50)
    return req
}

func (req * PutLogEventsRequest) AddEvent(message string, at time.Time) {
    req.LogEvents = append(req.LogEvents, LogEvents{
        Message: message,
        Timestamp: at.UnixNano() / 1000000,
    })
}

func (req * PutLogEventsRequest) VerifyInput() (error) {
    if len(req.LogGroupName) == 0 || len(req.LogGroupName) > 512 {
        return errors.New("Bad LogGroupName")
    }
    if len(req.LogStreamName) == 0 || len(req.LogStreamName) > 512 {
        return errors.New("Bad LogStreamName")
    }
    // TODO: regex to check LogGroupName
    return nil
}

func (req PutLogEventsRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := CheckForErrorResponse(response, statusCode); err != nil {
        return err
    }
    if statusCode != 200 {
        return errors.New(fmt.Sprintf("Bad Status code: %d", statusCode))
    }
    resp := new(PutLogEventsResponse)
    if err := json.Unmarshal(response, resp); err != nil {
        return err
    }
    return resp
}

// events and entire content cannot be above 32,768 bytes. If you are unsure of your
// content size, use RequestSplit. It will use multiple requests to send your logs.
func (req PutLogEventsRequest) Request() (*PutLogEventsResponse, error) {
    request, err := awsgo.NewAwsRequest(&req, req)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&req)
    if resp == nil {
        return nil, err
    }
    return resp.(*PutLogEventsResponse), err
}


func (req PutLogEventsRequest) RequestSplit() (*PutLogEventsResponse, error) {
    var err error
    var resp *PutLogEventsResponse
    start := 0
    calcSize := 0
    for i := range req.LogEvents{
        newSize := calcSize + len(req.LogEvents[i].Message) + 60 // 60 is an over approx json surrounding the message + timestamp
        if newSize >= 30000 || i - start == 1000 {
            resp, err = split(req, start, i)
            if err != nil {
                return nil, err
            }
            time.Sleep(100 * time.Millisecond)
            req.SequenceToken = resp.NextSequenceToken
            calcSize = newSize - calcSize
            start = i
        }
        calcSize = newSize
    }
    return split(req, start, len(req.LogEvents))
}

func split(req PutLogEventsRequest, start, end int) (*PutLogEventsResponse, error) {
    newReq := NewPutLogEventsRequest()
    newReq.LogEvents = req.LogEvents[start:end]
    newReq.LogGroupName = req.LogGroupName
    newReq.LogStreamName = req.LogStreamName
    newReq.SequenceToken = req.SequenceToken
    newReq.Key = req.Key
    return newReq.Request()
}
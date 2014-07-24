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



type GetLogEventsRequest struct {
    awsgo.RequestBuilder

    LogGroupName            string          `json:"logGroupName"`
    LogStreamName           string          `json:"logStreamName"`
    EndTime                 *int64          `json:"endTime,omitempty"`
    Limit                   *int            `json:"limit,omitempty"`
    NextToken               string          `json:"nextToken,omitempty"`
    StartFromHead           bool            `json:"startFromHead"`
    StartTime               *int64          `json:"startTime,omitempty"`

}

type LogEventsResponse struct {
    LogEvents
    IngestionTime           int64           `json:"ingestionTime"`
}

type GetLogEventsResponse struct {
    NextBackwardToken       string          `json:"nextBackwardToken"`
    NextForwardToken        string          `json:"nextForwardToken"`
    Events                  []LogEventsResponse `json:"events"`
}


func NewGetLogEventsRequest() *GetLogEventsRequest {
    req := new(GetLogEventsRequest)
    req.Host.Service = "logs"
    req.Host.Region = "us-east-1"
    req.Host.Domain = "amazonaws.com"
    req.Headers = make(map[string]string)
    req.Headers["X-Amz-Target"] = "Logs_20140328.GetLogEvents"
    req.Headers["Content-Type"] = "application/x-amz-json-1.1"
    req.RequestMethod = "POST"
    req.CanonicalUri = "/"
    req.StartFromHead = false
    return req
}

func (req * GetLogEventsRequest) SetTimeRange(start time.Time, end time.Time) {
    startInt := start.UnixNano() / 1000000
    req.StartTime = &startInt
    endInt := end.UnixNano() / 1000000
    req.EndTime = &endInt
}
func (req * GetLogEventsRequest) SetLimit(limit int) {
    req.Limit = &limit
}

func (req * GetLogEventsRequest) VerifyInput() (error) {
    if len(req.LogGroupName) == 0 || len(req.LogGroupName) > 512 {
        return errors.New("Bad Log Group Name")
    }
    if len(req.LogStreamName) == 0 || len(req.LogStreamName) > 512 {
        return errors.New("Bad LogStreamName")
    }
    // TODO: regex to check LogGroupName
    return nil
}

func (req GetLogEventsRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := CheckForErrorResponse(response, statusCode); err != nil {
        return err
    }
    if statusCode != 200 {
        return errors.New(fmt.Sprintf("Bad Status code: %d", statusCode))
    }
    resp := new(GetLogEventsResponse)
    if err := json.Unmarshal(response, resp); err != nil {
        return err
    }
    return resp
}


func (req GetLogEventsRequest) Request() (*GetLogEventsResponse, error) {
    request, err := awsgo.NewAwsRequest(&req, req)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&req)
    if resp == nil {
        return nil, err
    }
    return resp.(*GetLogEventsResponse), err
}

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
    "github.com/fromkeith/awsgo"
    "errors"
    "fmt"
    "net/url"
    "encoding/xml"
    "time"
    "strconv"
)

const (
    STATISTIC_Average = "Average"
    STATISTIC_Sum = "Sum"
    STATISTIC_SampleCount = "SampleCount"
    STATISTIC_Maximum = "Maximum"
    STATISTIC_Minimum = "Minimum"
)





type GetMetricStatisticsRequest struct {
    awsgo.RequestBuilder

    Dimensions      []MetricDimensions
    EndTime         time.Time
    MetricName      string
    Namespace       string
    Period          int
    StartTime       time.Time
    Statistics      []string
    Unit            string
}

type GetMetricResultDatapoint struct {
    Timestamp                   time.Time
    Unit                        string
    Average                     float64
    Sum                         float64
    SampleCount                 float64
    Maximum                     float64
    Minimum                     float64
}

type GetMetricStatisticsResult struct {
    Label                       string
    Datapoints                  []GetMetricResultDatapoint      `xml:"Datapoints>member"`
}


type GetMetricStatisticsResponse struct {
    ResponseMetadata                awsgo.ResponseMetaData
    GetMetricStatisticsResult       GetMetricStatisticsResult
}


func NewGetMetricStatisticsRequest() *GetMetricStatisticsRequest {
    req := new(GetMetricStatisticsRequest)
    req.Host.Service = "monitoring"
    req.Host.Region = ""
    req.Host.Domain = "amazonaws.com"
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    req.Headers = make(map[string]string)
    req.RequestMethod = "GET"
    req.CanonicalUri = "/"

    req.Period = 60
    return req
}


func (gir * GetMetricStatisticsRequest) VerifyInput() (error) {
    gir.Host.Service = "monitoring"
    if len(gir.Host.Region) == 0 {
        gir.Host.Region = "us-east-1"
    }

    if gir.EndTime == gir.StartTime {
        return errors.New("Invalid time range")
    }
    if gir.MetricName == "" {
        return errors.New("MetricName is empty")
    }
    if gir.Namespace == "" {
        return errors.New("Namespace is empty")
    }
    if gir.Period < 60 {
        return errors.New("Period is invalid")
    }
    if len(gir.Statistics) == 0 {
        return errors.New("Statistics cannot be empty")
    }
    if gir.StartTime.After(gir.EndTime) {
        return errors.New("Start time must be before endtime")
    }

    vals := url.Values{}
    vals.Set("Action", "GetMetricStatistics")
    vals.Set("Version", "2010-08-01")

    for i := range gir.Dimensions {
        vals.Set(fmt.Sprintf("Dimensions.member.%d.Name", i + 1), gir.Dimensions[i].Name)
        vals.Set(fmt.Sprintf("Dimensions.member.%d.Value", i + 1), gir.Dimensions[i].Value)
    }
    vals.Set("EndTime", awsgo.IsoDate(gir.EndTime))
    vals.Set("MetricName", gir.MetricName)
    vals.Set("Namespace", gir.Namespace)
    vals.Set("Period", strconv.FormatInt(int64(gir.Period), 10))
    vals.Set("StartTime", awsgo.IsoDate(gir.StartTime))
    for i := range gir.Statistics {
        vals.Set(fmt.Sprintf("Statistics.member.%d", i + 1), gir.Statistics[i])
    }
    if gir.Unit != "" {
        vals.Set("Unit", gir.Unit)
    }

    gir.CanonicalUri = "/?" + vals.Encode()
    return nil
}

func (gir GetMetricStatisticsRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    if err := awsgo.CheckForErrorXml(response); err != nil {
        return err
    }
    giResponse := new(GetMetricStatisticsResponse)
    if err := xml.Unmarshal(response, giResponse); err != nil {
        return err
    }
    return giResponse
}

func (gir GetMetricStatisticsRequest) Request() (*GetMetricStatisticsResponse, error) {
    request, err := awsgo.BuildEmptyContentRequest(&gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*GetMetricStatisticsResponse), err
}

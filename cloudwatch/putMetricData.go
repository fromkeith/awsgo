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
)

const (
    UNIT_SECONDS = "Seconds"
    UNIT_MICROSECONDS = "Microseconds"
    UNIT_MILLISECONDS = "Milliseconds"
    UNIT_BYTES = "Bytes"
    UNIT_KILOBYTES = "Kilobytes"
    UNIT_MEGABYTES = "Megabytes"
    UNIT_GIGABYTES = "Gigabytes"
    UNIT_TERABYTES = "Terabytes"
    UNIT_BITS = "Bits"
    UNIT_KILOBITS = "Kilobits"
    UNIT_MEGABITS = "Megabits"
    UNIT_GIGABITS = "Gigabits"
    UNIT_TERABITS = "Terabits"
    UNIT_PERCENT = "Percent"
    UNIT_COUNT = "Count"
    UNIT_BYTES_PER_SECOND = "Bytes/Second"
    UNIT_KILOBYTES_PER_SECOND = "Kilobytes/Second"
    UNIT_MEGABYTES_PER_SECOND = "Megabytes/Second"
    UNIT_GIGABYTES_PER_SECOND = "Gigabytes/Second"
    UNIT_TERABYTES_PER_SECOND = "Terabytes/Second"
    UNIT_BITS_PER_SECOND = "Bits/Second"
    UNIT_KILOBITS_PER_SECOND = "Kilobits/Second"
    UNIT_MEGABITS_PER_SECOND = "Megabits/Second"
    UNIT_GIGABITS_PER_SECOND = "Gigabits/Second"
    UNIT_TERABITS_PER_SECOND = "Terabits/Second"
    UNIT_COUNT_PER_SECOND = "Count/Second"
    UNIT_NONE = "None"
)


type StatisticSet struct {
    Maximum         float64
    Minimum         float64
    SampleCount     float64
    Sum             float64
}

type MetricDimensions struct {
    Name            string
    Value           string
}

type MetricDatum struct {
    Dimensions      []MetricDimensions
    MetricName      string
    StatisticValues * StatisticSet
    Timestamp       * time.Time
    Unit            string
    Value           * float64

}

type PutMetricRequest struct {
    awsgo.RequestBuilder

    MetricData      []MetricDatum
    Namespace       string
}


type PutMetricResponse struct {
    ResponseMetadata awsgo.ResponseMetaData
}


func NewPutMetricRequest() *PutMetricRequest {
    req := new(PutMetricRequest)
    req.Host.Service = "monitoring"
    req.Host.Region = ""
    req.Host.Domain = ""
    req.Key.AccessKeyId = ""
    req.Key.SecretAccessKey = ""
    req.Headers = make(map[string]string)
    req.RequestMethod = "GET"
    req.CanonicalUri = "/"
    return req
}


func addMetricDatumToUri(vals url.Values, datum MetricDatum, index int) (error) {
    if len(datum.MetricName) == 0 {
        return errors.New("MetricDatum items must specify a name!")
    }

    vals.Set(fmt.Sprintf("MetricData.member.%d.MetricName", index), datum.MetricName)

    if len(datum.Dimensions) > 0 {
        for i := range(datum.Dimensions) {
            vals.Set(
                fmt.Sprintf("MetricData.member.%d.Dimensions.member.%d.Name", index, i + 1),
                datum.Dimensions[i].Name,
            )
            vals.Set(
                fmt.Sprintf("MetricData.member.%d.Dimensions.member.%d.Value", index, i + 1),
                datum.Dimensions[i].Value,
            )
        }
    }
    if datum.StatisticValues != nil {
        vals.Set(
            fmt.Sprintf("MetricData.member.%d.StatisticValues.Maximum", index),
            fmt.Sprintf("%f", datum.StatisticValues.Maximum),
        )
        vals.Set(
            fmt.Sprintf("MetricData.member.%d.StatisticValues.Minimum", index),
            fmt.Sprintf("%f", datum.StatisticValues.Minimum),
        )
        vals.Set(
            fmt.Sprintf("MetricData.member.%d.StatisticValues.SampleCount", index),
            fmt.Sprintf("%f", datum.StatisticValues.SampleCount),
        )
        vals.Set(
            fmt.Sprintf("MetricData.member.%d.StatisticValues.Sum", index),
            fmt.Sprintf("%f", datum.StatisticValues.Sum),
        )
    }
    if datum.Timestamp != nil {
        vals.Set(
            fmt.Sprintf("MetricData.member.%d.Timestamp", index),
            awsgo.IsoDate(*datum.Timestamp),
        )
    }
    if len(datum.Unit) > 0 {
        vals.Set(
            fmt.Sprintf("MetricData.member.%d.Unit", index),
            datum.Unit,
        )
    }
    if datum.Value != nil {
        vals.Set(
            fmt.Sprintf("MetricData.member.%d.Value", index),
            fmt.Sprintf("%f", *datum.Value),
        )
    }
    return nil
}


func (gir * PutMetricRequest) VerifyInput() (error) {
    gir.Host.Service = "monitoring"
    if len(gir.Host.Region) == 0 {
        gir.Host.Region = "us-east-1"
    }

    vals := url.Values{}
    vals.Set("Action", "PutMetricData")
    vals.Set("Version", "2010-08-01")

    for i := range(gir.MetricData) {
        var err error
        err = addMetricDatumToUri(vals, gir.MetricData[i], i + 1)
        if err != nil {
            return err
        }
    }
    if len(gir.Namespace) == 0 {
        return errors.New("Namespace cannot be empty!")
    }

    vals.Set("Namespace", gir.Namespace)
    gir.CanonicalUri = "/?" + vals.Encode()
    //fmt.Println("PutMetric.Url: ", gir.CanonicalUri)

    return nil
}

func (gir PutMetricRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    giResponse := new(PutMetricResponse)
    //fmt.Println("PutMetric: ", string(response))
    //fmt.Println("PutMetric.StatusCode: ", statusCode)
    if statusCode != 200 {
        return errors.New(fmt.Sprintf("Bad status code: %d", statusCode))
    }
    xml.Unmarshal(response, giResponse)
    //json.Unmarshal([]byte(response), giResponse)
    return giResponse
}

func (gir PutMetricRequest) Request() (*PutMetricResponse, error) {    
    request, err := awsgo.BuildEmptyContentRequest(&gir)
    if err != nil {
        return nil, err
    }
    request.RequestSigningType = awsgo.RequestSigningType_AWS4
    resp, err := request.DoAndDemarshall(&gir)
    if resp == nil {
        return nil, err
    }
    return resp.(*PutMetricResponse), err
}

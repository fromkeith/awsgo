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


func addMetricDatumToUri(uri string, datum MetricDatum, index int) (string, error) {
    if len(datum.MetricName) == 0 {
        return "", errors.New("MetricDatum items must specify a name!")
    }
    uri = fmt.Sprintf("%s&MetricData.member.%d.MetricName=%s",
        uri, index, url.QueryEscape(datum.MetricName))
    if len(datum.Dimensions) > 0 {
        for i := range(datum.Dimensions) {
            uri = fmt.Sprintf("%s&MetricData.member.%d.Dimensions.member.%d.Name=%s" +
                "&MetricData.member.%d.Dimensions.member.%d.Value=%s",
                uri, index, i + 1, url.QueryEscape(datum.Dimensions[i].Name),
                index, i + 1, url.QueryEscape(datum.Dimensions[i].Value))
        }
    }
    if datum.StatisticValues != nil {
        uri = fmt.Sprintf("%s&MetricData.member.%d.StatisticValues.Maximum=%f" +
            "&MetricData.member.%d.StatisticValues.Minimum=%f" + 
            "&MetricData.member.%d.StatisticValues.SampleCount=%f" +
            "&MetricData.member.%d.StatisticValues.Sum=%f",
            uri,
            index, datum.StatisticValues.Maximum,
            index, datum.StatisticValues.Minimum,
            index, datum.StatisticValues.SampleCount,
            index, datum.StatisticValues.Sum)
    }
    if datum.Timestamp != nil {
        uri = fmt.Sprintf("%s&MetricData.member.%d.Timestamp=%s",
            uri, index, url.QueryEscape(awsgo.IsoDate(*datum.Timestamp)))
    }
    if len(datum.Unit) > 0 {
        uri = fmt.Sprintf("%s&MetricData.member.%d.Unit=%s",
            uri, index, url.QueryEscape(datum.Unit))
    }
    if datum.Value != nil {
        uri = fmt.Sprintf("%s&MetricData.member.%d.Value=%f",
            uri, index, *datum.Value)
    }
    return uri, nil
}


func (gir * PutMetricRequest) VerifyInput() (error) {
    gir.Host.Service = "monitoring"
    if len(gir.Host.Region) == 0 {
        gir.Host.Region = "us-east-1"
    }

    gir.CanonicalUri = "/?Action=PutMetricData&Version=2010-08-01"
    for i := range(gir.MetricData) {
        var err error
        gir.CanonicalUri, err = addMetricDatumToUri(gir.CanonicalUri, gir.MetricData[i], i + 1)
        if err != nil {
            return err
        }
    }
    if len(gir.Namespace) == 0 {
        return errors.New("Namespace cannot be empty!")
    }
    gir.CanonicalUri = fmt.Sprintf("%s&Namespace=%s", gir.CanonicalUri, url.QueryEscape(gir.Namespace))

    return gir.RequestBuilder.VerifyInput()
}

func (gir PutMetricRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{}) {
    giResponse := new(PutMetricResponse)
    //fmt.Println(string(response))
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
    resp, err := awsgo.DoRequest(&gir, request)
    if resp == nil {
        return nil, err
    }
    return resp.(*PutMetricResponse), err
}

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
package examples

import (
    "fmt"
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/dynamo"
    "github.com/fromkeith/awsgo/s3"
    "github.com/fromkeith/awsgo/sqs"
    "github.com/fromkeith/awsgo/ses"
    "github.com/fromkeith/awsgo/cloudwatch"
    "io/ioutil"
    "bytes"
    "crypto/md5"
    "time"
)

func TestGetItem() {
    itemRequest := dynamo.NewGetItemRequest()
    itemRequest.Search["Game"] = aws.NewStringItem("e5dd6f4d-5c80-4069-817e-646372bf5f74")
    itemRequest.TableName = "test.table"
    itemRequest.AttributesToGet = []string{"GameName"}

    itemRequest.Host.Region = "us-west-2"
    itemRequest.Host.Domain = "amazonaws.com"
    itemRequest.Key.Key, itemRequest.Key.SecretKey, _ = aws.GetSecurityKeys()

    //resp, _ := itemRequest.Request()
    fu, futureError := itemRequest.CoRequest()
    if futureError != nil {
        fmt.Println(futureError)
        return
    }
    resp, getErr := fu.Get()
    if getErr != nil {
        fmt.Println(getErr)
        return
    }
    itemCast := resp.Item["GameName"]
    switch itemCast := itemCast.(type) {
    case aws.AwsStringItem:
        fmt.Println(itemCast.Value)
    case string:
        fmt.Println("isstring")
    default:
        fmt.Println("Unknown type %T", itemCast)
    }
}

func TestUpdateItem() {
    itemRequest := dynamo.NewUpdateItemRequest()
    itemRequest.UpdateKey["Game"] = aws.NewStringItem("e5dd6f4d-5c80-4069-817e-646372bf5f74")
    itemRequest.Update["GameName"] = dynamo.AttributeUpdates{"PUT", aws.NewStringItem("sadfsdfew ewr a dsf")}
    itemRequest.Expected["Holinn"] = dynamo.ExpectedItem{true, aws.NewNumberItem(1)}
    itemRequest.TableName = "testtable"

    itemRequest.Host.Region = "us-west-2"
    itemRequest.Host.Domain = "amazonaws.com"
    itemRequest.Key.Key, itemRequest.Key.SecretKey, _ = aws.GetSecurityKeys()

    resp, err := itemRequest.Request()
    if err != nil {
        fmt.Printf("WasError %T %v \n", err, err)
        return
    }
    fmt.Println(resp)
}

func TestPutItem() {
    itemRequest := dynamo.NewPutItemRequest()
    itemRequest.Item["Game"] = aws.NewStringItem("helloThere!")
    itemRequest.TableName = "testtable"

    itemRequest.Host.Region = "us-west-2"
    itemRequest.Host.Domain = "amazonaws.com"
    itemRequest.Key.Key, itemRequest.Key.SecretKey, _ = aws.GetSecurityKeys()

    resp, err := itemRequest.Request()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(resp)
}

func TestBatchGetItem() {
    itemRequest := dynamo.NewBatchGetItemRequest()
    tableReq := dynamo.NewBatchGetIteamRequestTable()
    tableReq.Search = make([]map[string]interface{}, 1)
    tableReq.Search[0] = make(map[string]interface{})
    tableReq.Search[0]["Game"] = aws.NewStringItem("e5dd6f4d-5c80-4069-817e-646372bf5f74")
    tableReq.AttributesToGet = []string{"GameName"}
    itemRequest.RequestItems["testtable"] = tableReq

    itemRequest.Host.Region = "us-west-2"
    itemRequest.Host.Domain = "amazonaws.com"
    itemRequest.Key.Key, itemRequest.Key.SecretKey, _ = aws.GetSecurityKeys()

    resp, err := itemRequest.Request()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(resp)
}

func TestPutS3File() {
    putRequest := s3.NewPutObjectRequest()
    putRequest.ContentType = "text/plain"
    putRequest.Permissions = "private"
    putRequest.Path = "someexamplebucket/haha/test.jpg"
    fakePayload := "1234567890"
    putRequest.Length = int64(len(fakePayload))
    putRequest.Source = ioutil.NopCloser(bytes.NewBuffer([]byte(fakePayload)))

    putRequest.Host.Domain = "amazonaws.com"
    putRequest.Key.Key, putRequest.Key.SecretKey, _ = aws.GetSecurityKeys()

    resp, err := putRequest.Request()
    if err != nil {
        fmt.Println(err)
        return
    }
    md5Hasher := md5.New()
    md5Hasher.Write([]byte(fakePayload))
    ourHash := fmt.Sprintf("%x", md5Hasher.Sum(nil))
    if ourHash != resp.Hash {
        fmt.Println("Our hash does not match returned hash!", ourHash, resp.Hash)
    } else {
        fmt.Println("File uploaded!")
    }
}


func TestSqsSendMessage() {
    sendRequest := sqs.NewSendMessageRequest()
    sendRequest.MessageBody = "hello"
    sendRequest.TaskQueue = "/<queueNumber>/TestQueue"

    sendRequest.Host.Region = "us-west-2"
    sendRequest.Host.Domain = "amazonaws.com"
    sendRequest.Key.Key, sendRequest.Key.SecretKey, _ = aws.GetSecurityKeys()

    resp, err := sendRequest.Request()
    fmt.Println(resp, err)

}

func TestSqsReceiveMessage() {
    sendRequest := sqs.NewReceiveMessageRequest()
    sendRequest.TaskQueue = "/<queueNumber>/TestQueue"
    sendRequest.WaitTimeSeconds = 10
    sendRequest.MaxNumberOfMessages = 1

    sendRequest.Host.Region = "us-west-2"
    sendRequest.Host.Domain = "amazonaws.com"
    sendRequest.Key.Key, sendRequest.Key.SecretKey, _ = aws.GetSecurityKeys()

    resp, err := sendRequest.Request()
    fmt.Println(resp, err)

}

func TestSesSendEmail() {
    sendRequest := ses.NewSendEmailRequest()
    sendRequest.Message.Body.Html.Data = "<h1>Hello</h1> There"
    sendRequest.Message.Subject.Data = "Wassup"
    sendRequest.Destination.ToAddresses = []string{"example@example.com"}
    sendRequest.Source = "example@example.com"

    sendRequest.Host.Domain = "amazonaws.com"
    sendRequest.Key.Key, sendRequest.Key.SecretKey, _ = aws.GetSecurityKeys()

    resp, err := sendRequest.Request()
    fmt.Println(resp, err)
}

func TestPutMetric() {
    putMetricRequest := cloudwatch.NewPutMetricRequest()
    putMetricRequest.Namespace = "SimpleTest/Test"
    putMetricRequest.MetricData = make([]cloudwatch.MetricDatum, 3)
    putMetricRequest.MetricData[0].MetricName = "MyTestMetric"
    putMetricRequest.MetricData[0].Unit = cloudwatch.UNIT_COUNT
    putMetricRequest.MetricData[0].Value = new(float64)
    *(putMetricRequest.MetricData[0].Value) = 55.0
    putMetricRequest.MetricData[1].MetricName = "MyOtherTestMetric"
    putMetricRequest.MetricData[1].Dimensions = make([]cloudwatch.MetricDimensions, 1)
    putMetricRequest.MetricData[1].Dimensions[0].Name = "DimOrig"
    putMetricRequest.MetricData[1].Dimensions[0].Value = "Yes"
    putMetricRequest.MetricData[1].Value = new(float64)
    *(putMetricRequest.MetricData[1].Value) = -63.23
    putMetricRequest.MetricData[1].Timestamp = new(time.Time)
    *(putMetricRequest.MetricData[1].Timestamp) = time.Now()
    putMetricRequest.MetricData[2].MetricName = "MyStat"
    putMetricRequest.MetricData[2].StatisticValues = new(cloudwatch.StatisticSet)
    putMetricRequest.MetricData[2].StatisticValues.Maximum = 10.0
    putMetricRequest.MetricData[2].StatisticValues.Minimum = 1.0
    putMetricRequest.MetricData[2].StatisticValues.SampleCount = 10.0
    putMetricRequest.MetricData[2].StatisticValues.Sum = 33

    putMetricRequest.Host.Region = "us-west-2"
    putMetricRequest.Host.Domain = "amazonaws.com"
    putMetricRequest.Key.Key, putMetricRequest.Key.SecretKey, _ = aws.GetSecurityKeys()

    resp, err := putMetricRequest.Request()
    fmt.Println(resp, err)
}
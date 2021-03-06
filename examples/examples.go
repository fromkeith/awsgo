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
package main

import (
    "bytes"
    "crypto/md5"
    "fmt"
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/cloudwatch"
    "github.com/fromkeith/awsgo/dynamo"
    "github.com/fromkeith/awsgo/ec2"
    "github.com/fromkeith/awsgo/s3"
    "github.com/fromkeith/awsgo/ses"
    "github.com/fromkeith/awsgo/sns"
    "github.com/fromkeith/awsgo/sqs"
    "github.com/fromkeith/awsgo/swf"
    "io/ioutil"
    "log"
    "math/rand"
    "time"
)

const (
    TEST_TABLE_NAME = "fromkeith.conquer.testgame"// "test.table"
    TEST_SQS_QUEUE = "1234"
    TEST_ITEM_NAME = "Game"
    TEST_S3_BUCKET = "testbucket"
)

func TestGetItem() {
    itemRequest := dynamo.NewGetItemRequest()
    itemRequest.Search[TEST_ITEM_NAME] = "e5dd6f4d-5c80-4069-817e-646372bf5f74"
    itemRequest.TableName = TEST_TABLE_NAME
    itemRequest.AttributesToGet = []string{"Num", "NumArray", "String", "StringArray"}

    itemRequest.Host.Region = "us-west-2"
    itemRequest.Host.Domain = "amazonaws.com"
    itemRequest.Key, _ = awsgo.GetSecurityKeys()

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
    s := resp.Item["String"]
    switch itemCast := s.(type) {
    case string:
        fmt.Println(itemCast)
    default:
        panic(fmt.Sprintf("Unexpected type: %T", itemCast))
    }
    ss := resp.Item["StringArray"]
    switch itemCast := ss.(type) {
    case []string:
        if len(itemCast) != 3 {
            panic(fmt.Sprintf("Values should have 3 items, got %d", len(itemCast)))
        }
        fmt.Println(itemCast)
    default:
        panic(fmt.Sprintf("Unexpected type: %T", itemCast))
    }
    n := resp.Item["Num"]
    switch itemCast := n.(type) {
    case float64:
        fmt.Println(itemCast)
    default:
        panic(fmt.Sprintf("Unexpected type: %T", itemCast))
    }
    nn := resp.Item["NumArray"]
    switch itemCast := nn.(type) {
    case []float64:
        if len(itemCast) != 3 {
            panic(fmt.Sprintf("Values should have 3 items, got %d", len(itemCast)))
        }
        fmt.Println(itemCast)
    default:
        panic(fmt.Sprintf("Unexpected type: %T", itemCast))
    }
}

func TestUpdateItem() {
    itemRequest := dynamo.NewUpdateItemRequest()
    itemRequest.UpdateKey[TEST_ITEM_NAME] = "e5dd6f4d-5c80-4069-817e-646372bf5f74"
    rand.Seed(time.Now().Unix())
    newName := fmt.Sprintf("%d", rand.Int())
    itemRequest.Update["GameName"] = dynamo.AttributeUpdates{"PUT", newName}
    itemRequest.Expected["Holinn"] = dynamo.ExpectedItem{true, 1}
    itemRequest.TableName = TEST_TABLE_NAME
    itemRequest.ReturnValues = dynamo.ReturnValues_UPDATED_NEW

    itemRequest.Host.Region = "us-west-2"
    itemRequest.Host.Domain = "amazonaws.com"
    itemRequest.Key, _ = awsgo.GetSecurityKeys()

    resp, err := itemRequest.Request()
    if err != nil {
        fmt.Printf("WasError %T %v \n", err, err)
        return
    }

    switch itemCast := resp.BeforeAttributes["GameName"].(type) {
    case string:
        if itemCast != newName {
            panic("Updated value was not updated!")
        }
    default:
        panic(fmt.Sprintf("Unexpected type: %T", itemCast))
    }
    fmt.Println(resp)
}

func TestPutItem() {
    itemRequest := dynamo.NewPutItemRequest()
    itemRequest.Item[TEST_ITEM_NAME] = "helloThere!s!!"
    itemRequest.Item["AnotherVal"] = "asdf"
    itemRequest.TableName = TEST_TABLE_NAME
    itemRequest.ReturnValues = dynamo.ReturnValues_ALL_OLD
    itemRequest.ReturnItemCollectionMetrics = dynamo.ItemCollectionMetrics_SIZE

    itemRequest.Host.Region = "us-west-2"
    itemRequest.Host.Domain = "amazonaws.com"
    itemRequest.Key, _ = awsgo.GetSecurityKeys()

    resp, err := itemRequest.Request()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(resp)
}

func TestDeleteItem() {
    deleteRequest := dynamo.NewDeleteItemRequest()
    deleteRequest.DeleteKey[TEST_ITEM_NAME] = "helloThere!"
    deleteRequest.TableName = TEST_TABLE_NAME
    deleteRequest.ReturnValues = dynamo.ReturnValues_ALL_OLD

    deleteRequest.Host.Region = "us-west-2"
    deleteRequest.Host.Domain = "amazonaws.com"
    deleteRequest.Key, _ = awsgo.GetSecurityKeys()

    resp, err := deleteRequest.Request()
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
    tableReq.Search[0][TEST_ITEM_NAME] = "e5dd6f4d-5c80-4069-817e-646372bf5f74"
    tableReq.AttributesToGet = []string{"GameName"}
    itemRequest.RequestItems[TEST_TABLE_NAME] = tableReq

    itemRequest.Host.Region = "us-west-2"
    itemRequest.Host.Domain = "amazonaws.com"
    itemRequest.Key, _ = awsgo.GetSecurityKeys()

    resp, err := itemRequest.Request()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(resp)
}

func testBatchWriteItem_Put(keys []string) {
    itemRequest := dynamo.NewBatchWriteItemRequest()
    for i := range keys {
        itemRequest.AddPutRequest(TEST_TABLE_NAME,
            map[string]interface{}{
                TEST_ITEM_NAME : keys[i],
                "GameName" : "gg",
            })
    }
    itemRequest.Host.Region = "us-west-2"
    itemRequest.Host.Domain = "amazonaws.com"
    itemRequest.Key, _ = awsgo.GetSecurityKeys()
    resp, err := itemRequest.Request()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(resp)
}
func testBatchWriteItem_Delete(keys []string) {
    itemRequest := dynamo.NewBatchWriteItemRequest()

    for i := range keys {
        itemRequest.AddDeleteRequest(TEST_TABLE_NAME,
            map[string]interface{}{
                TEST_ITEM_NAME : keys[i],
            })
    }
    itemRequest.Host.Region = "us-west-2"
    itemRequest.Host.Domain = "amazonaws.com"
    itemRequest.Key, _ = awsgo.GetSecurityKeys()
    resp, err := itemRequest.Request()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(resp)
    if len(resp.UnprocessedItems) > 0 {
        fmt.Println("WAiting to retry unprocessed items: ", len(resp.UnprocessedItems))
        time.Sleep(5000)
        itemRequest2 := dynamo.NewBatchWriteItemRequest()
        itemRequest2.RequestItems = resp.UnprocessedItems
        itemRequest2.Host.Region = "us-west-2"
        itemRequest2.Host.Domain = "amazonaws.com"
        itemRequest2.Key, _ = awsgo.GetSecurityKeys()
        resp, err := itemRequest2.Request()
        if err != nil {
        fmt.Println(err)
            return
        }
        fmt.Println(resp)
    }
}

func TestBatchWriteItem() {
    keys := make([]string, 25)
    for i := range keys {
        keys[i] = fmt.Sprintf("test%d", i)
    }

    testBatchWriteItem_Put(keys)
    testBatchWriteItem_Delete(keys)
}

func TestQuery() {
    req := dynamo.NewQueryRequest()
    req.AddKeyCondition(TEST_ITEM_NAME,
        []interface{}{
            "test4",
        },
        dynamo.ComparisonOperator_EQ)
    req.Select = dynamo.Select_ALL_ATTRIBUTES
    req.TableName = TEST_TABLE_NAME
    req.Host.Region = "us-west-2"
    req.Host.Domain = "amazonaws.com"
    req.Key, _ = awsgo.GetSecurityKeys()

    resp, err := req.Request()
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
    putRequest.Path = fmt.Sprintf("%s/haha/test.jpg", TEST_S3_BUCKET)
    fakePayload := "1234567890"
    putRequest.Length = int64(len(fakePayload))
    putRequest.Source = ioutil.NopCloser(bytes.NewBuffer([]byte(fakePayload)))

    putRequest.Host.Domain = "amazonaws.com"
    putRequest.Key, _ = awsgo.GetSecurityKeys()

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
    sendRequest.TaskQueue = fmt.Sprintf("/%s/TestQueue", TEST_SQS_QUEUE)

    sendRequest.Host.Region = "us-west-2"
    sendRequest.Host.Domain = "amazonaws.com"
    sendRequest.Key, _ = awsgo.GetSecurityKeys()

    resp, err := sendRequest.Request()
    fmt.Println(resp, err)

}

func TestSqsReceiveMessage() {
    sendRequest := sqs.NewReceiveMessageRequest()
    sendRequest.TaskQueue = fmt.Sprintf("/%s/TestQueue", TEST_SQS_QUEUE)
    sendRequest.WaitTimeSeconds = 10
    sendRequest.MaxNumberOfMessages = 1

    sendRequest.Host.Region = "us-west-2"
    sendRequest.Host.Domain = "amazonaws.com"
    sendRequest.Key, _ = awsgo.GetSecurityKeys()

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
    sendRequest.Key, _ = awsgo.GetSecurityKeys()

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
    putMetricRequest.Key, _ = awsgo.GetSecurityKeys()

    resp, err := putMetricRequest.Request()
    fmt.Println(resp, err)
}

func TestDescribeInstance() {
    describe := ec2.NewDescribeInstancesRequest()
    describe.InstanceIds =  []string{"i-ffffffff"}
    describe.Host.Region = "us-west-2"
    describe.Key, _ = awsgo.GetSecurityKeys()
    resp, err := describe.Request()
    fmt.Println(resp, err)
}


func TestGetMetricStatistics() {
    getMetric := cloudwatch.NewGetMetricStatisticsRequest()
    getMetric.MetricName = "CPUUtilization"
    getMetric.Namespace = "AWS/EC2"
    getMetric.Statistics = []string{cloudwatch.STATISTIC_Average}
    getMetric.StartTime = time.Now().Add(-12 * time.Hour)
    getMetric.EndTime = time.Now()
    getMetric.Host.Region = "us-west-2"
    getMetric.Key, _ = awsgo.GetSecurityKeys()
    resp, err := getMetric.Request()
    fmt.Println(resp, err)
}

func TestCreateLogGroup() {
    createGroup := cloudwatch.NewCreateLogGroupRequest()
    createGroup.LogGroupName = "test_group"
    createGroup.Key, _ = awsgo.GetSecurityKeys()
    _, err := createGroup.Request()
    if err != nil  {
        if v, ok := err.(*cloudwatch.ErrorResult); ok {
            if v.Type != cloudwatch.ResourceAlreadyExistsExceptionType {
                fmt.Printf("createGroup Error: %v %#v\n", err, err)
                return
            }
        } else {
            fmt.Printf("createGroup Error: %v %#v\n", err, err)
            return
        }
    }
    createStream := cloudwatch.NewCreateLogStreamRequest()
    createStream.LogStreamName = "test_stream"
    createStream.LogGroupName = "test_group"
    createStream.Key, _ = awsgo.GetSecurityKeys()
    if _, err := createStream.Request(); err != nil {
        if v, ok := err.(*cloudwatch.ErrorResult); ok {
            if v.Type != cloudwatch.ResourceAlreadyExistsExceptionType {
                fmt.Printf("createStream Error: %v %#v\n", err, err)
                return
            }
        } else {
            fmt.Printf("createStream Error: %v %#v\n", err, err)
            return
        }
    }
    putRetention := cloudwatch.NewPutRetentionPolicyRequest()
    putRetention.LogGroupName = "test_group"
    putRetention.RetentionInDays = 1
    putRetention.Key, _ = awsgo.GetSecurityKeys()
    if _, err := putRetention.Request(); err != nil {
        fmt.Printf("putRetention Error: %v %#v\n", err, err)
        return
    }
    nextToken := ""
    for {
        log1 := cloudwatch.NewPutLogEventsRequest()
        log1.AddEvent("Test 1", time.Now())
        log1.AddEvent("Test 2", time.Now())
        log1.LogGroupName = "test_group"
        log1.LogStreamName = "test_stream"
        log1.SequenceToken = nextToken
        log1.Key, _ = awsgo.GetSecurityKeys()
        log1Resp, err := log1.Request()
        if err != nil {
            if v, ok := err.(*cloudwatch.ErrorResult); ok {
                if v.Type != cloudwatch.InvalidSequenceTokenExceptionType && v.Type != cloudwatch.DataAlreadyAcceptedExceptionType {
                    fmt.Printf("createStream Error: %v %#v\n", err, err)
                    return
                }
                nextToken = cloudwatch.ExtractNextTokenFromError(v)
                fmt.Println("Extracted: ", nextToken)
            } else {
                fmt.Printf("createStream Error: %v %#v\n", err, err)
                return
            }
        } else {
            nextToken = log1Resp.NextSequenceToken
            break
        }
    }
    log2 := cloudwatch.NewPutLogEventsRequest()
    log2.SequenceToken = nextToken
    log2.AddEvent("Test 3", time.Now())
    log2.AddEvent("Test 4", time.Now())
    log2.LogGroupName = "test_group"
    log2.LogStreamName = "test_stream"
    log2.Key, _ = awsgo.GetSecurityKeys()
    log2Resp, err := log2.Request()
    if err != nil {
        fmt.Printf("log2 Error: %v %#v\n", err, err)
        return
    }
    fmt.Println("Next Token: ", log2Resp.NextSequenceToken)
}

func TestGetLogEvents() {
    get := cloudwatch.NewGetLogEventsRequest()
    get.LogGroupName = "test_group"
    get.LogStreamName = "test_stream"
    get.SetLimit(100)
    get.SetTimeRange(time.Now().Add(-24 * time.Minute), time.Now())
    get.Key, _ = awsgo.GetSecurityKeys()
    resp, err := get.Request()
    fmt.Println(resp, err)
}

func DescribeLogStuff() {
    group := cloudwatch.NewDescribeLogGroupsRequest()
    group.Key, _ = awsgo.GetSecurityKeys()
    groupResp, err := group.Request()
    if err != nil {
        fmt.Printf("describeLogGroups Error: %v %#v\n", err, err)
        return
    }
    fmt.Printf("LogGroups : %#v\n", groupResp)
    for i := range groupResp.LogGroups {
        stream := cloudwatch.NewDescribeLogStreamsRequest()
        stream.LogGroupName = groupResp.LogGroups[i].LogGroupName
        stream.Key, _ = awsgo.GetSecurityKeys()
        streamResp, err := stream.Request()
        if err != nil {
            fmt.Printf("describeLogStreams Error: %v %#v\n", err, err)
            return
        }
        fmt.Printf("LogStreams : %#v\n", streamResp)
    }
}


func pollForDecisionTask() {
    req := swf.NewPollForActivityTaskRequest()
    req.Domain = "test"
    req.Identity = "example"
    req.TaskList.Name = "hoho"
    req.Host.Region = "us-west-2"
    req.Key, _ = awsgo.GetSecurityKeys()
    r, err := req.Request()
    log.Println("Err: ", err)
    log.Println("Resp: ", r)
}

func postSnsMessage() {
    pub := sns.NewPublishRequest()
    pub.Host.Region = "us-west-2"
    pub.TopicArn = ""
    pub.Message = "Wasssssssuppppp"
    pub.Subject = "Howdy me!"
    resp, err := pub.Request()
    log.Println("Err: ", err)
    log.Println("Resp: ", resp)
}

func main() {
    //TestGetMetricStatistics()
    //TestCreateLogGroup()
    //TestGetLogEvents()
    //DescribeLogStuff()

    //TestGetItem()
    //TestUpdateItem()
    //TestPutItem()
    //TestDeleteItem()
    //TestBatchGetItem()
    //TestBatchWriteItem()
    //TestQuery()
    //TestPutS3File()
    //TestSqsSendMessage()
    //TestSqsReceiveMessage()
    //TestSesSendEmail()
    //TestPutMetric()
    //TestDescribeInstance()

    //pollForDecisionTask()

    postSnsMessage()
}
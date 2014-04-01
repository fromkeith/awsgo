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
    "time"
    "errors"
    "fmt"
    "sync"
)



var offThreadSendChannel chan *PutMetricRequest
var offThreadChannelLock sync.Mutex



type timedEvent struct {
    startTime time.Time
    name string
    namespace string
    onThisThread bool
}

func NewTimedEvent(name, namespace string, onThisThread bool) timedEvent {
    return timedEvent{
        time.Now(),
        name,
        namespace,
        onThisThread,
    }
}

func (t timedEvent) Report() {
    SimpleKeyValueMetric(
        t.name,
        float64(float64(time.Now().Sub(t.startTime).Nanoseconds()) / float64(time.Millisecond)),
        UNIT_MILLISECONDS,
        t.namespace,
        t.onThisThread,
    )
}



/**
 * Creates a single metric and posts it.
 * Sets the metric name, value, unit and namespace for the given parameters.
 * Also sets timestamp to time.Now()
 * if sendOnThisThread is set to true then this request will block to send it.
 *  Otherwise it will be pushed into a channel to be sent by a worker created
 *  in CreateOffThreadSender().
 */
func SimpleKeyValueMetric(name string, value float64, unit string, namespace string, sendOnThisThread bool) error {
    putMetricRequest := NewPutMetricRequest()
    putMetricRequest.Namespace = namespace
    putMetricRequest.MetricData = make([]MetricDatum, 1)
    putMetricRequest.MetricData[0].MetricName = name
    putMetricRequest.MetricData[0].Unit = unit
    putMetricRequest.MetricData[0].Value = new(float64)
    *(putMetricRequest.MetricData[0].Value) = value
    putMetricRequest.MetricData[0].Timestamp = new(time.Time)
    *(putMetricRequest.MetricData[0].Timestamp) = time.Now()

    putMetricRequest.Host.Region = "us-west-2"
    putMetricRequest.Host.Domain = "amazonaws.com"

    if sendOnThisThread {
        putMetricRequest.Key, _ = awsgo.GetSecurityKeys()
        _, err := putMetricRequest.Request()
        return err
    } else {
        if offThreadSendChannel == nil {
            return errors.New("No sender has been created! Failing.")
        }
        offThreadSendChannel <- putMetricRequest
        return nil
    }
}

func MultiKeyValueMetrics(name []string, value []float64, unit []string, namespace string, sendOnThisThread bool) error {
    putMetricRequest := NewPutMetricRequest()
    putMetricRequest.Namespace = namespace

    var size int

    if size = len(name); size > len(value) {
        size = len(value)
    }
    if size > len(unit) {
        size = len(unit)
    }
    putMetricRequest.MetricData = make([]MetricDatum, size)
    for i := 0; i < size; i++ {
        putMetricRequest.MetricData[i].MetricName = name[i]
        putMetricRequest.MetricData[i].Unit = unit[i]
        putMetricRequest.MetricData[i].Value = new(float64)
        *(putMetricRequest.MetricData[i].Value) = value[i]
        putMetricRequest.MetricData[i].Timestamp = new(time.Time)
        *(putMetricRequest.MetricData[i].Timestamp) = time.Now()
    }

    putMetricRequest.Host.Region = "us-west-2"
    putMetricRequest.Host.Domain = "amazonaws.com"


    if sendOnThisThread {
        putMetricRequest.Key, _ = awsgo.GetSecurityKeys()
        _, err := putMetricRequest.Request()
        return err
    } else {
        if offThreadSendChannel == nil {
            createOffThreadSenderIfNotExists()
        }
        offThreadSendChannel <- putMetricRequest
        return nil
    }
}


func createOffThreadSendChannelIfNotExists() {
    offThreadChannelLock.Lock()
    defer offThreadChannelLock.Unlock()
    if offThreadSendChannel == nil {
        // have a bigish buffer so we actaully are not blocking
        offThreadSendChannel = make(chan *PutMetricRequest, 500)
    }
}

func createOffThreadSenderIfNotExists() {
    if offThreadSendChannel == nil {
        CreateOffThreadSender()
    }
}

/** Creates a worker to send metrics.
 * Must be called at least once before trying to set 'sendOnThisThread' false for helper methods.
 */
func CreateOffThreadSender() {
    createOffThreadSendChannelIfNotExists()
    go func () {
        for {
            putMetricRequest := <- offThreadSendChannel
            if putMetricRequest == nil {
                break
            }
            putMetricRequest.Key, _ = awsgo.GetSecurityKeys()
            _, err := putMetricRequest.Request()
            if err != nil {
                fmt.Println(err)
            }
        }
    }()
}

func CloseOffThreadSender() {
    if offThreadSendChannel != nil {
        close(offThreadSendChannel)
    }
}
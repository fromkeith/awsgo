/*
 * Copyright (c) 2014, fromkeith
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


package swfhelper


import (
    "code.google.com/p/go-uuid/uuid"
    "errors"
    "fmt"
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/ec2"
    "github.com/fromkeith/awsgo/swf"
    "log"
    "runtime/debug"
    "time"
)


type ActivityHandler func(*ActivityContext)

type ActivityContext struct {
    pollTask            *swf.PollForActivityTaskResponse
    Input               string
    Activity            swf.ActivityType
    region              string
    heartbeatTimer      *time.Timer
    // when the next heartbeat is sent out, this channel will be checked for the last message.
    HeartbeatDetails    chan string
    CancelRequested     chan bool

    owningPool          chan bool
    marshaler           Marshaler
    heartbeatTime       time.Duration
}

type ActivityWorker struct {
    Domain              string
    // default is: ec2 instace Id + uuid
    Identity            string
    TaskList            string

    Region              string

    activityHandlers        map[string]ActivityHandler
    activityHeartbeats      map[string]time.Duration

    // leave blank to use default json marshaller
    Marshaler           Marshaler

    // maximum number of workers we can spawn. Default is infinite
    MaxWorkers          int
    workerPool          chan bool
}


func (a *ActivityWorker) RegisterActivity(actType swf.ActivityType, handler ActivityHandler) {
    if len(a.activityHandlers) == 0 {
        a.activityHandlers = make(map[string]ActivityHandler)
    }
    key := fmt.Sprintf("%s==>%s", actType.Name, actType.Version)
    a.activityHandlers[key] = handler
}

// overrides our default heartbeat interval of 1 minute
func (a *ActivityWorker) SetActivityHeartbeatInterval(actType swf.ActivityType, dur time.Duration) {
    if len(a.activityHeartbeats) == 0 {
        a.activityHeartbeats = make(map[string]time.Duration)
    }
    key := fmt.Sprintf("%s==>%s", actType.Name, actType.Version)
    a.activityHeartbeats[key] = dur
}


// starts polling for activitiies, indefinitely.
func (a *ActivityWorker) Start() error {
    if a.Identity == "" {
        ec2Identity, err := ec2.InstanceId()
        if err != nil {
            return err
        }
        a.Identity = fmt.Sprintf("%s-%s", ec2Identity, uuid.New())
    }
    if a.Marshaler == nil {
        a.Marshaler = JsonMarshaler{}
    }
    if a.MaxWorkers > 0 {
        a.workerPool = make(chan bool, a.MaxWorkers)
        // fill it in
        for i := 0; i < a.MaxWorkers; i++ {
            a.workerPool <- true
        }
    }
    for {
        a.startActivityGetting()
    }
}
func (a *ActivityWorker) startActivityGetting() {
    defer func() {
        rec := recover()
        if rec != nil {
            log.Println("Paniced when running activity worker. ", rec)
            debug.PrintStack()
        }
    }()
    for {
        if a.MaxWorkers > 0 {
            // wait for 1
            <- a.workerPool
        }
        poll := swf.NewPollForActivityTaskRequest()
        poll.Domain = a.Domain
        poll.Identity = a.Identity
        poll.TaskList = swf.TaskList{Name: a.TaskList}

        poll.Host.Region = a.Region
        poll.Key, _ = awsgo.GetSecurityKeys()

        resp, err := poll.Request()
        if err != nil {
            log.Println("Error making poll for activity task request.", err)
            time.Sleep(1 * time.Second)
            a.workerPool <- true
            continue
        }
        // no activity to work on
        if resp.ActivityId == "" {
            a.workerPool <- true
            continue
        }
        a.handleActivityRequest(resp)
    }
}

func (a *ActivityWorker) handleActivityRequest(resp *swf.PollForActivityTaskResponse) {
    key := fmt.Sprintf("%s==>%s", resp.ActivityType.Name, resp.ActivityType.Version)
    if h, ok := a.activityHandlers[key]; !ok {
        log.Panicf("Could not find activity handler for key: %v", key)
    } else {
        act := ActivityContext{
            pollTask: resp,
            Input: resp.Input,
            Activity: resp.ActivityType,
            region: a.Region,
            HeartbeatDetails: make(chan string),
            CancelRequested: make(chan bool),
            owningPool: a.workerPool,
            marshaler: a.Marshaler,
            heartbeatTime: time.Minute,
        }
        if nd, ok := a.activityHeartbeats[key]; ok {
            act.heartbeatTime = nd
        }
        act.restartHeartbeat()
        go func () {
            // catch any unexpected panics
            defer func () {
                rec := recover()
                if rec != nil {
                    defer func () {
                        rec2 := recover()
                        if rec2 != nil {
                            log.Println("Panic trying to fail activity", rec2)
                        }
                    }()
                    act.Failed("Panic", fmt.Sprintf("%v", rec))
                }
            }()
            h(&act)
        }()
    }
}


func (a *ActivityContext) restartHeartbeat() {
    a.heartbeatTimer = time.AfterFunc(a.heartbeatTime, a.heartbeat)
}

func (a *ActivityContext) heartbeat() {

    details := ""
    var ok bool

forloop:
    for {
        select {
            case details, ok = <- a.HeartbeatDetails:
                if !ok {
                    break forloop
                }
                break
            default:
                break forloop
        }
    }

    heart := swf.NewRespondActivityTaskHeartbeatRequest()
    heart.TaskToken = a.pollTask.TaskToken
    heart.Details = details
    heart.Host.Region = a.region
    heart.Key, _ = awsgo.GetSecurityKeys()
    resp, err := heart.Request()
    if err != nil {
        log.Println("Error sending heartbeat: ", err)
        return
    }
    if resp.CancelRequested {
        go func () {
            a.CancelRequested <- true
        }()
    }
    a.restartHeartbeat()
}


func (a *ActivityContext) recycle() {
    close(a.HeartbeatDetails)
    close(a.CancelRequested)
    if a.owningPool != nil {
        a.owningPool <- true
    }
}

// called if you want to stop processing an activity, and stop its heartbeat
// but don't want to respond with a result to this activity.
// you may want to do this if the activity requires human intervention, and
// you don't want to consume resources waiting for that interaction.
func (a *ActivityContext) CloseWithoutResult() {
    a.heartbeatTimer.Stop()
    a.recycle()
}

func (a * ActivityContext) GetRawTask() *swf.PollForActivityTaskResponse {
    return a.pollTask
}


// mark this activity as succesfully completed
func (a *ActivityContext) Completed(result interface{}) {
    res, _ := a.marshaler.Marshal(result)
    a.heartbeatTimer.Stop()
    defer a.recycle()
    for i := 0; ; i++ {
        if err := a.markCompletedRequest(string(res)); err != nil {
            if i > 10 {
                log.Panicf("Failed to mark activity as completed!: %v", err)
            }
            log.Printf("Failed to mark activity as completed! Retrying.: %v", err)
            time.Sleep(time.Duration(i * i) * 200 * time.Millisecond)
            continue
        }
        return
    }
}

func (a *ActivityContext) markCompletedRequest(result string) error {
    com := swf.NewRespondActivityTaskCompletedRequest()
    com.Result = result
    com.TaskToken = a.pollTask.TaskToken
    com.Host.Region = a.region
    com.Key, _ = awsgo.GetSecurityKeys()
    _, err := com.Request()
    return err
}

// mark this activity as Failed
func (a *ActivityContext) Failed(reason, details string) {
    a.heartbeatTimer.Stop()
    defer a.recycle()
    for i := 0; ; i++ {
        if err := a.markFailedRequest(reason, details); err != nil {
            if i > 10 {
                log.Panicf("Failed to mark activity as failed!: %v", err)
            }
            log.Printf("Failed to mark activity as failed! Retrying.: %v", err)
            time.Sleep(time.Duration(i * i) * 200 * time.Millisecond)
            continue
        }
        return
    }
}

func (a *ActivityContext) markFailedRequest(reason, details string) error {
    fal := swf.NewRespondActivityTaskFailedRequest()
    fal.Reason = reason
    fal.Details = details
    fal.TaskToken = a.pollTask.TaskToken
    fal.Host.Region = a.region
    fal.Key, _ = awsgo.GetSecurityKeys()
    _, err := fal.Request()
    return err
}

// mark this activity as succesfully canceled
func (a *ActivityContext) SuccesfullyCancel(details string) {
    a.heartbeatTimer.Stop()
    defer a.recycle()

    for i := 0; ; i++ {
        if err := a.markCanceledRequest(details); err != nil {
            if i > 10 {
                log.Panicf("Failed to mark activity as cancel!: %v", err)
            }
            log.Printf("Failed to mark activity as cancel! Retrying.: %v", err)
            time.Sleep(time.Duration(i * i) * 200 * time.Millisecond)
            continue
        }
        return
    }
}

func (a *ActivityContext) markCanceledRequest(details string) error {
    canc := swf.NewRespondActivityTaskCanceledRequest()
    canc.Details = details
    canc.TaskToken = a.pollTask.TaskToken
    canc.Host.Region = a.region
    canc.Key, _ = awsgo.GetSecurityKeys()
    _, err := canc.Request()
    return err
}

func (a *ActivityContext) As(out interface{}) error {
    if a.marshaler == nil {
        return errors.New("Nothing to unmarshal with")
    }
    return a.marshaler.Unmarshal(a.Input, out)
}
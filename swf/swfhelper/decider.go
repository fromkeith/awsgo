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
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/swf"
    "github.com/fromkeith/awsgo/ec2"
    "code.google.com/p/go-uuid/uuid"
    "fmt"
    "log"
    "runtime/debug"
    "time"
    "strconv"
)


// an instance of a workflow
type SwfWorkflow struct {
    nextActivityId      int64

    history             []swf.HistoryEvent
    decisions           []swf.Decision
    taskToken           string

    region              string

    owningPool          chan *SwfWorkflow
}
// called on new coroutines when a new workflow needs to be executed.
type WorkflowHandler func(w *SwfWorkflow)

// a channel to get task results from. If this channel is closed, the we do not have a result yet.
type TaskResultChan chan TaskResult

type Task struct {
    // which activity to execute
    Activity                    swf.ActivityType
    // specifies the maximum time before which a worker processing a task of this type must report progress
    HeartbeatTimeout            *int64
    // maximum duration for this activity task.
    ScheduleToCloseTimeout      *int64
    // specifies the maximum duration the activity task can wait to be assigned to a worker
    ScheduleToStartTimeout      *int64
    // specifies the maximum duration a worker may take to process this activity task
    StartToCloseTimeout         *int64
    // specifies the name of the task list in which to schedule the activity task.
    TaskList                    string
}

type TaskResult struct {
    activityId                  string

    Result                      string

    // the type of failure.
    //      ScheduleActivityTaskFailed
    //          FailureCause: see 'cause' http://docs.aws.amazon.com/amazonswf/latest/apireference/API_ScheduleActivityTaskFailedEventAttributes.html
    //      ActivityTaskCanceled
    //          FailureCause: 'details' provided in cancel request
    //      ActivityTaskTimedOut
    //          FailureCause: 'details' last provided by a heartbeat
    //      ActivityTaskFailed
    //          FailureCause: 'details' provided in the failure if any. For now we don't return 'reason'
    FailureType                 string
    FailureCause                string
}

// will poll for new decision tasks, and delegate them in new coroutines via the workflowhandler.
type Decider struct {
    Domain          string
    // default is: ec2 instace Id + uuid
    Identity        string
    TaskList        string

    Region          string
    // maximum number of workers we can spawn. Default is infinite
    MaxWorkers      int

    workerPool          chan *SwfWorkflow
    workflowHandlers    map[string]WorkflowHandler
}

func (s * SwfWorkflow) Go(do Task, data string) TaskResultChan {
    thisId := fmt.Sprintf("activityId-%d", s.nextActivityId)
    s.nextActivityId ++

    response := make(TaskResultChan)

    // for now, just a full search (for i := range s.history)...
    // but we should be able to use some heuristics here..
    // as we require Go calls to be determinitic in ordering

    // find if we have scheduled it yet, and if we have a result yet
    var scheduledEventId int64 = -1
    for i := range s.history {
        if s.history[i].EventType == "ActivityTaskScheduled" {
            if s.history[i].ActivityTaskScheduledEventAttributes.ActivityId == thisId {
                scheduledEventId = s.history[i].EventId
            }
            continue
        }
        if s.history[i].EventType == "ActivityTaskCompleted" {
            if s.history[i].ActivityTaskCompletedEventAttributes.ScheduledEventId == scheduledEventId {
                go func() {
                    response <- TaskResult{
                        activityId: thisId,
                        Result: s.history[i].ActivityTaskCompletedEventAttributes.Result,
                    }
                    close(response)
                }()
                return response
            }
            continue
        }
        if s.history[i].EventType == "ScheduleActivityTaskFailed" {
            if s.history[i].ScheduleActivityTaskFailedEventAttributes.ActivityId == thisId {
                go func() {
                    response <- TaskResult{
                        activityId: thisId,
                        FailureType: "ScheduleActivityTaskFailed",
                        FailureCause: s.history[i].ScheduleActivityTaskFailedEventAttributes.Cause,
                    }
                    close(response)
                }()
                return response
            }
            continue
        }
        if s.history[i].EventType == "ActivityTaskCanceled" {
            if s.history[i].ActivityTaskCanceledEventAttributes.ScheduledEventId == scheduledEventId {
                go func() {
                    response <- TaskResult{
                        activityId: thisId,
                        FailureType: "ActivityTaskCanceled",
                        FailureCause: s.history[i].ActivityTaskCanceledEventAttributes.Details,
                    }
                    close(response)
                }()
                return response
            }
            continue
        }
        if s.history[i].EventType == "ActivityTaskTimedOut" {
            if s.history[i].ActivityTaskTimedOutEventAttributes.ScheduledEventId == scheduledEventId {
                go func() {
                    response <- TaskResult{
                        activityId: thisId,
                        FailureType: "ActivityTaskTimedOut",
                        FailureCause: s.history[i].ActivityTaskTimedOutEventAttributes.Details,
                    }
                    close(response)
                }()
                return response
            }
            continue
        }
        if s.history[i].EventType == "ActivityTaskFailed" {
            if s.history[i].ActivityTaskFailedEventAttributes.ScheduledEventId == scheduledEventId {
                go func() {
                    response <- TaskResult{
                        activityId: thisId,
                        FailureType: "ActivityTaskFailed",
                        FailureCause: s.history[i].ActivityTaskFailedEventAttributes.Details,
                    }
                    close(response)
                }()
                return response
            }
            continue
        }
    }
    if scheduledEventId != -1 {
        close(response)
        return response
    }
    var asTaskList *swf.TaskList
    if do.TaskList != "" {
        asTaskList = &swf.TaskList{
            Name: do.TaskList,
        }
    }
    // we haven't schedule it yet.. so lets do that
    s.decisions = append(s.decisions,
        swf.Decision{
            DecisionType: "ScheduleActivityTask",
            ScheduleActivityTaskDecisionAttributes: &swf.ScheduleActivityTaskDecisionAttributes{
                ActivityId: thisId,
                ActivityType: do.Activity,
                Control: "", // ignoring for now
                HeartbeatTimeout: intPtrToString(do.HeartbeatTimeout),
                Input: data,
                ScheduleToCloseTimeout: intPtrToString(do.ScheduleToCloseTimeout),
                ScheduleToStartTimeout: intPtrToString(do.ScheduleToStartTimeout),
                StartToCloseTimeout: intPtrToString(do.StartToCloseTimeout),
                TaskList: asTaskList,
            },
        },
    )
    // close the response as we haven't executed this task yet
    close(response)
    return response
}

func intPtrToString(v *int64) string {
    if v == nil {
        return ""
    }
    return strconv.FormatInt(*v, 10)
}


func (d *Decider) RegisterWorkflow(workflow swf.WorkflowType, handler WorkflowHandler) {
    if len(d.workflowHandlers) == 0 {
        d.workflowHandlers = make(map[string]WorkflowHandler)
    }
    key := fmt.Sprintf("%s==>%s", workflow.Name, workflow.Version)
    d.workflowHandlers[key] = handler
}

func (d *Decider) newWorker() *SwfWorkflow {
    return &SwfWorkflow{
        history: make([]swf.HistoryEvent, 0, 10),
        decisions: make([]swf.Decision, 0, 5),
        nextActivityId: 0,
        region: d.Region,
        owningPool: d.workerPool,
    }
}

// starts polling for decisions, indefinitely.
func (d *Decider) Start() error {
    if d.Identity == "" {
        ec2Identity, err := ec2.InstanceId()
        if err != nil {
            return err
        }
        d.Identity = fmt.Sprintf("%s-%s", ec2Identity, uuid.New())
    }
    if d.MaxWorkers > 0 {
        d.workerPool = make(chan *SwfWorkflow, d.MaxWorkers)
        for i := 0; i < d.MaxWorkers; i++ {
            d.workerPool <- d.newWorker()
        }
    }
    for {
        d.startDeciding()
    }
}
func (d *Decider) startDeciding() {
    defer func() {
        rec := recover()
        if rec != nil {
            log.Println("Paniced when running decider. ", rec)
            debug.PrintStack()
        }
    }()
    for {
        worker := <- d.workerPool
        worker.history = worker.history[0:0]
        worker.nextActivityId = 0
        worker.decisions = worker.decisions[0:0]
        worker.taskToken = ""

        poll := swf.NewPollForDecisionTaskRequest()
        poll.Domain = d.Domain
        poll.Identity = d.Identity
        poll.TaskList = swf.TaskList{Name: d.TaskList}

        poll.Host.Region = d.Region
        poll.Key, _ = awsgo.GetSecurityKeys()

        resp, err := poll.Request()
        if err != nil {
            log.Println("Error making poll for decision request.", err)
            time.Sleep(1 * time.Second)
            d.workerPool <- worker
            continue
        }
        if resp.TaskToken == "" {
            d.workerPool <- worker
            continue
        }
        d.handleDecisionTaskResponse(resp, worker)
    }
}

func (d *Decider) handleDecisionTaskResponse(resp *swf.PollForDecisionTaskResponse, worker *SwfWorkflow) {
    worker.history = d.fillInHistory(resp, worker.history)

    key := fmt.Sprintf("%s==>%s", resp.WorkflowType.Name, resp.WorkflowType.Version)
    if h, ok := d.workflowHandlers[key]; !ok {
        log.Panicf("Could not find workflow handler for key: %v", key)
    } else {
        worker.taskToken = resp.TaskToken
        go h(worker)
    }
}

func (d *Decider) fillInHistory(lastResp *swf.PollForDecisionTaskResponse, events []swf.HistoryEvent) []swf.HistoryEvent {
    for lastResp != nil {
        events = append(events, lastResp.Events...)

        if lastResp.NextPageToken == "" {
            break
        }

        poll := swf.NewPollForDecisionTaskRequest()
        poll.Domain = d.Domain
        poll.Identity = d.Identity
        poll.TaskList = swf.TaskList{Name: d.TaskList}
        poll.NextPageToken = lastResp.NextPageToken

        poll.Host.Region = d.Region
        poll.Key, _ = awsgo.GetSecurityKeys()

        for i := 0; ; i++ {
            resp, err := poll.Request()
            if err != nil {
                log.Println("Error making poll for decision request.", err)
                if i > 10 {
                    panic("Failed too many times to try get decidion task request.")
                }
                time.Sleep(time.Duration(i * i) * 200 * time.Millisecond)
                continue
            }
            lastResp = resp
            break
        }
    }
    return events
}

// Must be called for our decisions to be posted to the server
// in your handler you should always defer Decide() so that it gets executed.
func (w * SwfWorkflow) Decide() {
    // return to the pool
    defer func () {
        w.owningPool <- w
    }()

    for i := 0; ; i ++ {
        req := swf.NewRespondDecisionTaskCompletedRequest()
        req.Decisions = w.decisions
        req.TaskToken = w.taskToken
        req.Host.Region = w.region
        req.Key, _ = awsgo.GetSecurityKeys()
        _, err := req.Request()
        if err != nil {
            if i > 10 {
                log.Panicf("Error responding to decision task completed. Retrying", err)
            }
            log.Println("Error responding to decision task completed. Retrying", err)
            time.Sleep(time.Duration(i * i) * 100 * time.Millisecond)
            continue
        }
        return
    }
}

func (s* SwfWorkflow) Complete(result string) {
    s.decisions = append(s.decisions,
        swf.Decision{
            DecisionType: "CompleteWorkflowExecution",
            CompleteWorkflowExecutionDecisionAttributes: &swf.CompleteWorkflowExecutionDecisionAttributes{
                Result: result,
            },
        },
    )
}

func (s* SwfWorkflow) Fail(reason, details string) {
    s.decisions = append(s.decisions,
        swf.Decision{
            DecisionType: "FailWorkflowExecution",
            FailWorkflowExecutionDecisionAttributes: &swf.FailWorkflowExecutionDecisionAttributes{
                Details: details,
                Reason: reason,
            },
        },
    )
}
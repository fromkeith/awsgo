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

package swf

import (
	"encoding/json"
	"github.com/fromkeith/awsgo"
	"log"
)

type TaskList struct {
	Name string `json:"name"`
}

type WorkflowExecution struct {
	RunId      string `json:"runId"`
	WorkflowId string `json:"workflowId"`
}

type WorkflowType struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ActivityType struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type activityTaskCancelRequestedEventAttributes struct {
	ActivityId                   string `json:"activityId"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
}
type activityTaskCanceledEventAttributes struct {
	Details                      string `json:"details"`
	LatestCancelRequestedEventId int64  `json:"latestCancelRequestedEventId"`
	ScheduledEventId             int64  `json:"scheduledEventId"`
	StartedEventId               int64  `json:"startedEventId"`
}
type activityTaskCompletedEventAttributes struct {
	Result           string `json:"result"`
	ScheduledEventId int64  `json:"scheduledEventId"`
	StartedEventId   int64  `json:"startedEventId"`
}
type activityTaskFailedEventAttributes struct {
	Details          string `json:"details"`
	Reason           string `json:"reason"`
	ScheduledEventId int64  `json:"scheduledEventId"`
	StartedEventId   int64  `json:"startedEventId"`
}
type activityTaskScheduledEventAttributes struct {
	ActivityId                   string       `json:"activityId"`
	ActivityType                 ActivityType `json:"activityType"`
	Control                      string       `json:"control"`
	DecisionTaskCompletedEventId int64        `json:"decisionTaskCompletedEventId"`
	HeartbeatTimeout             string       `json:"heartbeatTimeout"`
	Input                        string       `json:"input"`
	ScheduleToCloseTimeout       string       `json:"scheduleToCloseTimeout"`
	ScheduleToStartTimeout       string       `json:"scheduleToStartTimeout"`
	StartToCloseTimeout          string       `json:"startToCloseTimeout"`
	TaskList                     TaskList      `json:"taskList"`
}
type activityTaskStartedEventAttributes struct {
	Identity         string `json:"identity"`
	ScheduledEventId int64  `json:"scheduledEventId"`
}
type activityTaskTimedOutEventAttributes struct {
	Details          string `json:"details"`
	ScheduledEventId int64  `json:"scheduledEventId"`
	StartedEventId   int64  `json:"startedEventId"`
	TimeoutType      string `json:"timeoutType"`
}
type cancelTimerFailedEventAttributes struct {
	Cause                        string `json:"cause"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	TimerId                      string `json:"timerId"`
}
type cancelWorkflowExecutionFailedEventAttributes struct {
	Cause                        string `json:"cause"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
}
type childWorkflowExecutionCanceledEventAttributes struct {
	Details           string `json:"details"`
	InitiatedEventId  int64  `json:"initiatedEventId"`
	StartedEventId    int64  `json:"startedEventId"`
	WorkflowExecution WorkflowExecution `json:"workflowExecution"`
	WorkflowType WorkflowType `json:"workflowType"`
}
type childWorkflowExecutionCompletedEventAttributes struct {
	InitiatedEventId  int64             `json:"initiatedEventId"`
	Result            string            `json:"result"`
	StartedEventId    int64             `json:"startedEventId"`
	WorkflowExecution WorkflowExecution `json:"workflowExecution"`
	WorkflowType      WorkflowType      `json:"workflowType"`
}
type childWorkflowExecutionFailedEventAttributes struct {
	Details           string            `json:"details"`
	InitiatedEventId  int64             `json:"initiatedEventId"`
	Reason            string            `json:"reason"`
	StartedEventId    int64             `json:"startedEventId"`
	WorkflowExecution WorkflowExecution `json:"workflowExecution"`
	WorkflowType      WorkflowType      `json:"workflowType"`
}
type childWorkflowExecutionStartedEventAttributes struct {
	InitiatedEventId  int64             `json:"initiatedEventId"`
	WorkflowExecution WorkflowExecution `json:"workflowExecution"`
	WorkflowType      WorkflowType      `json:"workflowType"`
}
type childWorkflowExecutionTerminatedEventAttributes struct {
	InitiatedEventId  int64             `json:"initiatedEventId"`
	StartedEventId    int64             `json:"startedEventId"`
	WorkflowExecution WorkflowExecution `json:"workflowExecution"`
	WorkflowType      WorkflowType      `json:"workflowType"`
}
type childWorkflowExecutionTimedOutEventAttributes struct {
	InitiatedEventId  int64             `json:"initiatedEventId"`
	StartedEventId    int64             `json:"startedEventId"`
	TimeoutType       string            `json:"timeoutType"`
	WorkflowExecution WorkflowExecution `json:"workflowExecution"`
	WorkflowType      WorkflowType      `json:"workflowType"`
}
type completeWorkflowExecutionFailedEventAttributes struct {
	Cause                        string `json:"cause"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
}
type continueAsNewWorkflowExecutionFailedEventAttributes struct {
	Cause                        string `json:"cause"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
}
type decisionTaskCompletedEventAttributes struct {
	ExecutionContext string `json:"executionContext"`
	ScheduledEventId int64  `json:"scheduledEventId"`
	StartedEventId   int64  `json:"startedEventId"`
}
type decisionTaskScheduledEventAttributes struct {
	StartToCloseTimeout string   `json:"startToCloseTimeout"`
	TaskList            TaskList `json:"taskList"`
}
type decisionTaskStartedEventAttributes struct {
	Identity         string `json:"identity"`
	ScheduledEventId int64  `json:"scheduledEventId"`
}
type decisionTaskTimedOutEventAttributes struct {
	ScheduledEventId int64  `json:"scheduledEventId"`
	StartedEventId   int64  `json:"startedEventId"`
	TimeoutType      string `json:"timeoutType"`
}

type externalWorkflowExecutionCancelRequestedEventAttributes struct {
	InitiatedEventId  int64             `json:"initiatedEventId"`
	WorkflowExecution WorkflowExecution `json:"workflowExecution"`
}
type externalWorkflowExecutionSignaledEventAttributes struct {
	InitiatedEventId  int64             `json:"initiatedEventId"`
	WorkflowExecution WorkflowExecution `json:"workflowExecution"`
}
type failWorkflowExecutionFailedEventAttributes struct {
	Cause                        string `json:"cause"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
}
type markerRecordedEventAttributes struct {
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	Details                      string `json:"details"`
	MarkerName                   string `json:"markerName"`
}
type recordMarkerFailedEventAttributes struct {
	Cause                        string `json:"cause"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	MarkerName                   string `json:"markerName"`
}
type requestCancelActivityTaskFailedEventAttributes struct {
	ActivityId                   string `json:"activityId"`
	Cause                        string `json:"cause"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
}
type requestCancelExternalWorkflowExecutionFailedEventAttributes struct {
	Cause                        string `json:"cause"`
	Control                      string `json:"control"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	InitiatedEventId             int64  `json:"initiatedEventId"`
	RunId                        string `json:"runId"`
	WorkflowId                   string `json:"workflowId"`
}
type requestCancelExternalWorkflowExecutionInitiatedEventAttributes struct {
	Control                      string `json:"control"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	RunId                        string `json:"runId"`
	WorkflowId                   string `json:"workflowId"`
}
type scheduleActivityTaskFailedEventAttributes struct {
	ActivityId                   string       `json:"activityId"`
	ActivityType                 ActivityType `json:"activityType"`
	Cause                        string       `json:"cause"`
	DecisionTaskCompletedEventId int64        `json:"decisionTaskCompletedEventId"`
}
type signalExternalWorkflowExecutionFailedEventAttributes struct {
	Cause                        string `json:"cause"`
	Control                      string `json:"control"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	InitiatedEventId             int64  `json:"initiatedEventId"`
	RunId                        string `json:"runId"`
	WorkflowId                   string `json:"workflowId"`
}
type signalExternalWorkflowExecutionInitiatedEventAttributes struct {
	Control                      string `json:"control"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	Input                        string `json:"input"`
	RunId                        string `json:"runId"`
	SignalName                   string `json:"signalName"`
	WorkflowId                   string `json:"workflowId"`
}
type startChildWorkflowExecutionFailedEventAttributes struct {
	Cause                        string       `json:"cause"`
	Control                      string       `json:"control"`
	DecisionTaskCompletedEventId int64        `json:"decisionTaskCompletedEventId"`
	InitiatedEventId             int64        `json:"initiatedEventId"`
	WorkflowId                   string       `json:"workflowId"`
	WorkflowType                 WorkflowType `json:"workflowType"`
}
type startChildWorkflowExecutionInitiatedEventAttributes struct {
	ChildPolicy                  string       `json:"childPolicy"`
	Control                      string       `json:"control"`
	DecisionTaskCompletedEventId int64        `json:"decisionTaskCompletedEventId"`
	ExecutionStartToCloseTimeout string       `json:"executionStartToCloseTimeout"`
	Input                        string       `json:"input"`
	TagList                      []string     `json:"tagList"`
	TaskList                     TaskList     `json:"taskList"`
	TaskStartToCloseTimeout      string       `json:"taskStartToCloseTimeout"`
	WorkflowId                   string       `json:"workflowId"`
	WorkflowType                 WorkflowType `json:"workflowType"`
}
type startTimerFailedEventAttributes struct {
	Cause                        string `json:"cause"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	TimerId                      string `json:"timerId"`
}
type timerCanceledEventAttributes struct {
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	StartedEventId               int64  `json:"startedEventId"`
	TimerId                      string `json:"timerId"`
}
type timerFiredEventAttributes struct {
	StartedEventId int64  `json:"startedEventId"`
	TimerId        string `json:"timerId"`
}
type timerStartedEventAttributes struct {
	Control                      string `json:"control"`
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	StartToFireTimeout           string `json:"startToFireTimeout"`
	TimerId                      string `json:"timerId"`
}
type workflowExecutionCancelRequestedEventAttributes struct {
	Cause                     string            `json:"cause"`
	ExternalInitiatedEventId  int64             `json:"externalInitiatedEventId"`
	ExternalWorkflowExecution WorkflowExecution `json:"externalWorkflowExecution"`
}
type workflowExecutionCanceledEventAttributes struct {
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	Details                      string `json:"details"`
}
type workflowExecutionCompletedEventAttributes struct {
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	Result                       string `json:"result"`
}
type workflowExecutionContinuedAsNewEventAttributes struct {
	ChildPolicy                  string       `json:"childPolicy"`
	DecisionTaskCompletedEventId int64        `json:"decisionTaskCompletedEventId"`
	ExecutionStartToCloseTimeout string       `json:"executionStartToCloseTimeout"`
	Input                        string       `json:"input"`
	NewExecutionRunId            string       `json:"newExecutionRunId"`
	TagList                      []string     `json:"tagList"`
	TaskList                     TaskList     `json:"taskList"`
	TaskStartToCloseTimeout      string       `json:"taskStartToCloseTimeout"`
	WorkflowType                 WorkflowType `json:"workflowType"`
}
type workflowExecutionFailedEventAttributes struct {
	DecisionTaskCompletedEventId int64  `json:"decisionTaskCompletedEventId"`
	Details                      string `json:"details"`
	Reason                       string `json:"reason"`
}
type workflowExecutionSignaledEventAttributes struct {
	ExternalInitiatedEventId  int64             `json:"externalInitiatedEventId"`
	ExternalWorkflowExecution WorkflowExecution `json:"externalWorkflowExecution"`
	Input                     string            `json:"input"`
	SignalName                string            `json:"signalName"`
}
type workflowExecutionStartedEventAttributes struct {
	ChildPolicy                  string            `json:"childPolicy"`
	ContinuedExecutionRunId      string            `json:"continuedExecutionRunId"`
	ExecutionStartToCloseTimeout string            `json:"executionStartToCloseTimeout"`
	Input                        string            `json:"input"`
	ParentInitiatedEventId       int64             `json:"parentInitiatedEventId"`
	ParentWorkflowExecution      WorkflowExecution `json:"parentWorkflowExecution"`
	TagList                      []string          `json:"tagList"`
	TaskList                     TaskList          `json:"taskList"`
	TaskStartToCloseTimeout      string            `json:"taskStartToCloseTimeout"`
	WorkflowType                 WorkflowType      `json:"workflowType"`
}
type workflowExecutionTerminatedEventAttributes struct {
	Cause       string `json:"cause"`
	ChildPolicy string `json:"childPolicy"`
	Details     string `json:"details"`
	Reason      string `json:"reason"`
}
type workflowExecutionTimedOutEventAttributes struct {
	ChildPolicy string `json:"childPolicy"`
	TimeoutType string `json:"timeoutType"`
}

type HistoryEvent struct {
	EventId        int64   `json:"eventId"`
	EventTimestamp float64 `json:"eventTimestamp"`
	EventType      string  `json:"eventType"`

	ActivityTaskCanceledEventAttributes                            *activityTaskCanceledEventAttributes                            `json:"activityTaskCanceledEventAttributes,omitempty"`
	ActivityTaskCancelRequestedEventAttributes                     *activityTaskCancelRequestedEventAttributes                     `json:"activityTaskCancelRequestedEventAttributes,omitempty"`
	ActivityTaskCompletedEventAttributes                           *activityTaskCompletedEventAttributes                           `json:"activityTaskCompletedEventAttributes,omitempty"`
	ActivityTaskFailedEventAttributes                              *activityTaskFailedEventAttributes                              `json:"activityTaskFailedEventAttributes,omitempty"`
	ActivityTaskScheduledEventAttributes                           *activityTaskScheduledEventAttributes                           `json:"activityTaskScheduledEventAttributes,omitempty"`
	ActivityTaskStartedEventAttributes                             *activityTaskStartedEventAttributes                             `json:"activityTaskStartedEventAttributes,omitempty"`
	ActivityTaskTimedOutEventAttributes                            *activityTaskTimedOutEventAttributes                            `json:"activityTaskTimedOutEventAttributes,omitempty"`
	CancelTimerFailedEventAttributes                               *cancelTimerFailedEventAttributes                               `json:"cancelTimerFailedEventAttributes,omitempty"`
	ChildWorkflowExecutionCanceledEventAttributes                  *childWorkflowExecutionCanceledEventAttributes                  `json:"childWorkflowExecutionCanceledEventAttributes,omitempty"`
	ChildWorkflowExecutionCompletedEventAttributes                 *childWorkflowExecutionCompletedEventAttributes                 `json:"childWorkflowExecutionCompletedEventAttributes,omitempty"`
	ChildWorkflowExecutionFailedEventAttributes                    *childWorkflowExecutionFailedEventAttributes                    `json:"childWorkflowExecutionFailedEventAttributes,omitempty"`
	ChildWorkflowExecutionStartedEventAttributes                   *childWorkflowExecutionStartedEventAttributes                   `json:"childWorkflowExecutionStartedEventAttributes,omitempty"`
	ChildWorkflowExecutionTerminatedEventAttributes                *childWorkflowExecutionTerminatedEventAttributes                `json:"childWorkflowExecutionTerminatedEventAttributes,omitempty"`
	ChildWorkflowExecutionTimedOutEventAttributes                  *childWorkflowExecutionTimedOutEventAttributes                  `json:"childWorkflowExecutionTimedOutEventAttributes,omitempty"`
	DecisionTaskCompletedEventAttributes                           *decisionTaskCompletedEventAttributes                           `json:"decisionTaskCompletedEventAttributes,omitempty"`
	DecisionTaskScheduledEventAttributes                           *decisionTaskScheduledEventAttributes                           `json:"decisionTaskScheduledEventAttributes,omitempty"`
	DecisionTaskStartedEventAttributes                             *decisionTaskStartedEventAttributes                             `json:"decisionTaskStartedEventAttributes,omitempty"`
	DecisionTaskTimedOutEventAttributes                            *decisionTaskTimedOutEventAttributes                            `json:"decisionTaskTimedOutEventAttributes,omitempty"`
	ExternalWorkflowExecutionCancelRequestedEventAttributes        *externalWorkflowExecutionCancelRequestedEventAttributes        `json:"externalWorkflowExecutionCancelRequestedEventAttributes,omitempty"`
	ExternalWorkflowExecutionSignaledEventAttributes               *externalWorkflowExecutionSignaledEventAttributes               `json:"externalWorkflowExecutionSignaledEventAttributes,omitempty"`
	MarkerRecordedEventAttributes                                  *markerRecordedEventAttributes                                  `json:"markerRecordedEventAttributes,omitempty"`
	RequestCancelActivityTaskFailedEventAttributes                 *requestCancelActivityTaskFailedEventAttributes                 `json:"requestCancelActivityTaskFailedEventAttributes,omitempty"`
	RequestCancelExternalWorkflowExecutionFailedEventAttributes    *requestCancelExternalWorkflowExecutionFailedEventAttributes    `json:"requestCancelExternalWorkflowExecutionFailedEventAttributes,omitempty"`
	RequestCancelExternalWorkflowExecutionInitiatedEventAttributes *requestCancelExternalWorkflowExecutionInitiatedEventAttributes `json:"requestCancelExternalWorkflowExecutionInitiatedEventAttributes,omitempty"`
	ScheduleActivityTaskFailedEventAttributes                      *scheduleActivityTaskFailedEventAttributes                      `json:"scheduleActivityTaskFailedEventAttributes,omitempty"`
	SignalExternalWorkflowExecutionFailedEventAttributes           *signalExternalWorkflowExecutionFailedEventAttributes           `json:"signalExternalWorkflowExecutionFailedEventAttributes,omitempty"`
	SignalExternalWorkflowExecutionInitiatedEventAttributes        *signalExternalWorkflowExecutionInitiatedEventAttributes        `json:"signalExternalWorkflowExecutionInitiatedEventAttributes,omitempty"`
	StartChildWorkflowExecutionFailedEventAttributes               *startChildWorkflowExecutionFailedEventAttributes               `json:"startChildWorkflowExecutionFailedEventAttributes,omitempty"`
	StartChildWorkflowExecutionInitiatedEventAttributes            *startChildWorkflowExecutionInitiatedEventAttributes            `json:"startChildWorkflowExecutionInitiatedEventAttributes,omitempty"`
	StartTimerFailedEventAttributes                                *startTimerFailedEventAttributes                                `json:"startTimerFailedEventAttributes,omitempty"`
	TimerCanceledEventAttributes                                   *timerCanceledEventAttributes                                   `json:"timerCanceledEventAttributes,omitempty"`
	TimerFiredEventAttributes                                      *timerFiredEventAttributes                                      `json:"timerFiredEventAttributes,omitempty"`
	TimerStartedEventAttributes                                    *timerStartedEventAttributes                                    `json:"timerStartedEventAttributes,omitempty"`
	WorkflowExecutionCanceledEventAttributes                       *workflowExecutionCanceledEventAttributes                       `json:"workflowExecutionCanceledEventAttributes,omitempty"`
	WorkflowExecutionCancelRequestedEventAttributes                *workflowExecutionCancelRequestedEventAttributes                `json:"workflowExecutionCancelRequestedEventAttributes,omitempty"`
	WorkflowExecutionCompletedEventAttributes                      *workflowExecutionCompletedEventAttributes                      `json:"workflowExecutionCompletedEventAttributes,omitempty"`
	WorkflowExecutionContinuedAsNewEventAttributes                 *workflowExecutionContinuedAsNewEventAttributes                 `json:"workflowExecutionContinuedAsNewEventAttributes,omitempty"`
	WorkflowExecutionFailedEventAttributes                         *workflowExecutionFailedEventAttributes                         `json:"workflowExecutionFailedEventAttributes,omitempty"`
	WorkflowExecutionSignaledEventAttributes                       *workflowExecutionSignaledEventAttributes                       `json:"workflowExecutionSignaledEventAttributes,omitempty"`
	WorkflowExecutionStartedEventAttributes                        *workflowExecutionStartedEventAttributes                        `json:"workflowExecutionStartedEventAttributes,omitempty"`
	WorkflowExecutionTerminatedEventAttributes                     *workflowExecutionTerminatedEventAttributes                     `json:"workflowExecutionTerminatedEventAttributes,omitempty"`
	WorkflowExecutionTimedOutEventAttributes                       *workflowExecutionTimedOutEventAttributes                       `json:"workflowExecutionTimedOutEventAttributes,omitempty"`
}

type PollForDecisionTaskRequest struct {
	awsgo.RequestBuilder

	Domain          string   `json:"domain"`
	Identity        string   `json:"identity"`
	MaximumPageSize int      `json:"maximumPageSize"`
	NextPageToken   string   `json:"nextPageToken,omitempty"`
	ReverseOrder    bool     `json:"reverseOrder,omitempty"`
	TaskList        TaskList `json:"taskList"`
}

type PollForDecisionTaskResponse struct {
	Events                 []HistoryEvent    `json:"events"`
	NextPageToken          string            `json:"nextPageToken"`
	PreviousStartedEventId int64             `json:"previousStartedEventId"`
	StartedEventId         int64             `json:"startedEventId"`
	TaskToken              string            `json:"taskToken"`
	WorkflowExecution      WorkflowExecution `json:"workflowExecution"`
	WorkflowType           WorkflowType      `json:"workflowType"`
}

func NewPollForDecisionTaskRequest() *PollForDecisionTaskRequest {
	req := new(PollForDecisionTaskRequest)
	req.MaximumPageSize = 100
	req.Host.Service = "swf"
	req.Host.Region = ""
	req.Host.Domain = "amazonaws.com"
	req.Key.AccessKeyId = ""
	req.Key.SecretAccessKey = ""
	req.Headers = make(map[string]string)
	req.Headers["X-Amz-Target"] = "SimpleWorkflowService.PollForDecisionTask"
	req.RequestMethod = "POST"
	req.CanonicalUri = "/"
	return req
}

func (req *PollForDecisionTaskRequest) VerifyInput() error {
	return nil
}

func (req PollForDecisionTaskRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) interface{} {
	log.Println("response: ", string(response))
	resp := new(PollForDecisionTaskResponse)
	if err := json.Unmarshal(response, resp); err != nil {
		return err
	}
	return resp
}

func (req PollForDecisionTaskRequest) Request() (*PollForDecisionTaskResponse, error) {
	request, err := awsgo.NewAwsRequest(&req, req)
	if err != nil {
		return nil, err
	}
	request.RequestSigningType = awsgo.RequestSigningType_AWS3
	resp, err := request.DoAndDemarshall(&req)
	if resp == nil {
		return nil, err
	}
	return resp.(*PollForDecisionTaskResponse), err
}

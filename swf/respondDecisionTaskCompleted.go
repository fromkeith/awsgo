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
	"github.com/fromkeith/awsgo"
	"log"
	"errors"
)

type CancelTimerDecisionAttributes struct {
	TimerId string `json:"timerId"`
}
type CancelWorkflowExecutionDecisionAttributes struct {
	Details string `json:"details"`
}
type CompleteWorkflowExecutionDecisionAttributes struct {
	Result string `json:"result"`
}
type ContinueAsNewWorkflowExecutionDecisionAttributes struct {
	ChildPolicy                  string   `json:"childPolicy"`
	ExecutionStartToCloseTimeout string   `json:"executionStartToCloseTimeout"`
	Input                        string   `json:"input"`
	TagList                      []string `json:"tagList"`
	TaskList                     TaskList `json:"taskList"`
	TaskStartToCloseTimeout      string   `json:"taskStartToCloseTimeout"`
	WorkflowTypeVersion          string   `json:"workflowTypeVersion"`
}
type FailWorkflowExecutionDecisionAttributes struct {
	Details string `json:"details"`
	Reason  string `json:"reason"`
}
type RecordMarkerDecisionAttributes struct {
	Details    string `json:"details"`
	MarkerName string `json:"markerName"`
}
type RequestCancelActivityTaskDecisionAttributes struct {
	ActivityId string `json:"activityId"`
}
type RequestCancelExternalWorkflowExecutionDecisionAttributes struct {
	Control    string `json:"control"`
	RunId      string `json:"runId"`
	WorkflowId string `json:"workflowId"`
}
type ScheduleActivityTaskDecisionAttributes struct {
	ActivityId             string       `json:"activityId"`
	ActivityType           ActivityType `json:"activityType"`
	Control                string       `json:"control,omitempty"`
	HeartbeatTimeout       string       `json:"heartbeatTimeout,omitempty"`
	Input                  string       `json:"input,omitempty"`
	ScheduleToCloseTimeout string       `json:"scheduleToCloseTimeout,omitempty"`
	ScheduleToStartTimeout string       `json:"scheduleToStartTimeout,omitempty"`
	StartToCloseTimeout    string       `json:"startToCloseTimeout,omitempty"`
	TaskList               *TaskList     `json:"taskList,omitempty"`
}
type SignalExternalWorkflowExecutionDecisionAttributes struct {
	Control    string `json:"control"`
	Input      string `json:"input"`
	RunId      string `json:"runId"`
	SignalName string `json:"signalName"`
	WorkflowId string `json:"workflowId"`
}
type StartChildWorkflowExecutionDecisionAttributes struct {
	ChildPolicy                  string       `json:"childPolicy"`
	Control                      string       `json:"control"`
	ExecutionStartToCloseTimeout string       `json:"executionStartToCloseTimeout"`
	Input                        string       `json:"input"`
	TagList                      []string     `json:"tagList"`
	TaskList                     TaskList     `json:"taskList"`
	TaskStartToCloseTimeout      string       `json:"taskStartToCloseTimeout"`
	WorkflowId                   string       `json:"workflowId"`
	WorkflowType                 WorkflowType `json:"workflowType"`
}
type StartTimerDecisionAttributes struct {
	Control            string `json:"control"`
	StartToFireTimeout string `json:"startToFireTimeout"`
	TimerId            string `json:"timerId"`
}

type Decision struct {
	DecisionType string `json:"decisionType"`

	CancelTimerDecisionAttributes                            *CancelTimerDecisionAttributes                            `json:"cancelTimerDecisionAttributes,omitempty"`
	CancelWorkflowExecutionDecisionAttributes                *CancelWorkflowExecutionDecisionAttributes                `json:"cancelWorkflowExecutionDecisionAttributes,omitempty"`
	CompleteWorkflowExecutionDecisionAttributes              *CompleteWorkflowExecutionDecisionAttributes              `json:"completeWorkflowExecutionDecisionAttributes,omitempty"`
	ContinueAsNewWorkflowExecutionDecisionAttributes         *ContinueAsNewWorkflowExecutionDecisionAttributes         `json:"continueAsNewWorkflowExecutionDecisionAttributes,omitempty"`
	FailWorkflowExecutionDecisionAttributes                  *FailWorkflowExecutionDecisionAttributes                  `json:"failWorkflowExecutionDecisionAttributes,omitempty"`
	RecordMarkerDecisionAttributes                           *RecordMarkerDecisionAttributes                           `json:"recordMarkerDecisionAttributes,omitempty"`
	RequestCancelActivityTaskDecisionAttributes              *RequestCancelActivityTaskDecisionAttributes              `json:"requestCancelActivityTaskDecisionAttributes,omitempty"`
	RequestCancelExternalWorkflowExecutionDecisionAttributes *RequestCancelExternalWorkflowExecutionDecisionAttributes `json:"requestCancelExternalWorkflowExecutionDecisionAttributes,omitempty"`
	ScheduleActivityTaskDecisionAttributes                   *ScheduleActivityTaskDecisionAttributes                   `json:"scheduleActivityTaskDecisionAttributes,omitempty"`
	SignalExternalWorkflowExecutionDecisionAttributes        *SignalExternalWorkflowExecutionDecisionAttributes        `json:"signalExternalWorkflowExecutionDecisionAttributes,omitempty"`
	StartChildWorkflowExecutionDecisionAttributes            *StartChildWorkflowExecutionDecisionAttributes            `json:"startChildWorkflowExecutionDecisionAttributes,omitempty"`
	StartTimerDecisionAttributes                             *StartTimerDecisionAttributes                             `json:"startTimerDecisionAttributes,omitempty"`
}

type RespondDecisionTaskCompletedRequest struct {
	awsgo.RequestBuilder

	Decisions        []Decision `json:"decisions"`
	ExecutionContext string     `json:"executionContext"`
	TaskToken        string     `json:"taskToken"`
}

type RespondDecisionTaskCompletedResponse struct {
}



func NewRespondDecisionTaskCompletedRequest() *RespondDecisionTaskCompletedRequest {
	req := new(RespondDecisionTaskCompletedRequest)
	req.Host.Service = "swf"
	req.Host.Region = ""
	req.Host.Domain = "amazonaws.com"
	req.Key.AccessKeyId = ""
	req.Key.SecretAccessKey = ""
	req.Headers = make(map[string]string)
	req.Headers["X-Amz-Target"] = "SimpleWorkflowService.RespondDecisionTaskCompleted"
	req.RequestMethod = "POST"
	req.CanonicalUri = "/"
	return req
}

func (req *RespondDecisionTaskCompletedRequest) VerifyInput() error {
	return nil
}

func (req RespondDecisionTaskCompletedRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) interface{} {
	if statusCode != 200 {
		log.Println("response: ", string(response))
		return errors.New("Bad response code!")
	}
	return new(RespondDecisionTaskCompletedResponse)
}

func (req RespondDecisionTaskCompletedRequest) Request() (*RespondDecisionTaskCompletedResponse, error) {
	request, err := awsgo.NewAwsRequest(&req, req)
	if err != nil {
		return nil, err
	}
	request.RequestSigningType = awsgo.RequestSigningType_AWS3
	resp, err := request.DoAndDemarshall(&req)
	if resp == nil {
		return nil, err
	}
	return resp.(*RespondDecisionTaskCompletedResponse), err
}
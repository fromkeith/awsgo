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

type PollForActivityTaskRequest struct {
	awsgo.RequestBuilder

	Domain   string   `json:"domain"`
	Identity string   `json:"identity"`
	TaskList TaskList `json:"taskList"`
}

type PollForActivityTaskResponse struct {
	ActivityId        string            `json:"activityId"`
	ActivityType      ActivityType      `json:"activityType"`
	Input             string            `json:"input"`
	StartedEventId    int64             `json:"startedEventId"`
	TaskToken         string            `json:"taskToken"`
	WorkflowExecution WorkflowExecution `json:"workflowExecution"`
}

func NewPollForActivityTaskRequest() *PollForActivityTaskRequest {
	req := new(PollForActivityTaskRequest)
	req.Host.Service = "swf"
	req.Host.Region = ""
	req.Host.Domain = "amazonaws.com"
	req.Key.AccessKeyId = ""
	req.Key.SecretAccessKey = ""
	req.Headers = make(map[string]string)
	req.Headers["X-Amz-Target"] = "SimpleWorkflowService.PollForActivityTask"
	req.RequestMethod = "POST"
	req.CanonicalUri = "/"
	return req
}

func (req *PollForActivityTaskRequest) VerifyInput() error {
	return nil
}

func (req PollForActivityTaskRequest) DeMarshalResponse(response []byte, headers map[string]string, statusCode int) interface{} {
	log.Println("response: ", string(response))
	resp := new(PollForActivityTaskResponse)
	if err := json.Unmarshal(response, resp); err != nil {
		return err
	}
	return resp
}

func (req PollForActivityTaskRequest) Request() (*PollForActivityTaskResponse, error) {
	request, err := awsgo.NewAwsRequest(&req, req)
	if err != nil {
		return nil, err
	}
	request.RequestSigningType = awsgo.RequestSigningType_AWS3
	resp, err := request.DoAndDemarshall(&req)
	if resp == nil {
		return nil, err
	}
	return resp.(*PollForActivityTaskResponse), err
}

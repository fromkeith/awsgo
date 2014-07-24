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
    "fmt"
    "encoding/json"
    "regexp"
)

const (
    ResourceAlreadyExistsExceptionType          = "ResourceAlreadyExistsException"
    InvalidSequenceTokenExceptionType           = "InvalidSequenceTokenException"
    DataAlreadyAcceptedExceptionType            = "DataAlreadyAcceptedException"
)

var (
    extractSequenceRegex    = regexp.MustCompile(`.*expected sequenceToken is: ([a-z0-9A-Z]+)$`)
    nextExtractSequenceRegex    = regexp.MustCompile(`.*with sequenceToken: ([a-z0-9A-Z]+)$`)
)


func ExtractNextTokenFromError(e *ErrorResult) string {
    if e.Type == InvalidSequenceTokenExceptionType {
        matches := extractSequenceRegex.FindStringSubmatch(e.Message)
        if matches == nil {
            return ""
        }
        return matches[1]
    }
    if e.Type == DataAlreadyAcceptedExceptionType {
        matches := nextExtractSequenceRegex.FindStringSubmatch(e.Message)
        if matches == nil {
            return ""
        }
        return matches[1]
    }
    return ""
}


// currently assuming it is the same format as dynamo...

type ErrorResult struct {
    Type        string  `json:"__type"`
    Message     string  `json:"message"`
    StatusCode  int
}

func (e * ErrorResult) Error() string {
    return fmt.Sprintf("%s : %s", e.Type, e.Message)
}

func CheckForErrorResponse(response []byte, statusCode int) error {
    errorResult := new(ErrorResult)
    err2 := json.Unmarshal([]byte(response), errorResult)
    if err2 == nil {
        if errorResult.Type != "" {
            errorResult.StatusCode = statusCode
            return errorResult
        }
    }
    if statusCode < 200 || statusCode > 299 {
        errorResult.Type = "UnknownError"
        errorResult.Message = string(response)
        errorResult.StatusCode = statusCode
        return errorResult
    }
    return nil
}
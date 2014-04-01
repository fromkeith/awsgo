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

package awsgo

import (
    "encoding/xml"
    "errors"
)

var (
    Verification_Error_AccessKeyEmpty = errors.New("Key.AccessKeyId cannot be empty")
    Verification_Error_SecretAccessKeyEmpty = errors.New("Key.SecretAccessKey cannot be empty")
)

type Error struct {
    Type    string
    Code    string
    Message string
}
type ErrorResponse struct {
    ErrorT   Error      `xml:"Error"`
    RequestId string
}

type UnmarhsallingError struct {
    ActualContent       string
    MarshallError       error
}

func (e * UnmarhsallingError) Error() string {
    return "Error unmarshalling response"
}

func (e * ErrorResponse) Error() string {
    return e.ErrorT.Code
}


// Looks at a raw response for an XML error response
func CheckForErrorXml(response []byte) error {
    errorResponse := new(ErrorResponse)
    xml.Unmarshal(response, errorResponse)
    if errorResponse.ErrorT.Message != "" {
        return errorResponse
    }
    return nil
}
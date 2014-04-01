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
    "bytes"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/x509"
    "fmt"
    "hash"
    "io/ioutil"
    "time"
)


// Meta data in a lot of aws requests
type ResponseMetaData struct {
    RequestId string
}



// The host of the service we are hitting.
// Urls are formed by taking Service.Region.Domain
type AwsHost struct {
    // Eg. dynamo
    Service string
    // Eg us-west-2
    Region string
    // Eg. amazonaws.com
    Domain string
    // If you want to hit your own custom test service.
    // Generally leave nil to use go's default cert chain.
    CustomCertificates []*x509.Certificate
}

// concat the parts together into a hostname.
func (h AwsHost) ToString() string {
    if h.Region == "" {
        return fmt.Sprintf("%s.%s",
            h.Service, h.Domain)
    }
    return fmt.Sprintf("%s.%s.%s",
        h.Service, h.Region, h.Domain)
}




// 20110909
func simpleDate(d time.Time) string {
    d = d.UTC()
    return fmt.Sprintf("%d%0.2d%0.2d",
        d.Year(), d.Month(), d.Day())
}
// 20130315T092054Z ISO 8601 basic format
func IsoDate(t time.Time) string {
  t = t.UTC()
  return fmt.Sprintf("%04d%02d%02dT%02d%02d%02dZ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second());
}

func createHMacHasher256() hash.Hash {
    return sha256.New()
}
func createHMacHasher1() hash.Hash {
    return sha1.New()
}

func BuildEmptyContentRequest(rb RequestBuilderInterface) (request AwsRequest, verifyError error) {
    return NewAwsRequest(rb, ioutil.NopCloser(bytes.NewBuffer([]byte(""))))
}

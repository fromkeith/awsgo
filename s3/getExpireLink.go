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

package s3

import (
    "github.com/fromkeith/awsgo"
    "time"
    "fmt"
    "net/url"
    "crypto/hmac"
    "crypto/sha1"
    "encoding/base64"
)

// Generates an expirable link to an s3 object
// Generated link will be: "//s3.amazonaws.com/{bucket}/{path}?{queryargs}"
func GenerateExpirableLink(bucket, path string, expires time.Time, creds awsgo.Credentials) string {
    var header string

    if creds.GetToken() != "" {
        header = fmt.Sprintf("x-amz-security-token:%s\n", creds.GetToken())
    }

    toSign := fmt.Sprintf(
        "GET\n\n\n%d\n%s/%s/%s",
        expires.Unix(),
        header,
        bucket,
        path,
    )
    hmacHasher := hmac.New(sha1.New, []byte(creds.SecretAccessKey))
    hmacHasher.Write([]byte(toSign))
    signature := base64.StdEncoding.EncodeToString(hmacHasher.Sum(nil))

    v := make(url.Values)
    v.Add("AWSAccessKeyId", creds.AccessKeyId)
    v.Add("Expires", fmt.Sprintf("%d", expires.Unix()))
    v.Add("Signature", signature)
    if creds.GetToken() != "" {
        v.Add("x-amz-security-token", creds.GetToken())
    }

    return fmt.Sprintf(
        "//s3.amazonaws.com/%s/%s?%s",
        bucket,
        path,
        v.Encode(),
    )
}


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
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "sync"
    "time"
)


// Aws credentials
type Credentials struct {
    AccessKeyId     string
    SecretAccessKey string
    token           string // used when you are using the IAM role. Otherwise should be empty
    expiration time.Time
}

var cachedCredentials Credentials
var credentialLock sync.Mutex

func determineSecurityRole() (string, error) {
    metaDataUrl := "http://169.254.169.254/latest/meta-data/iam/security-credentials"
    metaDataUri, _ := url.Parse(metaDataUrl)
    hreq := http.Request {
        URL: metaDataUri,
        Method: "GET",
        ProtoMajor: 1,
        ProtoMinor: 1,
        Close: true,
    }
    resp, err := http.DefaultClient.Do(&hreq)
    if err != nil {
        return "", errors.New("Failed to create HTTP Client")
    }
    defer resp.Body.Close()
    buf := bytes.NewBuffer(make([]byte, 0))
    io.Copy(buf, resp.Body)

    if resp.StatusCode != 200 {
        return "", errors.New(fmt.Sprintf("Got Status code: %d", resp.StatusCode))
    }
    return buf.String(), nil
}

type CredentialMetaData struct {
    Code string
    LastUpdated string
    Type string
    AccessKeyId string
    SecretAccessKey string
    Token string
    Expiration string
}

func credentialsAreLocal() bool {
    f, err := os.Open("awskeys.json")
    if err != nil {
        return false
    }
    defer f.Close()
    buf := bytes.NewBuffer(make([]byte, 0))
    io.Copy(buf, f)
    var tmp Credentials
    if err = json.Unmarshal([]byte(buf.String()), &tmp); err != nil {
        return false
    }
    if tmp.AccessKeyId == "" || tmp.SecretAccessKey == "" {
        return false
    }
    // expire in 10 hours
    tmp.expiration = time.Now().Add(time.Hour * 10)
    cachedCredentials = tmp
    return true
}

/** Returns security credentials either from a JSON file 'awskeys.json' or
 * from AWS Metadata service
 * @return credentials, error
 */
func GetSecurityKeys() (finalCred Credentials, err error)  {
    if cachedCredentials.AccessKeyId == "" || cachedCredentials.expiration.Unix() < time.Now().Unix() {
        credentialLock.Lock()
        defer credentialLock.Unlock()
        if cachedCredentials.AccessKeyId == "" || cachedCredentials.expiration.Unix() < time.Now().Unix()  {
            if !credentialsAreLocal() {
                var role string
                role, err = determineSecurityRole()
                if err != nil {
                    return
                }
                credentialUrl := fmt.Sprintf("http://169.254.169.254/latest/meta-data/iam/security-credentials/%s", role)
                credentialUri, _ := url.Parse(credentialUrl)
                hreq := http.Request {
                    URL: credentialUri,
                    Method: "GET",
                    ProtoMajor: 1,
                    ProtoMinor: 1,
                    Close: true,
                }
                resp, err2 := http.DefaultClient.Do(&hreq)
                if err2 != nil {
                    err = errors.New("Failed to create HTTP Client")
                    return
                }
                defer resp.Body.Close()

                buf := bytes.NewBuffer(make([]byte, 0))
                io.Copy(buf, resp.Body)

                if resp.StatusCode != 200 {
                    err = errors.New(fmt.Sprintf("Got Status code: %d", resp.StatusCode))
                    return
                }

                var credentials CredentialMetaData
                if err = json.Unmarshal([]byte(buf.String()), &credentials); err != nil {
                    return
                }

                if credentials.Code != "Success" {
                    err = errors.New("Failed to get security keys")
                    return
                }

                var tmp Credentials
                tmp.AccessKeyId = credentials.AccessKeyId
                tmp.SecretAccessKey = credentials.SecretAccessKey
                tmp.expiration, _ = time.Parse("2006-01-02T15:04:05Z", credentials.Expiration)
                tmp.token = credentials.Token
                cachedCredentials = tmp
            }
        }
    }
    finalCred = cachedCredentials
    err = nil
    return 
}

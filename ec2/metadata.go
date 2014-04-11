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

package ec2


import (
    "net/url"
    "net/http"
    "bytes"
    "encoding/json"
    "io"
    "errors"
    "fmt"
)


func MakeSimpleRequest(urlString string) (string, error) {
    metaDataUri, _ := url.Parse("http://169.254.169.254/latest" + urlString)
    hreq := http.Request {
        URL: metaDataUri,
        Method: "GET",
        ProtoMajor: 1,
        ProtoMinor: 1,
        Close: true,
    }
    resp, err := http.DefaultClient.Do(&hreq)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    buf := bytes.Buffer{}
    io.Copy(&buf, resp.Body)

    if resp.StatusCode != 200 {
        return "", errors.New(fmt.Sprintf("Got Status code: %d", resp.StatusCode))
    }
    return buf.String(), nil
}



// make the request to "http://169.254.169.254/latest/meta-data/iam/security-credentials"
func IamSecurityCredentials() (string, error) {
    return MakeSimpleRequest("/meta-data/iam/security-credentials")
}

// Get the AMI id
func AmiId() (string, error) {
    return MakeSimpleRequest("/meta-data/ami-id")
}

// Get the HostName
func HostName() (string, error) {
    return MakeSimpleRequest("/meta-data/hostname")
}

// Get the InstanceId
func InstanceId() (string, error) {
    return MakeSimpleRequest("/meta-data/instance-id")
}

// Get the AvailabilityZone
func AvailabilityZone() (string, error) {
    return MakeSimpleRequest("/meta-data/placement/availability-zone")
}

// Get the PublicHostName
func PublicHostName() (string, error) {
    return MakeSimpleRequest("/meta-data/public-hostname")
}


type InstanceIdentityT struct {
    Version             string      `json:"version"`
    InstanceId          string      `json:"instanceId"`
    BillingProducts     []string    `json:"billingProducts"`
    Architecture        string      `json:"architecture"`
    ImageId             string      `json:"imageId"`
    PendingTime         string      `json:"pendingTime"`
    InstanceType        string      `json:"instanceType"`
    KernelId            string      `json:"kernelId"`
    AccountId           string      `json:"accountId"`
    RamDiskId           string      `json:"ramdiskId"`
    Region              string      `json:"region"`
    AvailabilityZone    string      `json:"availabilityZone"`
    DevPayProductCodes  []string    `json:"devpayProductCodes"`
    PrivateIp           string      `json:"privateIp"`
}

// Get the instance identity
func InstanceIdentity() (InstanceIdentityT, error) {
    resp, err := MakeSimpleRequest("/dynamic/instance-identity/document")
    if err != nil {
        return InstanceIdentityT{}, err
    }
    var it InstanceIdentityT
    err = json.Unmarshal([]byte(resp), &it)
    if err != nil {
        return InstanceIdentityT{}, err
    }
    return it, nil
}


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
    "crypto/hmac"
    "crypto/md5"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/pmylund/sortutil"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"
    "crypto/tls"
)

const (
    RequestSigningType_AWS4 = 1
    RequestSigningType_REST = 2
    RequestSigningType_AWS2 = 3
)




// Base of a request. Used across all requests.
type RequestBuilder struct {
    // The Host we are hitting
    Host AwsHost                    `json:"-"`
    // The Credentials to use
    Key  Credentials                `json:"-"`
    // Any custom headers
    Headers map[string]string       `json:"-"`
    // The method we are using GET, PUT, POST, ...
    RequestMethod string            `json:"-"`
    // The uri we are hitting.
    CanonicalUri string             `json:"-"`
}

// Implemented by each AWS request
// Provides some standard steps to doing a request, and handling its response.
type RequestBuilderInterface interface {
    // verify the request before we send it
    VerifyInput() error
    // Get the underlying RequestBuilder in the struct
    GetRequestBuilder() RequestBuilder
    // Unmarshal the response
    DeMarshalResponse(response []byte, headers map[string]string, statusCode int) (interface{})
}

type AwsRequest struct {
    Host AwsHost
    Date time.Time
    Headers map[string]string
    Payload string
    PayloadReader io.ReadCloser
    Key Credentials
    RequestMethod string
    CanonicalUri string
    RequestSigningType int
    // generated
    signature string
    scope string
    payloadHash []byte
}

func (r RequestBuilder) GetRequestBuilder() RequestBuilder {
    return r
}

// verify the RequestBuilder base has what is required.
func verifyInput(r RequestBuilder) (error) {
    if len(r.Host.Domain) == 0 {
        r.Host.Domain = "amazonaws.com"
    }
    if len(r.Key.AccessKeyId) == 0 {
        return Verification_Error_AccessKeyEmpty
    }
    if len(r.Key.SecretAccessKey) == 0 {
        return Verification_Error_SecretAccessKeyEmpty
    }
    return nil
}

// Start creating the request.
// Take the base information, and the data we are going to transport.
func createAwsRequest(rb RequestBuilder, marsh interface{}) (request AwsRequest) {
    request.Host = rb.Host
    request.Key = rb.Key
    request.Headers = rb.Headers
    request.RequestMethod = rb.RequestMethod
    request.CanonicalUri = rb.CanonicalUri
    request.Date = time.Now()
    if r, ok := marsh.(io.ReadCloser); ok {
        request.PayloadReader = r
    } else if marsh != nil {
        pay, _ := json.Marshal(marsh)
        request.Payload = string(pay)
        if request.Headers == nil {
            request.Headers = make(map[string]string)
        }
        request.Headers["Content-Type"] = "application/x-amz-json-1.0"
        request.Headers["Content-Length"] = fmt.Sprintf("%d", len(request.Payload))
    }
    return
}

// Given the RequestBuilderInterface this will verify the underlying request
// and then create a new AwsRequest instance.
// returned Request is only valid if error is not nill
func NewAwsRequest(rb RequestBuilderInterface, marsh interface{}) (request AwsRequest, verifyError error) {
    verifyError = rb.VerifyInput()
    if verifyError != nil {
        return
    }
    verifyError = verifyInput(rb.GetRequestBuilder())
    if verifyError != nil {
        return
    }
    request = createAwsRequest(rb.GetRequestBuilder(), marsh)
    return
}



// Perform the actual request, calling the demarshall on the RequestBuilderInterface
// returns the result of the Demarshall, or other errors.
func (request AwsRequest) DoAndDemarshall(rb RequestBuilderInterface) (interface{}, error) {
    responseIo, responseHeaders, statusCode, err := request.Do()
    if err != nil {
        return nil, err
    }
    defer responseIo.Close()

    var responseContent []byte
    buf := bytes.Buffer{}
    if _, err = io.Copy(&buf, responseIo); err != nil {
        return nil, err
    }
    responseContent = []byte(buf.String())

    //fmt.Println("responseContent")

    val := rb.DeMarshalResponse(responseContent, responseHeaders, statusCode)
    if t, ok := val.(error); ok {
        return nil, t
    }
    return val, nil
}

// Performs the actual request
// Returns:
//      io.ReaderCloser - the unclosed response of the request
//      map[string]string - the headers in the response
//      int - that status code the response
//      error - any errors that occured
func (req AwsRequest) Do() (io.ReadCloser, map[string]string, int, error) {

    // add the required headers
    req.Headers["Host"] = strings.ToLower(req.Host.ToString())
    req.Headers["user-agent"] = "go-aws-client-0.1"
    req.Headers["x-amz-date"] = IsoDate(req.Date)
    if req.Key.token != "" {
        req.Headers["x-amz-security-token"] = req.Key.token
    }

    if err := signRequest(&req); err != nil {
        return nil, nil, 0, err
    }

    req.Headers["Connection"] = "Keep-Alive"

    // create headers for the actual request
    reqHeaders := http.Header{}
    for k, v := range req.Headers {
        reqHeaders.Add(k, v)
    }

    url_, err := getUrl(req)
    if err != nil {
        return nil, nil, 0, err
    }

    // the base request
    hreq := http.Request {
        URL: url_,
        Method: req.RequestMethod,
        ProtoMajor: 1,
        ProtoMinor: 1,
        Close: true, // test what we want this. I seem to remember needing close...
        Header: reqHeaders,
    }

    httpClient := addCustomCertsAndCreateClient(req)

    if val, ok := req.Headers["Content-Length"]; ok {
        hreq.ContentLength, _ = strconv.ParseInt(val, 10, 64)
    }
    if req.Payload != "" {
        //fmt.Println("Payload", req.Payload)
        hreq.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(req.Payload)))
    } else if req.PayloadReader != nil {
        hreq.Body = req.PayloadReader
    }
    resp, err := httpClient.Do(&hreq)
    if err != nil {
        return nil, nil, 0, err
    }

    responseHeaders := make(map[string]string)
    for k, v := range resp.Header {
        responseHeaders[strings.ToLower(k)] = strings.Join(v, ";")
    }
    return resp.Body, responseHeaders, resp.StatusCode, nil
}

func addCustomCertsAndCreateClient(req AwsRequest) (http.Client) {
    // add in any custom certs they want us to use
    var rootCA *x509.CertPool
    if len(req.Host.CustomCertificates) > 0 {
        rootCA = x509.NewCertPool()
        for i := range req.Host.CustomCertificates {
            rootCA.AddCert(req.Host.CustomCertificates[i])
        }
    }
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{
            RootCAs : rootCA,
        },
    }
    return http.Client{Transport: tr}
}

func getUrl(req AwsRequest) (*url.URL, error) {
    // create the url string. Not all aws request have regions.
    var urlStr string
    if req.Host.Region == "" {
        urlStr = fmt.Sprintf("https://%s.%s%s", req.Host.Service, req.Host.Domain, req.CanonicalUri)
    } else {
        urlStr = fmt.Sprintf("https://%s.%s.%s%s", req.Host.Service, req.Host.Region, req.Host.Domain, req.CanonicalUri)
    }
    return url.Parse(urlStr)
}

func signRequest(req *AwsRequest) error {
    if req.RequestSigningType == RequestSigningType_AWS4 {
        req.createSignature()
    } else if req.RequestSigningType == RequestSigningType_REST {
        req.createRestSignature()
    } else if req.RequestSigningType == RequestSigningType_AWS2 {
        req.createV2Signature()
    } else {
        return errors.New("Invalid request signing type")
    }
    return nil
}

// http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html
func (req * AwsRequest) createRestSignature() {
    payloadHash := ""
    md5Hasher := md5.New()
    if req.Payload != "" {
        md5Hasher.Write([]byte(req.Payload)) // TODO: check return code?
        payloadHash = string(md5Hasher.Sum(nil))
    }
    canonicalHeaders, _  := req.createCanonicalHeaders("x-amz-")
    canonicalResource := fmt.Sprintf("%s", req.CanonicalUri)

    stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s\n%s%s",
        req.RequestMethod, payloadHash, req.Headers["Content-Type"], "" /*old school date*/,
        canonicalHeaders, canonicalResource)

    hmacHasher := hmac.New(createHMacHasher1, []byte(req.Key.SecretAccessKey))
    hmacHasher.Write([]byte(stringToSign))

    signature := base64.StdEncoding.EncodeToString(hmacHasher.Sum(nil))
    authorization := fmt.Sprintf("AWS %s:%s", req.Key.AccessKeyId, signature)
    req.Headers["Authorization"] = authorization
}

func (req * AwsRequest) createCanonicalHeaders(prefixReq string) (canonicalHeaders string, signedHeaders string) {
    mapKeys := make([]string, len(req.Headers))
    i := 0
    for k := range req.Headers {
        mapKeys[i] = k
        i++
    }
    sortutil.CiAsc(mapKeys)
    canonicalHeaders = ""
    signedHeaders = ""
    for i := range mapKeys {
        if prefixReq != "" {
            if !strings.HasPrefix(strings.ToLower(mapKeys[i]), prefixReq) {
                continue
            }
        }
        if len(signedHeaders) > 0 {
            signedHeaders = fmt.Sprintf("%s;%s", signedHeaders, strings.ToLower(mapKeys[i]))
            canonicalHeaders = fmt.Sprintf("%s%s:%s\n", canonicalHeaders, strings.ToLower(mapKeys[i]), req.Headers[mapKeys[i]])
        } else {
            signedHeaders = fmt.Sprintf("%s", strings.ToLower(mapKeys[i]))
            canonicalHeaders = fmt.Sprintf("%s:%s\n", strings.ToLower(mapKeys[i]), req.Headers[mapKeys[i]])
        }
    }
    return
}

//http://docs.aws.amazon.com/general/latest/gr/signature-version-2.html
func  (req * AwsRequest) createV2Signature() {

    canonicalQueryString := ""
    fixedUrl := req.CanonicalUri
    if !strings.Contains(req.CanonicalUri, "?") {
        req.CanonicalUri = req.CanonicalUri + "?"
    } else {
        req.CanonicalUri = req.CanonicalUri + "&"
    }
    now := time.Now()

    req.CanonicalUri = fmt.Sprintf("%sAWSAccessKeyId=%s&SignatureMethod=HmacSHA256&SignatureVersion=2&Timestamp=%s",
        req.CanonicalUri,
        strings.Replace(url.QueryEscape(req.Key.AccessKeyId), "+", "%20", -1),
        strings.Replace(url.QueryEscape(now.Format(time.RFC3339)), "+", "%20", -1),
        )
    if req.Key.token != "" {
        req.CanonicalUri = fmt.Sprintf("%s&SecurityToken=%s",
            strings.Replace(url.QueryEscape(req.Key.token), "+", "%20", -1),)
    }

    urlSplit := strings.Split(req.CanonicalUri, "?")
    fixedUrl = urlSplit[0]
    sp := strings.Split(urlSplit[1], "&")
    sortutil.Asc(sp)
    for i := range sp {
        if len(canonicalQueryString) > 0 {
            unEscaped := strings.Replace(sp[i], "+", "%20", -1)
            canonicalQueryString = fmt.Sprintf("%s&%s", canonicalQueryString, unEscaped)
        } else {
            canonicalQueryString = sp[i]
        }
    }

    toSign := fmt.Sprintf("%s\n%s\n%s\n%s",
        req.RequestMethod,
        req.Host.ToString(),
        fixedUrl,
        canonicalQueryString)

    hmacHasher := hmac.New(createHMacHasher256, []byte(req.Key.SecretAccessKey))
    hmacHasher.Write([]byte(toSign))
    sig := hmacHasher.Sum(nil)

    req.CanonicalUri = fmt.Sprintf("%s&Signature=%s", req.CanonicalUri,
        strings.Replace(
            url.QueryEscape(
                base64.StdEncoding.EncodeToString(sig),
            ), "+", "%20", -1))
    fmt.Println("ToSign: ", toSign)
    fmt.Println("Request: ", req.CanonicalUri)
}


// http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
func (req * AwsRequest) createSignature() {
    hasher := sha256.New()
    hasher.Write([]byte(req.Payload)) // TODO: check return code?
    req.payloadHash = hasher.Sum(nil)

    // Stupid Amazon. Why have the service name and the url not match???
    fixedService := req.Host.Service
    if req.Host.Service == "email" {
        fixedService = "ses"
    }

    req.scope = fmt.Sprintf("%s/%s/%s/aws4_request", simpleDate(req.Date), req.Host.Region, fixedService)

    canonicalQueryString := ""
    fixedUrl := req.CanonicalUri
    if strings.Contains(req.CanonicalUri, "?") {
        urlSplit := strings.Split(req.CanonicalUri, "?")
        fixedUrl = urlSplit[0]
        sp := strings.Split(urlSplit[1], "&")
        sortutil.CiAsc(sp)
        for i := range sp {
            if len(canonicalQueryString) > 0 {
                unEscaped := strings.Replace(sp[i], "+", "%20", -1)
                canonicalQueryString = fmt.Sprintf("%s&%s", canonicalQueryString, unEscaped)
            } else {
                canonicalQueryString = sp[i]
            }
        }
    }

    req.Headers["x-amz-content-sha256"] = fmt.Sprintf("%x", req.payloadHash)

    canonicalHeaders, signedHeaders := req.createCanonicalHeaders("")

    canonicalReq := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%x",
        req.RequestMethod, fixedUrl, canonicalQueryString, canonicalHeaders, signedHeaders, req.payloadHash)
    hasher.Reset()
    hasher.Write([]byte(canonicalReq)) // TODO: check return code?
    hashedCanonicalReq := hasher.Sum(nil)


    stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%x",
        IsoDate(req.Date), req.scope, hashedCanonicalReq)

    //fmt.Println("Canon", canonicalReq)
    //fmt.Println("String To Sign:", stringToSign)

    hasher.Reset()
    hmacHasher := hmac.New(createHMacHasher256, []byte(fmt.Sprintf("AWS4%s", req.Key.SecretAccessKey)))
    hmacHasher.Write([]byte(simpleDate(req.Date)))
    hmacDate := hmacHasher.Sum(nil)

    hmacHasher = hmac.New(createHMacHasher256, hmacDate)
    hmacHasher.Write([]byte(req.Host.Region))
    hmacRegion := hmacHasher.Sum(nil)

    hmacHasher = hmac.New(createHMacHasher256, hmacRegion)
    hmacHasher.Write([]byte(fixedService))
    hmacService := hmacHasher.Sum(nil)

    hmacHasher = hmac.New(createHMacHasher256, hmacService)
    hmacHasher.Write([]byte("aws4_request"))
    signingKey := hmacHasher.Sum(nil)

    hmacHasher = hmac.New(createHMacHasher256, signingKey)
    hmacHasher.Write([]byte(stringToSign))
    req.signature = fmt.Sprintf("%x", hmacHasher.Sum(nil))

    req.Headers["Authorization"] =
        fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s/%s/%s/aws4_request, SignedHeaders=%s, Signature=%s",
            req.Key.AccessKeyId, simpleDate(req.Date), req.Host.Region, fixedService, signedHeaders, req.signature)
}


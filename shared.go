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
    "crypto/sha1"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "encoding/xml"
    "errors"
    "fmt"
    "github.com/pmylund/sortutil"
    "hash"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
    "time"
    "sync"
)

const (
    RequestSigningType_AWS4 = 1
    RequestSigningType_REST = 2
)

type AwsStringItem struct {
    Value string      `json:"S,omitempty"`
    Values []string   `json:"SS,omitempty"`
}
type AwsNumberItem struct {
    Value float64    `json:"N,string"`
    Values []float64 `json:"NN,string"`
}

type ResponseMetaData struct {
    RequestId string
}

func NewStringItem(items ... string) AwsStringItem {
    var s AwsStringItem
    if (len(items) == 1) {
        s.Value = items[0]
    } else {
        s.Values = make([]string, len(items))
        for i, val := range items {
            s.Values[i] = val
        }
    }
    return s
}

func NewNumberItem(items ... float64) AwsNumberItem {
    var s AwsNumberItem
    if (len(items) == 1) {
        s.Value = items[0]
    } else {
        s.Values = make([]float64, len(items))
        for i, val := range items {
            s.Values[i] = val
        }
    }
    return s
}


/** Converts from an unknown interface... like:
 *      string, []string, float, []float64
 *  into the expected awsgo.AwsStringItem or awsgo.AwsNumberItem
 */
func ConvertToAwsItem(unknown interface{}) interface{} {
    switch j := unknown.(type) {
        case string:
            return NewStringItem(j)
        case float64:
            return NewNumberItem(j)
        case int:
            return NewNumberItem(float64(j))
        case float32:
            return NewNumberItem(float64(j))
        case int64:
            return NewNumberItem(float64(j))
        case []string:
            return AwsStringItem{"", j}
        case []int:
            // we need to cast these over
            vals64 := make([]float64, len(j))
            for i := range j {
                vals64[i] = float64(j[i])
            }
            return AwsNumberItem{0, vals64}
        case []int64:
            // we need to cast these over
            vals64 := make([]float64, len(j))
            for i := range j {
                vals64[i] = float64(j[i])
            }
            return AwsNumberItem{0, vals64}
        case []float32:
            // we need to cast these over
            vals64 := make([]float64, len(j))
            for i := range j {
                vals64[i] = float64(j[i])
            }
            return AwsNumberItem{0, vals64}
        case []float64:
            return AwsNumberItem{0, j}
        case AwsNumberItem:
            return j
        case AwsStringItem:
            return j
        default:
            panic(fmt.Sprintf("Unknown data type: %v %T", j, j))
            return j
    }
    return unknown
}

func FromRawMapToEasyTypedMap(raw map[string]map[string]interface{}, item map[string]interface{}) {
    for key, value := range raw {
        if v, ok := value["S"]; ok {
            switch t := v.(type) {
            case string:
                item[key] = t
                break
            default:
                panic("Item map was type 'S' but did not have string content!")
            }
        }
        if v, ok := value["SS"]; ok {
            if t, ok := v.([]interface{}); ok {
                vals := make([]string, len(t))
                for i := range t {
                    if t2, ok := t[i].(string); ok {
                        vals[i] = t2
                    } else {
                        panic(fmt.Sprintf("Expected string in SS but got: %T", t[i]))
                    }
                }
                item[key] = vals
            } else {
                panic(fmt.Sprintf("Item map was type 'NS' but did not have []string content! (We expect it as string, but convert to []float). Got %T", t))
            }
        }
        if v, ok := value["N"]; ok {
            switch t := v.(type) {
            case string:
                f, _ := strconv.ParseFloat(t, 64)
                item[key] = f
                break
            default:
                panic(fmt.Sprintf("Item map was type 'N' but did not have string content! (We expect it as string, but convert to float). Got %T", t))
            }
        }
        if v, ok := value["NS"]; ok {
            if t, ok := v.([]interface{}); ok {
                nums := make([]float64, len(t))
                for i := range t {
                    if t2, ok := t[i].(string); ok {
                        nums[i], _ = strconv.ParseFloat(t2, 64)
                    } else {
                        panic(fmt.Sprintf("Expected string in NS but got: %T", t[i]))
                    }
                }
                item[key] = nums
            } else {
                panic(fmt.Sprintf("Item map was type 'NS' but did not have []string content! (We expect it as string, but convert to []float). Got %T", t))
            }
        }
    }
}

func CheckForErrorXml(response []byte) error {
    errorResponse := new(ErrorResponse)
    xml.Unmarshal(response, errorResponse)
    if errorResponse.ErrorT.Message != "" {
        return errorResponse
    }
    return nil
}


type Error struct {
    Type    string
    Code    string
    Message string
}
type ErrorResponse struct {
    ErrorT   Error      `xml:"Error"`
    RequestId string
}

func (e * ErrorResponse) Error() string {
    return e.ErrorT.Code
}


type AwsHost struct {
    Service string
    Region string
    Domain string
}

type RequestBuilder struct {
    Host AwsHost                `json:"-"`
    Key  Credentials             `json:"-"`
    Headers map[string]string `json:"-"`
    RequestMethod string               `json:"-"`
    CanonicalUri string         `json:"-"`
}

type RequestBuilderInterface interface {
    VerifyInput() error
    CreateJsonAwsRequest(marsh interface{}) AwsRequest
    CreateReaderAwsRequest(r io.ReadCloser) AwsRequest
    DeMarshalGetItemResponse([]byte, map[string]string) (interface{})
}

func (r * RequestBuilder) VerifyInput() (error) {
    if len(r.Host.Domain) == 0 {
        return errors.New("Host.Domain cannot be empty")
    }
    if len(r.Key.AccessKeyId) == 0 {
        return errors.New("Key.AccessKeyId cannot be empty")
    }
    if len(r.Key.SecretAccessKey) == 0 {
        return errors.New("Key.SecretAccessKey cannot be empty")
    }
    return nil
}

func (rb RequestBuilder) CreateJsonAwsRequest(marsh interface{}) (request AwsRequest) {
    request.Host = rb.Host
    request.Key = rb.Key
    request.Headers = rb.Headers
    request.RequestMethod = rb.RequestMethod
    request.CanonicalUri = rb.CanonicalUri
    request.Date = time.Now()
    pay, _ := json.Marshal(marsh)
    request.Payload = string(pay)
    if request.Headers == nil {
        request.Headers = make(map[string]string)
    }
    request.Headers["Content-Type"] = "application/x-amz-json-1.0"
    request.Headers["Content-Length"] = fmt.Sprintf("%d", len(request.Payload))
    return
}
func (rb RequestBuilder) CreateReaderAwsRequest(r io.ReadCloser) (request AwsRequest) {
    request.Host = rb.Host
    request.Key = rb.Key
    request.Headers = rb.Headers
    request.RequestMethod = rb.RequestMethod
    request.CanonicalUri = rb.CanonicalUri
    request.Date = time.Now()
    request.PayloadReader = r
    return
}

func DoRequest(rb RequestBuilderInterface, request AwsRequest) (interface{}, error) {
    response, responseHeaders, _, err := request.SendRequest()
    if err != nil {
        return nil, err
    }
    val := rb.DeMarshalGetItemResponse([]byte(response), responseHeaders)
    switch t := val.(type) {
    case error:
        return nil, t
    default:
        return t, nil
    }
}

func BuildRequest(rb RequestBuilderInterface, marsh interface{}) (request AwsRequest, verifyError error) {
    verifyError = rb.VerifyInput()
    if verifyError != nil {
        return request, verifyError
    }
    request = rb.CreateJsonAwsRequest(marsh)
    return request, nil
}

func BuildEmptyContentRequest(rb RequestBuilderInterface) (request AwsRequest, verifyError error) {
    return BuildReaderRequest(rb, ioutil.NopCloser(bytes.NewBuffer([]byte(""))))
}

func BuildReaderRequest(rb RequestBuilderInterface, r io.ReadCloser) (request AwsRequest, verifyError error) {
    verifyError = rb.VerifyInput()
    if verifyError != nil {
        return request, verifyError
    }
    request = rb.CreateReaderAwsRequest(r)
    return request, nil
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

func (h AwsHost) ToString() string {
    if h.Region == "" {
        return fmt.Sprintf("%s.%s",
            h.Service, h.Domain)
    }
    return fmt.Sprintf("%s.%s.%s",
        h.Service, h.Region, h.Domain)
}

// 20110909
func SimpleDate(d time.Time) string {
    d = d.UTC()
    return fmt.Sprintf("%d%0.2d%0.2d",
        d.Year(), d.Month(), d.Day())
}
// 20130315T092054Z ISO 8601 basic format
func IsoDate(t time.Time) string {
  t = t.UTC()
  return fmt.Sprintf("%04d%02d%02dT%02d%02d%02dZ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second());
}

func CreateHMacHasher256() hash.Hash {
    return sha256.New()
}
func CreateHMacHasher1() hash.Hash {
    return sha1.New()
}

func (req * AwsRequest) SendRequest() (string, map[string]string, int, error) {
    req.Headers["Host"] = strings.ToLower(req.Host.ToString())
    req.Headers["user-agent"] = "go-aws-client-0.1"
    req.Headers["x-amz-date"] = IsoDate(req.Date)
    if req.Key.Token != "" {
        req.Headers["x-amz-security-token"] = req.Key.Token
    }

    if req.RequestSigningType == RequestSigningType_AWS4 {
        req.CreateSignature(true)
    } else if req.RequestSigningType == RequestSigningType_REST {
        req.CreateRestSignature()
    } else {
        return "", nil, 0, errors.New("Invalid request signing type")
    }
    
    req.Headers["Connection"] = "Keep-Alive"

    reqHeaders := http.Header{}
    for k, v := range req.Headers {
        reqHeaders.Add(k, v)
    }
    urlStr := ""
    if req.Host.Region == "" {
        urlStr = fmt.Sprintf("https://%s.%s%s", req.Host.Service, req.Host.Domain, req.CanonicalUri)
    } else {
        urlStr = fmt.Sprintf("https://%s.%s.%s%s", req.Host.Service, req.Host.Region, req.Host.Domain, req.CanonicalUri)
    }
    //fmt.Println(urlStr)

    url, err := url.Parse(urlStr)
    if err != nil {
        return "", nil, 0, err
    }
    hreq := http.Request {
        URL: url,
        Method: req.RequestMethod,
        ProtoMajor: 1,
        ProtoMinor: 1,
        Close: true,
        Header: reqHeaders,
    }
    if val, ok := req.Headers["Content-Length"]; ok {
        hreq.ContentLength, _ = strconv.ParseInt(val, 10, 64)
    }
    if req.Payload != "" {
        //fmt.Println("Payload", req.Payload)
        hreq.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(req.Payload)))
    } else if req.PayloadReader != nil {
        hreq.Body = req.PayloadReader
    }
    resp, err := http.DefaultClient.Do(&hreq)
    if err != nil {
        return "", nil, 0, err
    }
    defer resp.Body.Close()

    var responseContent []byte

    buf := bytes.NewBuffer(make([]byte, 0))
    if _, err = io.Copy(buf, resp.Body); err != nil {
        fmt.Println("Failed to copy part of response", err);
    }
    responseContent = []byte(buf.String())

    responseHeaders := make(map[string]string)
    for k, v := range resp.Header {
        responseHeaders[strings.ToLower(k)] = strings.Join(v, ";")
    }
    //fmt.Println("Response", string(responseContent))
    return string(responseContent), responseHeaders, resp.StatusCode, nil
}

// http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html
func (req * AwsRequest) CreateRestSignature() {
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

    hmacHasher := hmac.New(CreateHMacHasher1, []byte(req.Key.SecretAccessKey))
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


// http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
func (req * AwsRequest) CreateSignature(calcEmptyHash bool) {
    hasher := sha256.New()
    if req.Payload != "" || calcEmptyHash {
        hasher.Write([]byte(req.Payload)) // TODO: check return code?
        req.payloadHash = hasher.Sum(nil)
    }

    // Stupid Amazon. Why have the service name and the url not match???
    fixedService := req.Host.Service
    if req.Host.Service == "email" {
        fixedService = "ses"
    }

    req.scope = fmt.Sprintf("%s/%s/%s/aws4_request", SimpleDate(req.Date), req.Host.Region, fixedService)

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
    if req.Payload != "" || calcEmptyHash {
        req.Headers["x-amz-content-sha256"] = fmt.Sprintf("%x", req.payloadHash)
    }
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
    hmacHasher := hmac.New(CreateHMacHasher256, []byte(fmt.Sprintf("AWS4%s", req.Key.SecretAccessKey)))
    hmacHasher.Write([]byte(SimpleDate(req.Date)))
    hmacDate := hmacHasher.Sum(nil)

    hmacHasher = hmac.New(CreateHMacHasher256, hmacDate)
    hmacHasher.Write([]byte(req.Host.Region))
    hmacRegion := hmacHasher.Sum(nil)

    hmacHasher = hmac.New(CreateHMacHasher256, hmacRegion)
    hmacHasher.Write([]byte(fixedService))
    hmacService := hmacHasher.Sum(nil)

    hmacHasher = hmac.New(CreateHMacHasher256, hmacService)
    hmacHasher.Write([]byte("aws4_request"))
    signingKey := hmacHasher.Sum(nil)

    hmacHasher = hmac.New(CreateHMacHasher256, signingKey)
    hmacHasher.Write([]byte(stringToSign))
    req.signature = fmt.Sprintf("%x", hmacHasher.Sum(nil))

    req.Headers["Authorization"] =
        fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s/%s/%s/aws4_request, SignedHeaders=%s, Signature=%s",
            req.Key.AccessKeyId, SimpleDate(req.Date), req.Host.Region, fixedService, signedHeaders, req.signature)
}



type Credentials struct {
    AccessKeyId     string
    SecretAccessKey string
    Token           string // used when you are using the IAM role. Otherwise should be empty
    Expiration time.Time
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
    tmp.Expiration = time.Now().Add(time.Hour * 10)
    cachedCredentials = tmp
    return true
}

/** Returns security credentials either from a JSON file 'awskeys.json' or
 * from AWS Metadata service
 * @return credentials, error
 */
func GetSecurityKeys() (finalCred Credentials, err error)  {
    if cachedCredentials.AccessKeyId == "" || cachedCredentials.Expiration.Unix() < time.Now().Unix() {
        credentialLock.Lock()
        defer credentialLock.Unlock()
        if cachedCredentials.AccessKeyId == "" || cachedCredentials.Expiration.Unix() < time.Now().Unix()  {
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
                tmp.Expiration, _ = time.Parse("2006-01-02T15:04:05Z", credentials.Expiration)
                tmp.Token = credentials.Token
                cachedCredentials = tmp
            }
        }
    }
    finalCred = cachedCredentials
    err = nil
    return 
}


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
package dynamo

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "fmt"
    "io/ioutil"
    "strings"
    "crypto/x509"
    "github.com/fromkeith/awsgo"
)

func easyFloatCompare(a, b float64) bool {
    return a - 0.00001 < b && a + 0.00001 > b
}


func doGetItemTest(req *GetItemRequest, handler http.HandlerFunc) (*GetItemResponse, error) {
    ts := httptest.NewTLSServer(handler)
    defer ts.Close()
    certAsx509, _ := x509.ParseCertificate(ts.TLS.Certificates[0].Certificate[0])

    req.Host.Domain = strings.TrimPrefix(ts.URL, "https://127.0.")
    req.Host.Region = "0"
    req.Host.Service = "127"
    req.Key.AccessKeyId = "akey"
    req.Key.SecretAccessKey = "skey"
    req.Host.CustomCertificates = []*x509.Certificate{certAsx509}

    resp, err := req.Request()
    return resp, err
}

func verifySimpleRequestBody(t * testing.T, r * http.Request) {
    expectedBody := `
        {"ConsistentRead":"false","Key":{"blah":{"S":"asdf"}},"TableName":"asd","ReturnConsumedCapacity":"NONE"}
    `
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        t.Fatal("couldn't read content! error: %v", err)
    }
    if string(body) != strings.TrimSpace(expectedBody) {
        t.Fatal("Expected request body: '%s'. Got: '%s'", strings.TrimSpace(expectedBody), string(body))
    }
}


func Test_WorkingGetSingleItem_String(t * testing.T) {
    expectedValue := "Hello"
    expectedKey := "one"

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        verifySimpleRequestBody(t, r)
        fmt.Fprintf(w, `
            {
                "Item" : {
                    "%s" : { "S" : "%s"}
                }
            }
            `, expectedKey, expectedValue)
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"
    itemReq.Search["blah"] = "asdf"

    resp, err := doGetItemTest(itemReq, handler)
    if err != nil {
        t.Fatalf("Error should be nil. Got: %v", err)
    }
    if _, ok := resp.Item[expectedKey]; !ok {
        t.Fatalf("Response should have item named %s !", expectedKey)
    }
    if val, ok := resp.Item[expectedKey].(string); !ok {
        t.Fatalf("Response item '%s' should be a string. Got: '%T'", expectedKey, resp.Item[expectedKey])
    } else if val != expectedValue {
        t.Fatalf("Response item '%s' should be '%s'. Got: '%s'", expectedKey, expectedValue, val)
    }
}

func Test_WorkingGetSingleItem_StringArray(t * testing.T) {
    expectedValue := []string{"Hello", "There", "How", "Are", "You", "?"}
    expectedKey := "multi"

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        verifySimpleRequestBody(t, r)
        fmt.Fprintf(w, `
            {
                "Item" : {
                    "%s" : { "SS" : ["%s", "%s", "%s", "%s", "%s", "%s"] }
                }
            }
            `,
            expectedKey,
            expectedValue[0],
            expectedValue[1],
            expectedValue[2],
            expectedValue[3],
            expectedValue[4],
            expectedValue[5],
            )
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"
    itemReq.Search["blah"] = "asdf"

    resp, err := doGetItemTest(itemReq, handler)
    if err != nil {
        t.Fatalf("Error should be nil. Got: %v", err)
    }
    if _, ok := resp.Item[expectedKey]; !ok {
        t.Fatalf("Response should have item named %s !", expectedKey)
    }
    if val, ok := resp.Item[expectedKey].([]string); !ok {
        t.Fatalf("Response item '%s' should be a []string. Got: '%T'", expectedKey, resp.Item[expectedKey])
    } else {
        if len(val) != len(expectedValue) {
            t.Fatalf("Response item '%s' should be len('%d'). Got: len('%d')", expectedKey, len(expectedValue), len(val))
        }
        expectedValueMap := make(map[string]string)
        for i := range expectedValue {
            expectedValueMap[expectedValue[i]] = expectedValue[i]
        }
        for i := range val {
            if _, ok := expectedValueMap[val[i]]; !ok {
                t.Fatalf("Got unexpected value '%s' not in set %v", val[i], expectedValue)
            }
        }
    }
}

func Test_WorkingGetSingleItem_Number(t * testing.T) {
    expectedValue := 123.23
    expectedKey := "one"

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        verifySimpleRequestBody(t, r)
        fmt.Fprintf(w, `
            {
                "Item" : {
                    "%s" : { "N" : "%f"}
                }
            }
            `, expectedKey, expectedValue)
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"
    itemReq.Search["blah"] = "asdf"

    resp, err := doGetItemTest(itemReq, handler)
    if err != nil {
        t.Fatalf("Error should be nil. Got: %v", err)
    }
    if _, ok := resp.Item[expectedKey]; !ok {
        t.Fatalf("Response should have item named %s !", expectedKey)
    }
    if val, ok := resp.Item[expectedKey].(float64); !ok {
        t.Fatalf("Response item '%s' should be a float64. Got: '%T'", expectedKey, resp.Item[expectedKey])
    } else if val != expectedValue {
        t.Fatalf("Response item '%s' should be '%f'. Got: '%f'", expectedKey, expectedValue, val)
    }
}

func Test_WorkingGetSingleItem_NumberArray(t * testing.T) {
    expectedValue := []float64{32.23, 45.3433, 3221.3, 32, 555.44, -0.1232}
    expectedKey := "multi"

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        verifySimpleRequestBody(t, r)

        fmt.Fprintf(w, `
            {
                "Item" : {
                    "%s" : { "NS" : ["%f", "%f", "%f", "%f", "%f", "%f"] }
                }
            }
            `,
            expectedKey,
            expectedValue[0],
            expectedValue[1],
            expectedValue[2],
            expectedValue[3],
            expectedValue[4],
            expectedValue[5],
            )
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"
    itemReq.Search["blah"] = "asdf"

    resp, err := doGetItemTest(itemReq, handler)
    if err != nil {
        t.Fatalf("Error should be nil. Got: %v", err)
    }
    if _, ok := resp.Item[expectedKey]; !ok {
        t.Fatalf("Response should have item named %s !", expectedKey)
    }
    if val, ok := resp.Item[expectedKey].([]float64); !ok {
        t.Fatalf("Response item '%s' should be a []string. Got: '%T'", expectedKey, resp.Item[expectedKey])
    } else {
        if len(val) != len(expectedValue) {
            t.Fatalf("Response item '%s' should be len('%d'). Got: len('%d')", expectedKey, len(expectedValue), len(val))
        }
        expectedValueMap := make(map[float64]float64)
        for i := range expectedValue {
            expectedValueMap[expectedValue[i]] = expectedValue[i]
        }
        for i := range val {
            if _, ok := expectedValueMap[val[i]]; !ok {
                t.Fatalf("Got unexpected value '%f' not in set %v", val[i], expectedValue)
            }
        }
    }
}


func Test_MissingTableName(t * testing.T) {
    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        t.Fatal("Never should have made it to request stage")
    })

    itemReq := NewGetItemRequest()
    itemReq.Search["blah"] = "asdf"

    resp, err := doGetItemTest(itemReq, handler)
    if err == nil {
        t.Fatalf("Error should not be nil. Got response: %v", resp)
    }
    if err != Verification_Error_TableNameEmpty {
        t.Fatalf("Got wrong error. Expected: %v.. got: %v", Verification_Error_TableNameEmpty, err)
    }
}

func Test_MissingSearch(t * testing.T) {
    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        t.Fatal("Never should have made it to request stage")
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"

    resp, err := doGetItemTest(itemReq, handler)
    if err == nil {
        t.Fatalf("Error should not be nil. Got response: %v", resp)
    }
    if err != Verification_Error_SearchEmpty {
        t.Fatalf("Got wrong error. Expected: %v.. got: %v", Verification_Error_SearchEmpty, err)
    }
}

func Test_MissingRegion(t * testing.T) {
    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        t.Fatal("Never should have made it to request stage")
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"
    itemReq.Search["asd"] = "asd"

    ts := httptest.NewTLSServer(handler)
    defer ts.Close()
    certAsx509, _ := x509.ParseCertificate(ts.TLS.Certificates[0].Certificate[0])

    itemReq.Host.Domain = strings.TrimPrefix(ts.URL, "https://127.0.")
    //itemReq.Host.Region = "0" // test is to keep this empty.
    itemReq.Host.Service = "127"
    itemReq.Key.AccessKeyId = "akey"
    itemReq.Key.SecretAccessKey = "skey"
    itemReq.Host.CustomCertificates = []*x509.Certificate{certAsx509}

    resp, err := itemReq.Request()

    if err == nil {
        t.Fatalf("Error should not be nil. Got response: %v", resp)
    }
    if err != Verification_Error_RegionEmpty {
        t.Fatalf("Got wrong error. Expected: %v.. got: %v", Verification_Error_RegionEmpty, err)
    }
}

func Test_MissingService(t * testing.T) {
    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        t.Fatal("Never should have made it to request stage")
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"
    itemReq.Search["asd"] = "asd"

    ts := httptest.NewTLSServer(handler)
    defer ts.Close()
    certAsx509, _ := x509.ParseCertificate(ts.TLS.Certificates[0].Certificate[0])

    itemReq.Host.Domain = strings.TrimPrefix(ts.URL, "https://127.0.")
    itemReq.Host.Region = "0"
    itemReq.Host.Service = "" // test is to set this to empty.
    itemReq.Key.AccessKeyId = "akey"
    itemReq.Key.SecretAccessKey = "skey"
    itemReq.Host.CustomCertificates = []*x509.Certificate{certAsx509}

    resp, err := itemReq.Request()

    if err == nil {
        t.Fatalf("Error should not be nil. Got response: %v", resp)
    }
    if err != Verification_Error_ServiceEmpty {
        t.Fatalf("Got wrong error. Expected: %v.. got: %v", Verification_Error_ServiceEmpty, err)
    }
}

func Test_MissingAccessKey(t * testing.T) {
    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        t.Fatal("Never should have made it to request stage")
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"
    itemReq.Search["asd"] = "asd"

    ts := httptest.NewTLSServer(handler)
    defer ts.Close()
    certAsx509, _ := x509.ParseCertificate(ts.TLS.Certificates[0].Certificate[0])

    itemReq.Host.Domain = strings.TrimPrefix(ts.URL, "https://127.0.")
    itemReq.Host.Region = "0"
    itemReq.Host.Service = "127"
    //itemReq.Key.AccessKeyId = "akey" // test is to have this empty
    itemReq.Key.SecretAccessKey = "skey"
    itemReq.Host.CustomCertificates = []*x509.Certificate{certAsx509}

    resp, err := itemReq.Request()

    if err == nil {
        t.Fatalf("Error should not be nil. Got response: %v", resp)
    }
    if err != awsgo.Verification_Error_AccessKeyEmpty {
        t.Fatalf("Got wrong error. Expected: %v.. got: %v", awsgo.Verification_Error_AccessKeyEmpty, err)
    }
}

func Test_MissingSecretAccessKey(t * testing.T) {
    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        t.Fatal("Never should have made it to request stage")
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"
    itemReq.Search["asd"] = "asd"

    ts := httptest.NewTLSServer(handler)
    defer ts.Close()
    certAsx509, _ := x509.ParseCertificate(ts.TLS.Certificates[0].Certificate[0])

    itemReq.Host.Domain = strings.TrimPrefix(ts.URL, "https://127.0.")
    itemReq.Host.Region = "0"
    itemReq.Host.Service = "127"
    itemReq.Key.AccessKeyId = "akey"
    //itemReq.Key.SecretAccessKey = "skey" // test is to have this empty
    itemReq.Host.CustomCertificates = []*x509.Certificate{certAsx509}

    resp, err := itemReq.Request()

    if err == nil {
        t.Fatalf("Error should not be nil. Got response: %v", resp)
    }
    if err != awsgo.Verification_Error_SecretAccessKeyEmpty {
        t.Fatalf("Got wrong error. Expected: %v.. got: %v", awsgo.Verification_Error_SecretAccessKeyEmpty, err)
    }
}

func Test_InvalidJson(t * testing.T) {
    invalidJson := "hello. This isn't json yo!"

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        fmt.Fprintf(w, invalidJson)
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"
    itemReq.Search["blah"] = "asdf"

    resp, err := doGetItemTest(itemReq, handler)
    if err == nil {
        t.Fatalf("Error should not be nil. Got resp: %v", resp)
    }
    if val, ok := err.(*awsgo.UnmarhsallingError); !ok {
        t.Fatalf("Exepected awsgo.UnmarshallingError. Got: %T", err)
    } else if val.ActualContent != invalidJson {
        t.Fatalf("Expected error content to be '%s' got '%s'", invalidJson, val.ActualContent)
    }
}

func Test_EmtpyReturnSet(t * testing.T) {

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        fmt.Fprintf(w, "{}")
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"
    itemReq.Search["blah"] = "asdf"

    resp, err := doGetItemTest(itemReq, handler)
    if err != nil {
        t.Fatalf("Error should be nil. Got: %v", err)
    }
    if len(resp.Item) > 0 {
        t.Fatalf("Response should not have any items! Got %v", resp.Item)
    }
}

func Test_ConsumedCapacityDemarshal(t * testing.T) {
    expectedBody := `
        {"ConsistentRead":"false","Key":{"blah":{"S":"asdf"}},"TableName":"asd","ReturnConsumedCapacity":"TOTAL"}
    `

    expectedCapacityUnits := float64(0.5)
    expected_GS_Really := float64(1)
    expected_GS_OkayWhat := float64(3)
    expected_LS_Something := float64(2.12)
    expected_Table_Capcity := float64(2.1)
    expectedTableName := "asd"

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        defer r.Body.Close()
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            t.Fatal("couldn't read content! error: %v", err)
        }
        if string(body) != strings.TrimSpace(expectedBody) {
            t.Fatal("Expected request body: '%s'. Got: '%s'", strings.TrimSpace(expectedBody), string(body))
        }

        fmt.Fprintf(w, `
            {
                "ConsumedCapacity" : {
                    "CapacityUnits" : %f,
                    "GlobalSecondaryIndexes" : {
                        "Really" : { "CapacityUnits" : %f },
                        "OkayWhat" : { "CapacityUnits" : %f }
                    },
                    "LocalSecondaryIndexes" : {
                        "Something" : { "CapacityUnits" : %f }
                    },
                    "Table" : {
                        "CapacityUnits" : %f
                    },
                    "TableName" : "%s"
                }
            }
            `,
            expectedCapacityUnits,
            expected_GS_Really,
            expected_GS_OkayWhat,
            expected_LS_Something,
            expected_Table_Capcity,
            expectedTableName,
        )
    })

    itemReq := NewGetItemRequest()
    itemReq.ReturnConsumedCapacity = ConsumedCapacity_TOTAL
    itemReq.TableName = "asd"
    itemReq.Search["blah"] = "asdf"

    resp, err := doGetItemTest(itemReq, handler)
    if err != nil {
        t.Fatalf("Error should be nil. Got: %v", err)
    }
    if len(resp.Item) > 0 {
        t.Errorf("Response should not have any items! Got %v", resp.Item)
    }
    if resp.ConsumedCapacity.TableName != expectedTableName {
        t.Errorf("Expected tablename '%s' Got: '%s'", expectedTableName, resp.ConsumedCapacity.TableName)
    }
    if resp.ConsumedCapacity == nil {
        t.Errorf("Consumed Capacity should not be nil")
    }
    if resp.ConsumedCapacity.CapacityUnits != expectedCapacityUnits {
        t.Errorf("Expected %f capacity units. Got %f", resp.ConsumedCapacity.CapacityUnits)
    }
    if len(resp.ConsumedCapacity.GlobalSecondaryIndexes) != 2 {
        t.Errorf("Expected only 2 global secondary indexes. Got %d", len(resp.ConsumedCapacity.GlobalSecondaryIndexes))
    } else {
        if ca, ok := resp.ConsumedCapacity.GlobalSecondaryIndexes["Really"]; !ok {
            t.Errorf("Expected GS index 'Really'.")
        } else if !easyFloatCompare(ca.CapacityUnits, expected_GS_Really) {
            t.Errorf("Expected GS index 'Really' to be %f. Got: %f.", expected_GS_Really, ca.CapacityUnits)
        }
        if ca, ok := resp.ConsumedCapacity.GlobalSecondaryIndexes["OkayWhat"]; !ok {
            t.Errorf("Expected GS index 'OkayWhat'.")
        } else if !easyFloatCompare(ca.CapacityUnits, expected_GS_OkayWhat) {
            t.Errorf("Expected GS index 'OkayWhat' to be %f. Got: %f.", expected_GS_OkayWhat, ca.CapacityUnits)
        }
    }
    if len(resp.ConsumedCapacity.LocalSecondaryIndexes) != 1 {
        t.Errorf("Expected only 1 local secondary indexes. Got %d", len(resp.ConsumedCapacity.LocalSecondaryIndexes))
    } else if ca, ok := resp.ConsumedCapacity.LocalSecondaryIndexes["Something"]; !ok {
        t.Errorf("Expected LS index 'Something'.")
    } else if !easyFloatCompare(ca.CapacityUnits, expected_LS_Something) {
        t.Errorf("Expected LS index 'Something' to be %f. Got: %f.", expected_LS_Something, ca.CapacityUnits)
    }
    if resp.ConsumedCapacity.Table == nil {
        t.Errorf("Expected Table be not nil")
    }else if !easyFloatCompare(resp.ConsumedCapacity.Table.CapacityUnits, expected_Table_Capcity) {
        t.Errorf("Expected Table to be %f. Got %f", expected_Table_Capcity, resp.ConsumedCapacity.Table.CapacityUnits)
    }
}


func Test_BadConsumedCapacityString(t * testing.T) {
    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        expectedBody := `
            {"ConsistentRead":"false","Key":{"blah":{"S":"asdf"}},"TableName":"asd","ReturnConsumedCapacity":"invalidField"}
        `
        defer r.Body.Close()
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            t.Fatal("couldn't read content! error: %v", err)
        }
        if string(body) != strings.TrimSpace(expectedBody) {
            t.Fatal("Expected request body: '%s'. Got: '%s'", strings.TrimSpace(expectedBody), string(body))
        }
        http.Error(w, `
            {"__type":"com.amazon.coral.validate#ValidationException","message":"1 validation error detected: Value 'invalidField' at 'returnConsumedCapacity' failed to satisfy constraint: Member must satisfy enum value set: [INDEXES, TOTAL, NONE]"}
        `, 400)
    })

    itemReq := NewGetItemRequest()
    itemReq.Search["blah"] = "asdf"
    itemReq.ReturnConsumedCapacity = "invalidField"
    itemReq.TableName = "asd"

    resp, err := doGetItemTest(itemReq, handler)
    if err == nil {
        t.Fatalf("Error should not be nil. Got response: %v", resp)
    }
    if errorResult, ok := err.(*ErrorResult); !ok {
        t.Fatalf("Got wrong error. Expected: ErrorResult. got: %T : %v", err, err)
    } else if errorResult.Type != ValidationException {
        t.Fatalf("Got wrong error. Expected: %v. got: %v", ValidationException, err)
    } else if errorResult.StatusCode != 400 {
        t.Fatalf("Expected 400 status code. Got %v", errorResult.StatusCode)
    }
}

func Test_BadKeyName(t * testing.T) {
    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        verifySimpleRequestBody(t, r)
        http.Error(w, `
            {"__type":"com.amazon.coral.validate#ValidationException","message":"The provided key element does not match the schema"}
        `, 400)
    })

    itemReq := NewGetItemRequest()
    itemReq.Search["blah"] = "asdf" // in this test 'blah' isn't in the table
    itemReq.TableName = "asd"

    resp, err := doGetItemTest(itemReq, handler)
    if err == nil {
        t.Fatalf("Error should not be nil. Got response: %v", resp)
    }
    if errorResult, ok := err.(*ErrorResult); !ok {
        t.Fatalf("Got wrong error. Expected: ErrorResult. got: %T : %v", err, err)
    } else if errorResult.Type != ValidationException {
        t.Fatalf("Got wrong error. Expected: %v. got: %v", ValidationException, err)
    } else if errorResult.StatusCode != 400 {
        t.Fatalf("Expected 400 status code. Got %v", errorResult.StatusCode)
    }
}


func Test_ServerError(t * testing.T) {

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        verifySimpleRequestBody(t, r)
        http.Error(w, "", 500) // really don't know what this looks like...
    })

    itemReq := NewGetItemRequest()
    itemReq.TableName = "asd"
    itemReq.Search["blah"] = "asdf"

    resp, err := doGetItemTest(itemReq, handler)
    if err == nil {
        t.Fatalf("Error should  not be nil. Got resp: %v", resp)
    }
    if errorResult, ok := err.(*ErrorResult); !ok {
        t.Fatalf("Expected type ErrorResult. Got %T : %v", err, err)
    } else if errorResult.StatusCode != 500 {
        t.Fatalf("Expected 500 status code. Got %v", errorResult.StatusCode)
    } else if errorResult.Type != UnknownServerError {
        t.Fatalf("Expected UnknownServerError. Got %v", errorResult.Type)
    }
}

func Test_AccessDeniedToTable(t * testing.T) {
    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        verifySimpleRequestBody(t, r)
        http.Error(w, `
            {"__type":"com.amazon.coral.service#AccessDeniedException","Message":"User: arn:aws:iam::34444444:user/tester is not authorized to perform: dynamodb:GetItem on resource: arn:aws:dynamodb:us-west-2:33333333:table/asd"}
        `, 400)
    })

    itemReq := NewGetItemRequest()
    itemReq.Search["blah"] = "asdf"
    itemReq.TableName = "asd"

    resp, err := doGetItemTest(itemReq, handler)
    if err == nil {
        t.Fatalf("Error should not be nil. Got response: %v", resp)
    }
    if errorResult, ok := err.(*ErrorResult); !ok {
        t.Fatalf("Got wrong error. Expected: ErrorResult. got: %T : %v", err, err)
    } else if errorResult.Type != AccessDeniedException {
        t.Fatalf("Got wrong error. Expected: %v. got: %v", AccessDeniedException, err)
    } else if errorResult.StatusCode != 400 {
        t.Fatalf("Expected 400 status code. Got %v", errorResult.StatusCode)
    }
}
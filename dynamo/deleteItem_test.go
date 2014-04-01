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
    //"github.com/fromkeith/awsgo"
    "encoding/json"
    "bytes"
)


func doDeleteItemTest(req *DeleteItemRequest, handler http.HandlerFunc) (*DeleteItemResponse, error) {
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


func Test_WorkingDeleteSingleItem_String(t * testing.T) {

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        expectedRequestBody := `
        {
            "Key" : {
                "blah" : {
                    "S" : "asdf"
                }
            },
            "ReturnConsumedCapacity" : "NONE",
            "ReturnItemCollectionMetrics" : "NONE",
            "ReturnValues" : "NONE",
            "TableName": "asd"
        }
        `
        expectedCompactBuf := bytes.Buffer{}
        json.Compact(&expectedCompactBuf, []byte(expectedRequestBody))

        defer r.Body.Close()
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            t.Fatal("couldn't read content! error: %v", err)
        }

        if expectedCompactBuf.String() != string(body) {
            t.Errorf("Bodies don't match. Expected: %s. Got %s", expectedCompactBuf.String(), string(body))
        }
        fmt.Fprintf(w, "{}")
    })

    itemReq := NewDeleteItemRequest()
    itemReq.TableName = "asd"
    itemReq.DeleteKey["blah"] = "asdf"

    resp, err := doDeleteItemTest(itemReq, handler)
    if err != nil {
        t.Fatalf("Error should be nil. Got: %v", err)
    }
    if len(resp.Attributes) != 0 {
        t.Errorf("Attributes should be 0 length")
    }
    if resp.ConsumedCapacity != nil {
        t.Errorf("ConsumedCapacity should be nil")
    }
    if resp.ItemCollectionMetrics != nil {
        t.Errorf("ItemCollectionMetrics should be nil")
    }
}


func Test_WorkingDeleteMultiItems(t * testing.T) {

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        expectedRequestBody := `
        {
            "Key" : {
                "beatrice" : {
                    "N" : "55.202000"
                },
                "blah" : {
                    "S" : "Some string yo!"
                },
                "burt" : {
                    "NS" : [
                        "234.430000", "43.000000", "43.000000", "77.000000"
                    ]
                },
                "hulk" : {
                    "SS" : [
                        "smash", "everything", "ya?"
                    ]
                }
            },
            "ReturnConsumedCapacity" : "NONE",
            "ReturnItemCollectionMetrics" : "NONE",
            "ReturnValues" : "NONE",
            "TableName": "hello"
        }
        `
        expectedCompactBuf := bytes.Buffer{}
        json.Compact(&expectedCompactBuf, []byte(expectedRequestBody))

        defer r.Body.Close()
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            t.Fatal("couldn't read content! error: %v", err)
        }

        if expectedCompactBuf.String() != string(body) {
            t.Errorf("Bodies don't match. Expected: %s. Got %s", expectedCompactBuf.String(), string(body))
        }
        fmt.Fprintf(w, "{}")
    })

    itemReq := NewDeleteItemRequest()
    itemReq.TableName = "hello"
    itemReq.DeleteKey["blah"] = "Some string yo!"
    itemReq.DeleteKey["hulk"] = []string{"smash", "everything", "ya?"}
    itemReq.DeleteKey["burt"] = []float64{234.43,43,43,77}
    itemReq.DeleteKey["beatrice"] = 55.202

    resp, err := doDeleteItemTest(itemReq, handler)
    if err != nil {
        t.Fatalf("Error should be nil. Got: %v", err)
    }
    if len(resp.Attributes) != 0 {
        t.Errorf("Attributes should be 0 length")
    }
    if resp.ConsumedCapacity != nil {
        t.Errorf("ConsumedCapacity should be nil")
    }
    if resp.ItemCollectionMetrics != nil {
        t.Errorf("ItemCollectionMetrics should be nil")
    }
}



func Test_DeleteItem_BadResponse(t * testing.T) {
    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        http.Error(w, "sdfsdf", 400)
    })

    itemReq := NewDeleteItemRequest()
    itemReq.TableName = "asd"
    itemReq.DeleteKey["blah"] = "asdf"

    resp, err := doDeleteItemTest(itemReq, handler)
    if err == nil {
        t.Fatalf("Error should not be nil.")
    }
    if resp != nil {
        t.Fatalf("Response should be nil")
    }
}


func Test_WorkingDeleteSingleItem_WithExpected(t * testing.T) {

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        expectedRequestBody := `
        {
            "Expected" : {
                "notthere" : {
                    "Exists": "false"
                },
                "there" : {
                    "Exists" : "true",
                    "Value" : {
                        "S" : "dsf"
                    }
                },
                "thereArray" : {
                    "Exists" : "true",
                    "Value" : {
                        "SS" : [
                            "dsf", "sdfs"
                        ]
                    }
                },
                "thereNum" : {
                    "Exists" : "true",
                    "Value" : {
                        "N" : "23.000000"
                    }
                },
                "thereNumArray" : {
                    "Exists" : "true",
                    "Value" : {
                        "NS" : [
                            "23.000000", "43.000000"
                        ]
                    }
                }
            },
            "Key" : {
                "blah" : {
                    "S" : "as"
                }
            },
            "ReturnConsumedCapacity" : "NONE",
            "ReturnItemCollectionMetrics" : "NONE",
            "ReturnValues" : "NONE",
            "TableName": "asd"
        }
        `
        expectedCompactBuf := bytes.Buffer{}
        json.Compact(&expectedCompactBuf, []byte(expectedRequestBody))

        defer r.Body.Close()
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            t.Fatal("couldn't read content! error: %v", err)
        }

        if expectedCompactBuf.String() != string(body) {
            t.Errorf("Bodies don't match. Expected: %s. Got %s", expectedCompactBuf.String(), string(body))
        }
        fmt.Fprintf(w, "{}")
    })

    itemReq := NewDeleteItemRequest()
    itemReq.TableName = "asd"
    itemReq.DeleteKey["blah"] = "as"
    itemReq.Expected["notthere"] = ExpectedItem{false, nil}
    itemReq.Expected["there"] = ExpectedItem{true, "dsf"}
    itemReq.Expected["thereArray"] = ExpectedItem{true, []string{"dsf", "sdfs"}}
    itemReq.Expected["thereNum"] = ExpectedItem{true, 23}
    itemReq.Expected["thereNumArray"] = ExpectedItem{true, []float64{23,43}}

    resp, err := doDeleteItemTest(itemReq, handler)
    if err != nil {
        t.Fatalf("Error should be nil. Got: %v", err)
    }
    if len(resp.Attributes) != 0 {
        t.Errorf("Attributes should be 0 length")
    }
    if resp.ConsumedCapacity != nil {
        t.Errorf("ConsumedCapacity should be nil")
    }
    if resp.ItemCollectionMetrics != nil {
        t.Errorf("ItemCollectionMetrics should be nil")
    }
}


func Test_DeleteItem_ReturnItemCollectionMetrics(t * testing.T) {

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        expectedRequestBody := `
        {
            "Key" : {
                "blah" : {
                    "S" : "asdf"
                }
            },
            "ReturnConsumedCapacity" : "NONE",
            "ReturnItemCollectionMetrics" : "SIZE",
            "ReturnValues" : "NONE",
            "TableName": "asd"
        }
        `
        expectedCompactBuf := bytes.Buffer{}
        json.Compact(&expectedCompactBuf, []byte(expectedRequestBody))

        defer r.Body.Close()
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            t.Fatal("couldn't read content! error: %v", err)
        }

        if expectedCompactBuf.String() != string(body) {
            t.Errorf("Bodies don't match. Expected: %s. Got %s", expectedCompactBuf.String(), string(body))
        }
        fmt.Fprintf(w, `{
            "ItemCollectionMetrics" : {
                "ItemCollectionKey" : {
                    "blah" : {
                        "S" : "Hello There"
                    }
                },
                "SizeEstimateRangeGB" : [
                    "1.2"
                ]
            }
        }`)
    })

    itemReq := NewDeleteItemRequest()
    itemReq.TableName = "asd"
    itemReq.DeleteKey["blah"] = "asdf"
    itemReq.ReturnItemCollectionMetrics = ReturnItemCollection_SIZE

    resp, err := doDeleteItemTest(itemReq, handler)
    if err != nil {
        t.Fatalf("Error should be nil. Got: %v", err)
    }
    if len(resp.Attributes) != 0 {
        t.Errorf("Attributes should be 0 length")
    }
    if resp.ConsumedCapacity != nil {
        t.Errorf("ConsumedCapacity should be nil")
    }
    if resp.ItemCollectionMetrics == nil {
        t.Errorf("ItemCollectionMetrics should not be nil")
    }
    if v, ok := resp.ItemCollectionMetrics.ItemCollectionKey["blah"]; !ok {
        t.Errorf("Expected 'blah' key")
    } else if vs, ok := v.(string); !ok {
        t.Errorf("Expected key to be string. got %T", v)
    } else if vs != "Hello There" {
        t.Errorf("Wrong value for vs")
    }
    if len(resp.ItemCollectionMetrics.SizeEstimateRangeGB) != 1 {
        t.Errorf("Expected 1 size item")
    } else if resp.ItemCollectionMetrics.SizeEstimateRangeGB[0] != "1.2" {
        t.Errorf("Expected value 1.2 got: %s", resp.ItemCollectionMetrics.SizeEstimateRangeGB[0])
    }
}


func Test_DeleteItem_ReturnAttributes(t * testing.T) {

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        expectedRequestBody := `
        {
            "Key" : {
                "blah" : {
                    "S" : "asdf"
                },
                "huh" : {
                    "N": "57.000000"
                },
                "huh2" : {
                    "NS": ["4343.000000", "44.000000"]
                }
            },
            "ReturnConsumedCapacity" : "NONE",
            "ReturnItemCollectionMetrics" : "NONE",
            "ReturnValues" : "ALL_NEW",
            "TableName": "asd"
        }
        `
        expectedCompactBuf := bytes.Buffer{}
        json.Compact(&expectedCompactBuf, []byte(expectedRequestBody))

        defer r.Body.Close()
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            t.Fatal("couldn't read content! error: %v", err)
        }

        if expectedCompactBuf.String() != string(body) {
            t.Errorf("Bodies don't match. Expected: %s. Got %s", expectedCompactBuf.String(), string(body))
        }
        fmt.Fprintf(w, `{
            "Attributes" : {
                "blah" : {
                    "S": "asdf"
                },
                "huh" : {
                    "N": "57"
                },
                "huh2" : {
                    "NS": ["4343", "44"]
                }
            }
        }`)
    })

    itemReq := NewDeleteItemRequest()
    itemReq.TableName = "asd"
    itemReq.DeleteKey["blah"] = "asdf"
    itemReq.DeleteKey["huh"] = 57
    itemReq.DeleteKey["huh2"] = []float64{4343, 44}
    itemReq.ReturnValues = ReturnValues_ALL_NEW

    resp, err := doDeleteItemTest(itemReq, handler)
    if err != nil {
        t.Fatalf("Error should be nil. Got: %v", err)
    }
    if len(resp.Attributes) == 0 {
        t.Errorf("Attributes should not be 0 length")
    }
    if resp.ConsumedCapacity != nil {
        t.Errorf("ConsumedCapacity should be nil")
    }
    if resp.ItemCollectionMetrics != nil {
        t.Errorf("ItemCollectionMetrics should be nil")
    }

    if v, ok := resp.Attributes["blah"]; !ok {
        t.Errorf("Blah not found in return")
    } else if vs, ok := v.(string); !ok {
        t.Errorf("blah should be type string. got: %T", v)
    } else if vs != "asdf" {
        t.Errorf("Expecting asdf got: %s", vs)
    }

    if v, ok := resp.Attributes["huh"]; !ok {
        t.Errorf("huh not found in return")
    } else if vs, ok := v.(float64); !ok {
        t.Errorf("huh should be type float64. got: %T", v)
    } else if vs != 57 {
        t.Errorf("Expecting asdf got: %d", vs)
    }

    if v, ok := resp.Attributes["huh2"]; !ok {
        t.Errorf("huh2 not found in return")
    } else if vs, ok := v.([]float64); !ok {
        t.Errorf("huh2 should be type []float64. got: %T", v)
    } else if len(vs) != 2 {
        t.Errorf("Expecting asdf length 2. got: %d", len(vs))
    }
}


func Test_DeleteItem_NoItem(t * testing.T) {

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        t.Fatal("Should not have been completed")
    })

    itemReq := NewDeleteItemRequest()
    itemReq.TableName = "asd"

    _, err := doDeleteItemTest(itemReq, handler)
    if err == nil {
        t.Fatalf("Error should not be nil.")
    }
}

func Test_DeleteITem_NoTable(t * testing.T) {

    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        t.Fatal("Should not have been completed")
    })

    itemReq := NewDeleteItemRequest()
    itemReq.DeleteKey["asd"] = "Asd"

    _, err := doDeleteItemTest(itemReq, handler)
    if err == nil {
        t.Fatalf("Error should not be nil.")
    }
}


func Test_DeleteItem_BadResponse_GoodStatusCode(t * testing.T) {
    handler := http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
        http.Error(w, "sdfsdf", 200)
    })

    itemReq := NewDeleteItemRequest()
    itemReq.TableName = "asd"
    itemReq.DeleteKey["blah"] = "asdf"

    resp, err := doDeleteItemTest(itemReq, handler)
    if err == nil {
        t.Fatalf("Error should not be nil.")
    }
    if resp != nil {
        t.Fatalf("Response should be nil")
    }
}
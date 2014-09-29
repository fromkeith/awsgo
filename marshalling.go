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
    "fmt"
    "strconv"
)


type awsStringItem struct {
    Value string      `json:"S,omitempty"`
    Values []string   `json:"SS,omitempty"`
}
type awsNumberItem struct {
    Value float64       `json:"-"`
    Values []float64    `json:"-"`
    ValuesStr []string  `json:"NS,omitempty"`
    ValueStr string     `json:"N,omitempty"`
}


func newStringItem(items ... string) awsStringItem {
    var s awsStringItem
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

func newNumberItem(items ... float64) awsNumberItem {
    var s awsNumberItem
    if (len(items) == 1) {
        s.Value = items[0]
        s.ValueStr = fmt.Sprintf("%f", items[0])
    } else {
        s.Values = make([]float64, len(items))
        s.ValuesStr = make([]string, len(items))
        for i, val := range items {
            s.Values[i] = val
            s.ValuesStr[i] = fmt.Sprintf("%f", val)
        }
    }
    return s
}

// Converts from an unknown interface... like:
//     string, []string, float, []float64
// into the expected awsgo.awsStringItem or awsgo.awsNumberItem
func ConvertToAwsItem(unknown interface{}) interface{} {
    switch j := unknown.(type) {
        case string:
            return newStringItem(j)
        case float64:
            return newNumberItem(j)
        case int:
            return newNumberItem(float64(j))
        case uint:
            return newNumberItem(float64(j))
        case float32:
            return newNumberItem(float64(j))
        case int64:
            return newNumberItem(float64(j))
        case uint64:
            return newNumberItem(float64(j))
        case int32:
            return newNumberItem(float64(j))
        case uint32:
            return newNumberItem(float64(j))
        case []string:
            return awsStringItem{"", j}
        case []int:
            // we need to cast these over
            vals64 := make([]float64, len(j))
            for i := range j {
                vals64[i] = float64(j[i])
            }
            return newNumberItem(vals64...)
        case []uint:
            // we need to cast these over
            vals64 := make([]float64, len(j))
            for i := range j {
                vals64[i] = float64(j[i])
            }
            return newNumberItem(vals64...)
        case []int64:
            // we need to cast these over
            vals64 := make([]float64, len(j))
            for i := range j {
                vals64[i] = float64(j[i])
            }
            return newNumberItem(vals64...)
        case []uint64:
            // we need to cast these over
            vals64 := make([]float64, len(j))
            for i := range j {
                vals64[i] = float64(j[i])
            }
            return newNumberItem(vals64...)
        case []int32:
            // we need to cast these over
            vals64 := make([]float64, len(j))
            for i := range j {
                vals64[i] = float64(j[i])
            }
            return newNumberItem(vals64...)
        case []uint32:
            // we need to cast these over
            vals64 := make([]float64, len(j))
            for i := range j {
                vals64[i] = float64(j[i])
            }
            return newNumberItem(vals64...)
        case []float32:
            // we need to cast these over
            vals64 := make([]float64, len(j))
            for i := range j {
                vals64[i] = float64(j[i])
            }
            return newNumberItem(vals64...)
        case []float64:
            return newNumberItem(j...)
        case awsNumberItem:
            return j
        case awsStringItem:
            return j
        default:
            panic(fmt.Sprintf("Unknown data type: %v %T", j, j))
            return j
    }
    return unknown
}

// converts from raw JSON map to the expected types Eg. float64, string
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

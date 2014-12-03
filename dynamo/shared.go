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
    "encoding/json"
    "fmt"
    "strconv"
    "reflect"
    "github.com/fromkeith/awsgo"
    "errors"
    "time"
)

// Variable Constants
const (
    ConsumedCapacity_TOTAL = "TOTAL"
    ConsumedCapacity_NONE = "NONE"
    ConsumedCapacity_INDEXES = "INDEXES"

    ReturnItemCollection_SIZE = "SIZE"
    ReturnItemCollection_NONE = "NONE"

    // apply to update & put
    ReturnValues_NONE = "NONE"
    ReturnValues_ALL_OLD = "ALL_OLD"
    ReturnValues_UPDATED_OLD = "UPDATED_OLD"
    ReturnValues_ALL_NEW = "ALL_NEW"
    ReturnValues_UPDATED_NEW = "UPDATED_NEW"


    ItemCollectionMetrics_SIZE = "SIZE"
    ItemCollectionMetrics_NONE = "NONE"

    ComparisonOperator_EQ = "EQ"
    ComparisonOperator_LE = "LE"
    ComparisonOperator_LT = "LT"
    ComparisonOperator_GE = "GE"
    ComparisonOperator_GT = "GT"
    ComparisonOperator_BEGINS_WITH = "BEGINS_WITH"
    ComparisonOperator_BETWEEN = "BETWEEN"

    Select_ALL_ATTRIBUTES = "ALL_ATTRIBUTES"
    Select_ALL_PROJECTED_ATTRIBUTES = "ALL_PROJECTED_ATTRIBUTES"
    Select_COUNT = "COUNT"
    Select_SPECIFIC_ATTRIBUTES  = "SPECIFIC_ATTRIBUTES"
)


// Targets
const (
    GetItemTarget = "DynamoDB_20120810.GetItem"
    PutItemTarget = "DynamoDB_20120810.PutItem"
    BatchGetItemTarget = "DynamoDB_20120810.BatchGetItem"
    UpdateItemTarget = "DynamoDB_20120810.UpdateItem"
    BatchWriteItemTarget = "DynamoDB_20120810.BatchWriteItem"
    QueryTarget = "DynamoDB_20120810.Query"
    DeleteItemTarget = "DynamoDB_20120810.DeleteItem"
    ScanTarget = "DynamoDB_20120810.Scan"
    DescribeTableTarget = "DynamoDB_20120810.DescribeTable"
    UpdateTableTarget = "DynamoDB_20120810.UpdateTable"
)
// Known Errors
const (
    ConditionalCheckFailed = "com.amazonaws.dynamodb.v20120810#ConditionalCheckFailedException"
    SerializationException = "com.amazon.coral.service#SerializationException"
    ValidationException = "com.amazon.coral.validate#ValidationException"
    UnknownServerError = "UnknownServerError"
    AccessDeniedException = "com.amazon.coral.service#AccessDeniedException"
    ThroughputException = "com.amazonaws.dynamodb.v20120810#ProvisionedThroughputExceededException"
)

type CapacityUnitsStruct struct {
    CapacityUnits   float64
}

type CapacityResult struct {
    CapacityUnits   float64
    TableName       string
    Table           *CapacityUnitsStruct
    GlobalSecondaryIndexes  map[string]CapacityUnitsStruct
    LocalSecondaryIndexes   map[string]CapacityUnitsStruct
}

type ExpectedItem struct {
    Exists  bool        `json:",string"`
    Value   interface{} `json:",omitempty"`
}

type KeyConditions struct {
    AttributeValueList      []interface{}
    ComparisonOperator      string
}


type ItemCollectionMetricsStruct struct {
    RawItemCollectionKey        map[string]map[string]interface{}   `json:"ItemCollectionKey"`
    ItemCollectionKey           map[string]interface{}  `json:"-"`
    SizeEstimateRangeGB         []string
}

type ErrorResult struct {
    Type        string  `json:"__type"`
    Message     string  `json:"message"`
    StatusCode  int
}

func (e * ErrorResult) Error() string {
    return fmt.Sprintf("%s : %s", e.Type, e.Message)
}

func CheckForErrorResponse(response []byte, statusCode int) error {
    errorResult := new(ErrorResult)
    err2 := json.Unmarshal([]byte(response), errorResult)
    if err2 == nil {
        if errorResult.Type != "" {
            errorResult.StatusCode = statusCode
            return errorResult
        }
    }
    if statusCode < 200 || statusCode > 299 {
        errorResult.Type = UnknownServerError
        errorResult.StatusCode = statusCode
        return errorResult
    }
    return nil
}

// returns the value in the map if it exists, otherise 'elze' value is returned
func AsStringOr(item map[string]interface{}, key, elze string) string {
    if v, ok := item[key].(string); ok {
        return v
    }
    return elze
}

// returns the value in the map if it exists, otherise 'elze' value is returned
func AsFloatOr(item map[string]interface{}, key string, elze float64) float64 {
    if v, ok := item[key].(float64); ok {
        return v
    }
    return elze
}

// parses the value to be a boolean. using strconv.
// if it is a float then it returns true on the value != 0
// if the value doesn't exit, returns elze
func AsBoolOr(item map[string]interface{}, key string, elze bool) bool {
    if v, ok := item[key].(string); ok {
        if b, bok := strconv.ParseBool(v); bok == nil {
            return b
        }
        return elze
    }
    if v, ok := item[key].(float64); ok {
        return v != 0
    }
    return elze
}


// Unmarshalls a JSON response from AWS.
func Unmarshal(in map[string]map[string]interface{}, out interface{}) error {
    reflectVal := reflect.ValueOf(out)
    if !reflectVal.IsValid() {
        return errors.New("Out is not valid")
    }
    reflectType := reflectVal.Type()
    if reflectType.Kind() != reflect.Ptr {
        return errors.New("Out is not valid pointer")
    }
    reflectVal = reflectVal.Elem()
    reflectType = reflectVal.Type()
    if reflectType.Kind() != reflect.Struct {
        return errors.New("Out is not valid pointer to a struct")
    }
    for i := 0; i < reflectType.NumField(); i++ {
        f := reflectType.Field(i)
        if f.PkgPath != "" {
            continue // unexported
        }
        tag := f.Tag.Get("dynamo")
        if tag == "-" {
            continue
        }
        name := tag
        if name == "" {
            name = f.Name
        }
        switch f.Type.Kind() {
        case reflect.String:
            if asStr, ok := in[name]["S"].(string); ok {
                reflectVal.FieldByIndex(f.Index).SetString(asStr)
            }
            break
        case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
            if asNum, ok := in[name]["N"].(string); ok {
                asInt, err := strconv.ParseInt(asNum, 10, 64)
                if err != nil {
                    return err
                }
                reflectVal.FieldByIndex(f.Index).SetInt(asInt)
            }
            break
        case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
            if asNum, ok := in[name]["N"].(string); ok {
                asInt, err := strconv.ParseUint(asNum, 10, 64)
                if err != nil {
                    return err
                }
                reflectVal.FieldByIndex(f.Index).SetUint(asInt)
            }
            break
        case reflect.Float32, reflect.Float64:
            if asNum, ok := in[name]["N"].(string); ok {
                asFloat, err := strconv.ParseFloat(asNum, 64)
                if err != nil {
                    return err
                }
                reflectVal.FieldByIndex(f.Index).SetFloat(asFloat)
            }
            break
        case reflect.Array, reflect.Slice:
            err := decodeArray(in[name], reflectVal.FieldByIndex(f.Index))
            if err != nil {
                return err
            }
            break
        default:
            hasItem, ok := in[name]
            if !ok {
                break
            }
            if _, ok := reflectVal.FieldByIndex(f.Index).Interface().(time.Time); ok {
                asTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", hasItem["S"].(string))
                if err != nil {
                    return err
                }
                reflectVal.FieldByIndex(f.Index).Set(reflect.ValueOf(asTime))
                break
            }
            if asStr, ok := hasItem["S"].(string); ok {
                as := reflect.New(f.Type)
                err := json.Unmarshal([]byte(asStr), as.Interface())
                if err != nil {
                    return err
                }
                reflectVal.FieldByIndex(f.Index).Set(as.Elem())
            } else {
                return errors.New("Cannot decode field: " + name)
            }
        }
    }
    return nil
}

func decodeArray(in map[string]interface{}, v reflect.Value) error {
    if in == nil {
        return nil
    }
    k := v.Type().Elem().Kind()
    t := v.Type()
    switch k {
    case reflect.String:
        if asStr, ok := in["SS"].([]interface{}); ok {
            v.Set(reflect.MakeSlice(t, len(asStr), len(asStr)))
            for i := range asStr {
                v.Index(i).SetString(asStr[i].(string))
            }
        }
        return nil
        break
    case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
        if asNum, ok := in["NS"].([]interface{}); ok {
            v.Set(reflect.MakeSlice(t, len(asNum), len(asNum)))
            for i := range asNum {
                asInt, err := strconv.ParseInt(asNum[i].(string), 10, 64)
                if err != nil {
                    return err
                }
                v.Index(i).SetInt(asInt)
            }
        }
        return nil
    case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
        if asNum, ok := in["NS"].([]interface{}); ok {
            v.Set(reflect.MakeSlice(t, len(asNum), len(asNum)))
            for i := range asNum {
                asInt, err := strconv.ParseUint(asNum[i].(string), 10, 64)
                if err != nil {
                    return err
                }
                v.Index(i).SetUint(asInt)
            }
        }
        return nil
    case reflect.Float32, reflect.Float64:
        if asNum, ok := in["NS"].([]interface{}); ok {
            v.Set(reflect.MakeSlice(t, len(asNum), len(asNum)))
            for i := range asNum {
                asFloat, err := strconv.ParseFloat(asNum[i].(string), 64)
                if err != nil {
                    return err
                }
                v.Index(i).SetFloat(asFloat)
            }
        }
        return nil
    }
    if asStr, ok := in["S"].(string); ok {
        as := reflect.New(t)
        err := json.Unmarshal([]byte(asStr), as.Interface())
        if err != nil {
            return err
        }
        v.Set(as.Elem())
        return err
    }
    return errors.New("Could not decode item")
}


// takes an struct and returns a map that can be used write or put item with
//      you can rename a field via: `dynamo:"rename"` tag
//      fields can be omitted via: `dynamo:"-"` tag
//      empty strings are not marshalled, same for empty arrays
//      non basic types (eg interfaces, structs, pointers) are marshalled as json string
func Marshal(v interface{}) map[string]interface{} {
    reflectVal := reflect.ValueOf(v)
    if !reflectVal.IsValid() {
        return nil
    }
    reflectType := reflectVal.Type()
    if reflectType.Kind() != reflect.Struct {
        return nil
    }
    result := make(map[string]interface{})
    for i := 0; i < reflectType.NumField(); i++ {
        f := reflectType.Field(i)
        if f.PkgPath != "" {
            continue // unexported
        }
        tag := f.Tag.Get("dynamo")
        if tag == "-" {
            continue
        }
        name := tag
        if name == "" {
            name = f.Name
        }
        switch f.Type.Kind() {
        case reflect.String:
            val := reflectVal.FieldByIndex(f.Index).String()
            if val == "" {
                continue
            }
            result[name] = awsgo.AwsStringItem{
                Value: val,
            }
            break
        case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
            val := reflectVal.FieldByIndex(f.Index).Int()
            result[name] = awsgo.AwsNumberItem{
                Value: float64(val), // backwards compat
                ValueStr: strconv.FormatInt(val, 10),
            }
            break
        case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
            val := reflectVal.FieldByIndex(f.Index).Uint()
            result[name] = awsgo.AwsNumberItem{
                Value: float64(val), // backwards compat
                ValueStr: strconv.FormatUint(val, 10),
            }
            break
        case reflect.Float32, reflect.Float64:
            val := reflectVal.FieldByIndex(f.Index).Float()
            result[name] = awsgo.AwsNumberItem{
                Value: val, // backwards compat
                ValueStr: fmt.Sprintf("%f", val),
            }
            break
        case reflect.Array, reflect.Slice:
            val := encodeArray(reflectVal.FieldByIndex(f.Index))
            if val == nil {
                continue
            }
            result[name] = val
            break
        default:
            if asTime, ok := reflectVal.FieldByIndex(f.Index).Interface().(time.Time); ok {
                result[name] = awsgo.AwsStringItem{
                    Value: asTime.String(),
                }
                break
            }
            val, _ := json.Marshal(reflectVal.FieldByIndex(f.Index).Interface())
            result[name] = awsgo.AwsStringItem{
                Value: string(val),
            }
        }
    }
    return result
}

func encodeArray(v reflect.Value) interface{} {
    if v.Len() == 0 {
        return nil
    }
    switch v.Type().Elem().Kind() {
    case reflect.String:
        return awsgo.AwsStringItem{
            Values: v.Interface().([]string),
        }
    case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
        strArray := make([]string, v.Len())
        for i := range strArray {
            val := v.Index(i).Int()
            strArray[i] = strconv.FormatInt(val, 10)
        }
        return awsgo.AwsNumberItem{
            ValuesStr: strArray,
        }
    case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
        strArray := make([]string, v.Len())
        for i := range strArray {
            val := v.Index(i).Uint()
            strArray[i] = strconv.FormatUint(val, 10)
        }
        return awsgo.AwsNumberItem{
            ValuesStr: strArray,
        }
    case reflect.Float32, reflect.Float64:
        strArray := make([]string, v.Len())
        for i := range strArray {
            val := v.Index(i).Float()
            strArray[i] = fmt.Sprintf("%f", val)
        }
        return awsgo.AwsNumberItem{
            ValuesStr: strArray,
        }
    }
    enc, _ := json.Marshal(v.Interface())
    return awsgo.AwsStringItem{
        Value: string(enc),
    }
}



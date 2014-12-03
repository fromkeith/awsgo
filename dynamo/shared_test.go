package dynamo

import (
    "testing"
    "github.com/fromkeith/awsgo"
    "encoding/json"
    "log"
)



type BasicString struct {
    Bob             string
    Jim             string
}

// simulates it as dynamo json, and our decoding of it
func resToJsonAndBack(in map[string]interface{}) (out map[string]map[string]interface{}) {
    b, err := json.Marshal(in)
    if err != nil {
        panic(err)
    }
    log.Println(string(b))
    err = json.Unmarshal(b, &out)
    if err != nil {
        panic(err)
    }
    return
}

func TestMarshalStringStruct(t *testing.T) {
    b := BasicString{
        Bob: "Is cool",
        Jim: "Is way cooler yo!",
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    bob, ok := res["Bob"].(awsgo.AwsStringItem)
    t.Logf("Bob is: %v", res["Bob"])
    if !ok || bob.Value != b.Bob {
        t.Fail()
    }
    jim, ok := res["Jim"].(awsgo.AwsStringItem)
    t.Logf("Jim is: %v", res["Jim"])
    if !ok || jim.Value != b.Jim {
        t.Fail()
    }
    if len(res) != 2 {
        t.Logf("Bad resulting size: %d", len(res))
        t.Fail()
    }
}

func TestUnmarshalStringStruct(t *testing.T) {
    b := BasicString{
        Bob: "Is cool",
        Jim: "Is way cooler yo!",
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    outB := BasicString{}
    err := Unmarshal(resToJsonAndBack(res), &outB)
    if err != nil {
        t.Logf("Error unmarshalling: %v", err)
        t.Fail()
    }
    if outB.Bob != b.Bob {
        t.Logf("Expected Bob to be '%s' got '%s'", b.Bob, outB.Bob)
        t.Fail()
    }
    if outB.Jim != b.Jim {
        t.Logf("Expected Jim to be '%s' got '%s'", b.Jim, outB.Jim)
        t.Fail()
    }
}

type BasicStringRenamed struct {
    Bob             string      `dynamo:"Sherry"`
    Jim             string
}

func TestMarshalStringStructRenameField(t *testing.T) {
    b := BasicStringRenamed{
        Bob: "Is cool",
        Jim: "Is way cooler yo!",
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    _, ok := res["Bob"].(awsgo.AwsStringItem)
    t.Logf("Bob (should be empty) is: %v", res["Bob"])
    if ok {
        t.Fail()
    }
    _, ok = res["Sherry"].(awsgo.AwsStringItem)
    t.Logf("Sherry is: %v", res["Sherry"])
    if !ok {
        t.Fail()
    }
    _, ok = res["Jim"].(awsgo.AwsStringItem)
    t.Logf("Jim is: %v", res["Jim"])
    if !ok {
        t.Fail()
    }

    if len(res) != 2 {
        t.Logf("Bad resulting size: %d", len(res))
        t.Fail()
    }
}

func TestUnmarshalStringStructRenameField(t *testing.T) {
    b := BasicStringRenamed{
        Bob: "Is cool",
        Jim: "Is way cooler yo!",
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    out := BasicStringRenamed{}
    if err := Unmarshal(resToJsonAndBack(res), &out); err != nil {
        t.Logf("Erorr unmarshallign: %v", err)
        t.Fail()
    }
    if b.Bob != out.Bob {
        t.Logf("Bobs do not match")
        t.Fail()
    }
    if b.Jim != out.Jim {
        t.Logf("Jims do not match")
        t.Fail()
    }
}



type BasicStringFieldOmitted struct {
    Bob             string          `dynamo:"-"`
    Jim             string
}

func TestMarshalStringStructFieldOmitted(t *testing.T) {
    b := BasicStringFieldOmitted{
        Bob: "Is cool",
        Jim: "Is way cooler yo!",
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    _, ok := res["Bob"].(awsgo.AwsStringItem)
    t.Logf("Bob is: %v", res["Bob"])
    if ok {
        t.Fail()
    }
    _, ok = res["Jim"].(awsgo.AwsStringItem)
    t.Logf("Jim is: %v", res["Jim"])
    if !ok {
        t.Fail()
    }
    if len(res) != 1 {
        t.Logf("Bad resulting size: %d", len(res))
        t.Fail()
    }
}


func TestUnmarshalStringStructFieldOmitted(t *testing.T) {
    b := BasicStringFieldOmitted{
        Bob: "Is cool",
        Jim: "Is way cooler yo!",
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    out := BasicStringFieldOmitted{}
    if err := Unmarshal(resToJsonAndBack(res), &out); err != nil {
        t.Logf("Erorr unmarshallign: %v", err)
        t.Fail()
    }
    if b.Bob == out.Bob {
        t.Logf("Bobs should not match")
        t.Fail()
    }
    if b.Jim != out.Jim {
        t.Logf("Jims do not match")
        t.Fail()
    }
}


type BasicStringFieldHidden struct {
    bob             string
    Jim             string
}

func TestMarshalStringStructFieldHidden(t *testing.T) {
    b := BasicStringFieldHidden{
        bob: "Is cool",
        Jim: "Is way cooler yo!",
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    _, ok := res["Bob"].(awsgo.AwsStringItem)
    t.Logf("Bob is: %v", res["Bob"])
    if ok {
        t.Fail()
    }
    _, ok = res["bob"].(awsgo.AwsStringItem)
    t.Logf("bob is: %v", res["bob"])
    if ok {
        t.Fail()
    }
    _, ok = res["Jim"].(awsgo.AwsStringItem)
    t.Logf("Jim is: %v", res["Jim"])
    if !ok {
        t.Fail()
    }
    if len(res) != 1 {
        t.Logf("Bad resulting size: %d", len(res))
        t.Fail()
    }
}

func TestUnmarshalStringStructFieldHidden(t *testing.T) {
    b := BasicStringFieldHidden{
        bob: "Is cool",
        Jim: "Is way cooler yo!",
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    out := BasicStringFieldHidden{}
    if err := Unmarshal(resToJsonAndBack(res), &out); err != nil {
        t.Logf("Erorr unmarshallign: %v", err)
        t.Fail()
    }
    if b.bob == out.bob {
        t.Logf("Bobs should not match")
        t.Fail()
    }
    if b.Jim != out.Jim {
        t.Logf("Jims do not match")
        t.Fail()
    }
}




type BasicInt struct {
    Amanda          int
    Barbara         int16
    Carol           int32
    Demi            int64
    Erica           int8
}

func TestMarshalIntStruct(t *testing.T) {
    b := BasicInt{
        Amanda: 5,
        Barbara: 6,
        Carol: 7,
        Demi: 8544654645634,
        Erica: 9,
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    am, ok := res["Amanda"].(awsgo.AwsNumberItem)
    t.Logf("Amanda is: %v", res["Amanda"])
    if !ok || am.ValueStr != "5" {
        t.Fail()
    }
    bar, ok := res["Barbara"].(awsgo.AwsNumberItem)
    t.Logf("Barbara is: %v", res["Barbara"])
    if !ok || bar.ValueStr != "6" {
        t.Fail()
    }
    _, ok = res["Carol"].(awsgo.AwsNumberItem)
    t.Logf("Carol is: %v", res["Carol"])
    if !ok {
        t.Fail()
    }
    dem, ok := res["Demi"].(awsgo.AwsNumberItem)
    t.Logf("Demi is: %v", res["Demi"])
    if !ok || dem.ValueStr != "8544654645634"{
        t.Fail()
    }
    er, ok := res["Erica"].(awsgo.AwsNumberItem)
    t.Logf("Erica is: %v", res["Erica"])
    if !ok || er.ValueStr != "9" {
        t.Fail()
    }
    if len(res) != 5 {
        t.Logf("Bad resulting size: %d", len(res))
        t.Fail()
    }
}



func TestUnmarshalIntStruct(t *testing.T) {
    b := BasicInt{
        Amanda: 5,
        Barbara: 6,
        Carol: 7,
        Demi: 8544654645634,
        Erica: 9,
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    out := BasicInt{}
    if err := Unmarshal(resToJsonAndBack(res), &out); err != nil {
        t.Logf("Failed to umarshal: %v", err)
        t.Fail()
    }
    if b.Amanda != out.Amanda {
        t.Logf("Amanda does not match")
        t.Fail()
    }
    if b.Barbara != out.Barbara {
        t.Logf("Barbara does not match")
        t.Fail()
    }
    if b.Carol != out.Carol {
        t.Logf("Carol does not match")
        t.Fail()
    }
    if b.Demi != out.Demi {
        t.Logf("Demi does not match")
        t.Fail()
    }
    if b.Erica != out.Erica {
        t.Logf("Erica does not match")
        t.Fail()
    }
}

type BasicUint struct {
    Amanda          uint
    Barbara         uint16
    Carol           uint32
    Demi            uint64
    Erica           uint8
}

func TestMarshalUintStruct(t *testing.T) {
    b := BasicUint{
        Amanda: 5,
        Barbara: 6,
        Carol: 7,
        Demi: 8,
        Erica: 9,
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    _, ok := res["Amanda"].(awsgo.AwsNumberItem)
    t.Logf("Amanda is: %v", res["Amanda"])
    if !ok {
        t.Fail()
    }
    _, ok = res["Barbara"].(awsgo.AwsNumberItem)
    t.Logf("Barbara is: %v", res["Barbara"])
    if !ok {
        t.Fail()
    }
    _, ok = res["Carol"].(awsgo.AwsNumberItem)
    t.Logf("Carol is: %v", res["Carol"])
    if !ok {
        t.Fail()
    }
    _, ok = res["Demi"].(awsgo.AwsNumberItem)
    t.Logf("Demi is: %v", res["Demi"])
    if !ok {
        t.Fail()
    }
    _, ok = res["Erica"].(awsgo.AwsNumberItem)
    t.Logf("Erica is: %v", res["Erica"])
    if !ok {
        t.Fail()
    }
    if len(res) != 5 {
        t.Logf("Bad resulting size: %d", len(res))
        t.Fail()
    }
}




func TestUnmarshalUintStruct(t *testing.T) {
    b := BasicUint{
        Amanda: 5,
        Barbara: 6,
        Carol: 7,
        Demi: 8,
        Erica: 9,
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    out := BasicUint{}
    if err := Unmarshal(resToJsonAndBack(res), &out); err != nil {
        t.Logf("Execption: %v", err)
        t.Fail()
    }
    if b.Amanda != out.Amanda {
        t.Logf("Amanda does not match")
        t.Fail()
    }
    if b.Barbara != out.Barbara {
        t.Logf("Barbara does not match")
        t.Fail()
    }
    if b.Carol != out.Carol {
        t.Logf("Carol does not match")
        t.Fail()
    }
    if b.Demi != out.Demi {
        t.Logf("Demi does not match")
        t.Fail()
    }
    if b.Erica != out.Erica {
        t.Logf("Erica does not match")
        t.Fail()
    }
}

type BasicFloat struct {
    Bob             float64
    Jim             float32
}

func TestMarshalFloatStruct(t *testing.T) {
    b := BasicFloat{
        Bob: 34.432455,
        Jim: 55647.345,
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    _, ok := res["Bob"].(awsgo.AwsNumberItem)
    t.Logf("Bob is: %v", res["Bob"])
    if !ok {
        t.Fail()
    }
    _, ok = res["Jim"].(awsgo.AwsNumberItem)
    t.Logf("Jim is: %v", res["Jim"])
    if !ok {
        t.Fail()
    }
    if len(res) != 2 {
        t.Logf("Bad resulting size: %d", len(res))
        t.Fail()
    }
}

func TestUnmarshalFloatStruct(t *testing.T) {
    b := BasicFloat{
        Bob: 34.432455,
        Jim: 55647.345,
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    out := BasicFloat{}
    if err := Unmarshal(resToJsonAndBack(res), &out); err != nil {
        t.Logf("Exception: %v", err)
        t.Fail()
    }
    if b.Bob != out.Bob {
        t.Logf("Bob does not match")
        t.Fail()
    }
    if b.Jim != out.Jim {
        t.Logf("Jim does not match")
        t.Fail()
    }
}



type BasicArrays struct {
    Bob             []string
    Jim             []int64
}

func TestMarshalBasicArraysStruct(t *testing.T) {
    b := BasicArrays{
        Bob: []string{"Is cool", "beyond", "awesome"},
        Jim: []int64{64, 342, 52},
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    bob, ok := res["Bob"].(awsgo.AwsStringItem)
    t.Logf("Bob is: %v", res["Bob"])
    if !ok {
        t.Fail()
    }
    if len(bob.Values) != 3 {
        t.Logf("Bob should be 3 length: %d", len(bob.Values))
        t.Fail()
    }
    j, ok := res["Jim"].(awsgo.AwsNumberItem)
    t.Logf("Jim is: %v", res["Jim"])
    if !ok {
        t.Fail()
    }
    if len(j.ValuesStr) != 3 {
        t.Logf("Jim should be 3 length: %d", len(j.ValuesStr))
        t.Fail()
    }
    if len(res) != 2 {
        t.Logf("Bad resulting size: %d", len(res))
        t.Fail()
    }
}

func TestUnmarshalBasicArraysStruct(t *testing.T) {
    b := BasicArrays{
        Bob: []string{"Is cool", "beyond", "awesome"},
        Jim: []int64{64, 342, 52},
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    t.Logf("%#v", res)
    out := BasicArrays{}
    if err := Unmarshal(resToJsonAndBack(res), &out) ; err != nil {
        t.Logf("Exception: %v", err)
        t.Fail()
    }
    if len(b.Bob) != len(out.Bob) {
        t.Logf("lengths don't match %d. %d", len(b.Bob), len(out.Bob))
        t.FailNow()
    }
    if len(b.Jim) != len(out.Jim) {
        t.Logf("lengths don't match %d. %d", len(b.Jim), len(out.Jim))
        t.FailNow()
    }
    for i := range b.Bob {
        if b.Bob[i] != out.Bob[i] {
            t.Logf("Element %d does not match.", i)
            t.FailNow()
        }
    }
    for i := range b.Jim {
        if b.Jim[i] != out.Jim[i] {
            t.Logf("Element %d does not match.", i)
            t.FailNow()
        }
    }
}


// we expect no items to be return as strings are empty
func TestMarshalEmptyStringStruct(t *testing.T) {
    b := BasicString{
        Bob: "",
        Jim: "",
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    _, ok := res["Bob"].(string)
    t.Logf("Bob is: %v", res["Bob"])
    if ok {
        t.Fail()
    }
    _, ok = res["Jim"].(string)
    t.Logf("Jim is: %v", res["Jim"])
    if ok {
        t.Fail()
    }
    if len(res) != 0 {
        t.Logf("Bad resulting size: %d", len(res))
        t.Fail()
    }
}


// we expect no items to be return as strings are empty
func TestMarshalEmptyArrayStruct(t *testing.T) {
    b := BasicArrays{
        Bob: []string{},
        Jim: []int64{},
    }
    res := Marshal(b)
    if res == nil {
        t.Log("Result is nil!")
        t.Fail()
    }
    _, ok := res["Bob"].([]string)
    t.Logf("Bob is: %v", res["Bob"])
    if ok {
        t.Fail()
    }
    _, ok = res["Jim"].([]int64)
    t.Logf("Jim is: %v", res["Jim"])
    if ok {
        t.Fail()
    }
    if len(res) != 0 {
        t.Logf("Bad resulting size: %d", len(res))
        t.Fail()
    }
}

type Level1 struct {
    Me          string
    Part2       Level2
    Weird       []Level2
}
type Level2 struct {
    Them        string
}

func TestMarshalMutliLevelStruct(t *testing.T) {
    b := Level1{
        Me: "Weeeeee",
        Part2: Level2{
            Them: "Boooo",
        },
        Weird: []Level2{
            Level2{
                Them: "once",
            },
            Level2{
                Them: "twice",
            },
        },
    }
    res := Marshal(b)
    if len(res) != 3 {
        t.Logf("size should be 2: %v", res)
        t.Fail()
    }
    if _, ok := res["Me"].(awsgo.AwsStringItem); !ok {
        t.Logf("Me not found")
        t.Fail()
    }
    if asStr, ok := res["Part2"].(awsgo.AwsStringItem); !ok {
        t.Logf("Part2 not found")
        t.Fail()
    } else if asStr.Value != `{"Them":"Boooo"}` {
        t.Logf("Part2 = %s", asStr.Value)
        t.Fail()
    }
    if asStr, ok := res["Weird"].(awsgo.AwsStringItem); !ok {
        t.Logf("Weird not found")
        t.Fail()
    } else if asStr.Value != `[{"Them":"once"},{"Them":"twice"}]` {
        t.Logf("Weird = %s", asStr.Value)
        t.Fail()
    }
}


func TestUnmarshalMutliLevelStruct(t *testing.T) {
    b := Level1{
        Me: "Weeeeee",
        Part2: Level2{
            Them: "Boooo",
        },
        Weird: []Level2{
            Level2{
                Them: "once",
            },
            Level2{
                Them: "twice",
            },
        },
    }
    res := Marshal(b)
    if len(res) != 3 {
        t.Logf("size should be 3: %v", res)
        t.Fail()
    }
    out := Level1{}
    if err := Unmarshal(resToJsonAndBack(res), &out); err != nil {
        t.Logf("Exeption: %v", err)
        t.Fail()
    }
    if out.Me != b.Me {
        t.Logf("Me no match")
        t.Fail()
    }
    if out.Part2.Them != b.Part2.Them {
        t.Logf("Them no match")
        t.Fail()
    }
    if len(b.Weird) != len(out.Weird) {
        t.Logf("Weird does not have same lengths: %d, %d", len(b.Weird), len(out.Weird))
        t.FailNow()
    }
    for i := range b.Weird {
        if b.Weird[i] != out.Weird[i] {
            t.Logf("Index %d bad", i)
            t.FailNow()
        }
    }
}
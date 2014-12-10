/*
 * Copyright (c) 2014, fromkeith
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

package main


import (
    "github.com/fromkeith/awsgo/swf"
    "github.com/fromkeith/awsgo/swf/swfhelper"
    "strconv"
    "math/rand"
    "log"
    "time"
)

func main() {

    decider := &swfhelper.Decider{
        Region: "us-west-2",
        Domain: "yewwdev",
        TaskList: "hello1",
        Identity: "atHome1",
        MaxWorkers: 1,
    }
    decider.RegisterWorkflow(swf.WorkflowType{Name:"hello", Version:"1"}, BasicDecider)
    go decider.Start()

    act := &swfhelper.ActivityWorker{
        Region: "us-west-2",
        Domain: "yewwdev",
        TaskList: "actlist",
        Identity: "atHome1-worker",
        MaxWorkers: 1,
    }
    act.RegisterActivity(swf.ActivityType{Name:"getnum", Version:"1"}, ActivityOne)
    act.Start()
}



func BasicDecider(w *swfhelper.SwfWorkflow) {
    defer w.Decide()

    task := swfhelper.Task{
        Activity: swf.ActivityType{
            Name:"getnum", Version:"1",
        },
    }
    helloTask := w.Go(task, "hello1")
    res, ok := <- helloTask
    if !ok {
        return
    }
    log.Println("1. Error:", res.FailureType, res.FailureCause)
    log.Println("1. Result:", res.Result)

    helloTask2 := w.Go(task, "hello2")
    helloTask3 := w.Go(task, "hello3")
    res2, ok := <- helloTask2
    if !ok {
        log.Println("Task2 not ready")
        return
    }
    res3, ok := <- helloTask3
    if !ok {
        log.Println("Task3 not ready")
        return
    }

    log.Println("2. Error:", res2.FailureType, res2.FailureCause)
    log.Println("2. Result:", res2.Result)
    log.Println("3. Error:", res3.FailureType, res3.FailureCause)
    log.Println("3. Result:", res3.Result)

    w.Complete("Just amazing")
}

func ActivityOne(a *swfhelper.ActivityContext) {
    myRes := int64(rand.Int())
    defer a.Completed(strconv.FormatInt(myRes, 10))
    defer log.Println("Hey:", a.Input, myRes)
    time.Sleep(5 * time.Second + time.Second * time.Duration(rand.Float32() * 5))
}
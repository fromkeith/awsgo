/*

The SwfHelper is meant to make the decision process execution task easier.
As decision tasks can go to any decider, they need to be able to easily
recreate the history, and make decisions.


Of the swf helper there are two main types. Deciders, and Activities.

Deciders
=========
Deciders must be deterministic in flow. A decider is registered by given the workflow
type and a function callback.

    // this creates the decider manager/poller
    // it will poll the specified region, domain and task list to get decision tasks
    decider := &swfhelper.Decider{
        Region: "us-west-2",
        Domain: "yewwdev",
        TaskList: "hello1",
        Identity: "atHome1",
        MaxWorkers: 1,
    }
    // We must register a function to a workflow. Any decision tasks that come in,
    // and match this type, and version will be handled by the 'BasicDecider' function
    decider.RegisterWorkflow(swf.WorkflowType{Name:"hello", Version:"1"}, BasicDecider)
    // we now block constantly polling SWF for decision tasks
    decider.Start()

The decider itself, is just a function like so:

    // the swfhelper.SwfWorkflow is the object we will manipulate to our needs.
    // This function will be called multiple times during the lifetime of a single
    // workflow. Log lines, or anything not contained in an activity will be executed
    // multiple times, as each time a new decision is required for a workflow, this function
    // will be called. Thus the ordering of w.Go calls must be deterministic.
    func BasicDecider(w *swfhelper.SwfWorkflow) {
        // we must always make a decision after a workflow
        // otherwise swfhelper will not respond with a decision
        defer w.Decide()

        // first we define a task/activity that we want to execute
        task := swfhelper.Task{
            Activity: swf.ActivityType{
                Name:"getnum", Version:"1",
            },
        }
        // we then ask for the execution of it, passing in the 'hello1' string as a variable
        // this will return a channel. If we have not executed the task yet, or the task is
        // not yet complete, a closed channel will be returned. If we are wanting to wait
        // for the response of this task, we must return on a closed channel.
        helloTask := w.Go(task, "hello1")
        res, ok := <- helloTask
        if !ok {
            // return as the channel is closed, and we need the response to continue
            return
        }
        // logging the output of the first task
        // note: that this log line wil appear again after helloTask2 is completed
        log.Println("1. Error:", res.FailureType, res.FailureCause)
        log.Println("1. Result:", res.Result)

        // here we will execute two tasks in parallel. and then wait for their responses
        helloTask2 := w.Go(task, "hello2")
        helloTask3 := w.Go(task, "hello3")
        // wait for hello2
        res2, ok := <- helloTask2
        if !ok {
            log.Println("Task2 not ready")
            return
        }
        // wait for hello3
        res3, ok := <- helloTask3
        if !ok {
            log.Println("Task3 not ready")
            return
        }

        log.Println("2. Error:", res2.FailureType, res2.FailureCause)
        log.Println("2. Result:", res2.Result)
        log.Println("3. Error:", res3.FailureType, res3.FailureCause)
        log.Println("3. Result:", res3.Result)

        // now that all the tasks are done, complete the workflow
        w.Complete("Just amazing")
    }

The SwfWorkflow object uses the function name 'Go' to launch an activity. It will
return a channel. This channel will be closed if the result (success or failure) has
not yet been reached on the task. You should never block in the decider, unless you
want no decisions (Eg. 'Go', 'Completed', 'StartTimer') to be posted. Decision tasks
are not sticky, so if you have multiple deciders running any of them can handle the next
decision for a single workflow. Ordering of calls to 'Go' must be deterministic, and be
exactly the same, no matter the order of response from SWF. This is so we can properly
identify which response matches with which call.



Activities
==========
Activities do not need to be determintic. They are execute once per task,
and return their a response to the decider.

    // the worker here will listen to the domain, region and task list for activity tasks
    act := &swfhelper.ActivityWorker{
        Region: "us-west-2",
        Domain: "yewwdev",
        TaskList: "actlist",
        Identity: "atHome1-worker",
        MaxWorkers: 1,
    }
    // we register the function 'ActivityOne' to handle the 'getnum' swf activity.
    act.RegisterActivity(swf.ActivityType{Name:"getnum", Version:"1"}, ActivityOne)
    // block polling for new activity tasks
    act.Start()

The actual activity is defined as such:

    // swfhelper.ActivityContext provides us with the input to the activity, and ways
    // to respond to it.
    func ActivityOne(a *swfhelper.ActivityContext) {
        myRes := int64(rand.Int())
        // we defer the completed decision, as we must make a decision at the end
        defer a.Completed(strconv.FormatInt(myRes, 10))
        defer log.Println("Hey:", a.Input, myRes)
        time.Sleep(5 * time.Second + time.Second * time.Duration(rand.Float32() * 5))
    }

Activities and deciders do not need to run on the same machine, or executable.
*/
package swfhelper


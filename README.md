# Zenaton examples for Ruby
This repository contains examples of workflows built with Zenaton. These examples illustrates how Zenaton orchestrates tasks that are executed on different workers.

## Installation
Download this repo
```
git clone https://github.com/zenaton/examples-go.git; cd examples-go
```
then add an .env file
```
cp .env.example .env
```
and populate it with your application id and api token found [here](https://zenaton.com/app/api).

### Running Locally
Install dependencies
```
go get github.com/zenaton/zenaton-go
```
Then, you need to install a Zenaton worker
```
curl https://install.zenaton.com | sh
```
and start it:
```
zenaton start
```
and make it listen to your configuration:
```
zenaton listen --boot=boot/boot.go
```
Your all set!


*Your workflows will be processed by your worker, so you won't see anything except the stdout and stderr, respectively `zenaton.out` and `zenaton.err`. Look at these files :)*

## Example 1 : Sequential tasks execution
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/sequential.go) showcases
- a sequential execution of two tasks. The second task is executed only when the first one is processed.
- In a sequential task execution, you can get the output of a task. The result of the first task can be used by the second one.

<p align="center">
    <img src="support/sequential_workflow.png" alt="Sequential Workflow Diagram" />
</p>

```
go run sequential/main.go
```

## Example 2: Parallel tasks execution
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/parallel.go) showcases
- a parallel execution of 2 tasks
- a third task that is executed only after *both* first two tasks were processed

<p align="center">
    <img src="support/parallel_workflow.png" alt="Parallel Workflow Diagram" />
</p>

```go
go launch_parallel.go
```

## Example 3: Asynchronous tasks execution
[this example](https://github.com/zenaton/examples-go/tree/master/workflows/asynchronous.go) showcases
- An asynchronous execution of a task A (fire and forget)
- Then a sequential execution of Task B

<p align="center">
    <img src="support/asynchronous_workflow.png" alt="Asynchronous Workflow Diagram" />
</p>

```go
go launch_asynchronous.go
```

When a task is dispatched asynchronously, the workflow continues its execution without waiting for the task completion. Consequently, a task asynchronous dispatching always returns a null value.

## Example 4: Event
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/event.go) showcases
- how to change a workflow's behaviour based on an external event

<p align="center">
    <img src="support/event_workflow.png" alt="Event Workflow Diagram" />
</p>

```go
go launch_event.go
```

## Example 5: Wait
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/wait.go) showcases
- how the provided `Wait` task can be used to pause the workflow for a specified duration

<p align="center">
    <img src="support/wait_workflow.png" alt="Wait Workflow Diagram" />
</p>

```go
go launch_wait.go
```

## Example 6: Wait Event
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/wait_event.go) showcases
- how the provided `Wait` task can also be used to pause the workflow up to receiving a specific external event

<p align="center">
    <img src="support/waitEvent_workflow.png" alt="WaitEvent Workflow Diagram" />
</p>

```go
go launch_wait_event.go
```

## Example 7: Recursive Workflow
[This example](https://github.com/zenaton/examples-go/tree/master/recursive/recursive.go) showcases
- how launching events or workflows directly from orchestrated tasks allows you to schedule recurring workflows

```go
go launch_recursive.go
```

## Example 8: Workflow Versions
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/version.go) showcases
- how to update your workflow implementation, even while previous versions are still running

```go
go launch_version.go
```


Todo:

Necessary:
    in zenaton listen I need to write the paths differently depending on if I'm in the gopath or not :/ meh
    OnFailure and OnTimout in Worker/V1/workflow.go
    what is the purpose of CheckClasses (ask Gilles or Igor)
    figure out why zenaton unlisten --go doesn't work (ask Luis)
    loader.checkdir
        should it return an error?
        Should it check that the directory is accessible?
            	//try {
            	//fs.accessSync(dir, fs.constants.R_OK);
            	//} catch (e) {
            	//this.failure('Can not read \'' + dir + '\' directory');
            	//}
    fix zenaton listen todos
    make sure this is impossible to get index out of bounds in task.execute
    add MaxProcessingTime func() int64 (both workflows and tasks?)

    do we need to check that the data sent on an event is the same type as the data recieved?
        from Builder.send:
        //todo: do we want to have a different method for each type? or use this empty interface?
        //onEventType := reflect.TypeOf(b.workflow.OnEvent)
        //fmt.Println("onEventType.In(1): ", onEventType.In(1))
        //fmt.Println("reflect.TypeOf(eventData) ", reflect.TypeOf(eventData))
        //if onEventType.In(1) != reflect.TypeOf(eventData) {
        //panic("incompatible types")
        //}
    OnFailure func(*task.Task, error)
        OnTimeout func(*task.Task)
        (for both tasks and workfows?)
    figure out if client.SendEvent should return something.
    Complete Versioning
    Figure out if I need to have a .env (ask gilles)
    Does the boot file need to be in a package called "boot" as described in the example boot file?
    Add integration test to test for outputs
    Figure out errors
        should loader.getclasses return an error?
        in Worker.Process (agent library), send a SpecialError. Do I need the stack trace (ask Gilles)
        also in Worker.Process:
            //} catch (e) {
            //if (e instanceof ZenatonError) {
            //this.microserver.failWorker(e);
            //this.microserver.reset();
            //throw e;
            //}
            //
            //this.microserver.failWork(e);
            //this.microserver.reset();
            //throw e;
            //}
    Write documentation
    go fmt.
    clean up log statements
    test on windows (ask Gilles if this is important for this release)
    //todo check accessibility:
    	//try {
    	//fs.accessSync(dir, fs.constants.R_OK);
    	//} catch (e) {
    	//this.failure('Can not read \'' + dir + '\' directory');
    	//}


Can be done later:
    confirm that debug logs look good
    Clean up test_listen. Needs to be refactored, better naming of variables etc.
    Many methods should not be accessible to a user but must be because the libraries are split. For example, the workflow and task managers
        or GetTimestampOrDuration
    would be nice if the handle func could take many arguments, instead of just one. would have to think how that would be done (maybe pass in argments into execute?)
    Make it impossible to create a workflow or task, without using the constructor method for safety
    write documentation in code
    Refactor interfaces directory
    take custom serialization out and only leave json for now
    should have the equivilent documentation of this from cadence:
        This sample workflow demonstrates how to use multiple Cadence corotinues (instead of native goroutine) to process a
        chunk of a large work item in parallel, and then merge the intermediate result to generate the final result.
        In cadence workflow, you should not use go routine. Instead, you use corotinue via workflow.Go method.

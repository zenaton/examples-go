> ⚠️ This repository is abandoned.

# Zenaton examples for Go
This repository contains examples of workflows built with Zenaton. These examples illustrates how Zenaton orchestrates tasks that are executed on different workers.

## Installation
Download this repo
```
go get github.com/zenaton/examples-go
```
and the client library
```
go get github.com/zenaton/zenaton-go/...
```
cd into example directory
```
cd $(go env GOPATH)/src/github.com/zenaton/examples-go
```
then add an .env file
```
cp -n .env.example .env
```
and populate it with your application id and api token found [here](https://zenaton.com/app/api).

### Running Locally
Then, you need to install a Zenaton worker
```
curl https://install.zenaton.com | sh
```
and make it listen to your configuration:
```
zenaton listen --boot=boot/boot.go
```
Your all set!


*Your workflows will be processed by your worker, so you won't see anything except the stdout and stderr, respectively `zenaton.out` and `zenaton.err`. Look at these files :)*

## Example 1 : Sequential tasks execution
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/sequential.go) showcases
- a sequential execution of three tasks. The second and third tasks are executed only when the previous one is processed.
- In a sequential task execution, you can get the output of a task. The result of a task can be used by the next one.

<p align="center">
    <img
        src="https://raw.githubusercontent.com/zenaton/resources/master/examples/images/png/flow_sequential.png"
        alt="Sequential Workflow Diagram"
        width="400px"
    />
</p>

```
go run sequential/main.go
```

## Example 2: Parallel tasks execution
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/parallel.go) showcases
- a parallel execution of 2 tasks
- a third task that is executed only after *both* first two tasks were processed

<p align="center">
    <img
        src="https://raw.githubusercontent.com/zenaton/resources/master/examples/images/png/flow_parallel.png"
        alt="Parallel Workflow Diagram"
        width="400px"
    />
</p>

```
go run parallel/main.go
```

## Example 3: Asynchronous tasks execution
[this example](https://github.com/zenaton/examples-go/tree/master/workflows/asynchronous.go) showcases
- Asynchronous executions of Task A and Task B (fire and forget)
- Then sequential executions of Task C and Task D

<p align="center">
    <img
        src="https://raw.githubusercontent.com/zenaton/resources/master/examples/images/png/flow_async.png"
        alt="Asynchronous Workflow Diagram"
        width="400px"
    />
</p>

```
go run asynchronous/main.go
```

When a task is dispatched asynchronously, the workflow continues its execution without waiting for the task completion. Consequently, a task asynchronous dispatching does not return a value from the task.

## Example 4: Event
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/event.go) showcases
- how to change a workflow's behaviour based on an external event

<p align="center">
    <img
        src="https://raw.githubusercontent.com/zenaton/resources/master/examples/images/png/flow_react_event.png"
        alt="Event Workflow Diagram"
        width="400px"
    />
</p>

```
go run event/main.go
```

## Example 5: Wait
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/wait.go) showcases
- how the provided `Wait` task can be used to pause the workflow for a specified duration

<p align="center">
    <img
        src="https://raw.githubusercontent.com/zenaton/resources/master/examples/images/png/flow_wait.png"
        alt="Wait Workflow Diagram"
        width="400px"
    />
</p>

```
go run wait/main.go
```

## Example 6: Wait Event
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/wait_event.go) showcases
- how the provided `Wait` task can also be used to pause the workflow up to receiving a specific external event

<p align="center">
    <img
        src="https://raw.githubusercontent.com/zenaton/resources/master/examples/images/png/flow_wait_event.png"
        alt="WaitEvent Workflow Diagram"
        width="400px"
    />
</p>

```
go run waitevent/main.go
```

## Example 7: Error Workflow
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/error.go) showcases
- how to recover from a faulty task
 <p align="center">
    <img
        src="https://raw.githubusercontent.com/zenaton/resources/master/examples/images/png/flow_error.png"
        alt="Error Workflow Diagram"
        width="400px"
    />
</p>

```
go run error/main.go
```


## Example 7: Recursive Workflow
[This example](https://github.com/zenaton/examples-go/tree/master/recursive/recursive.go) showcases
- how launching events or workflows directly from orchestrated tasks allows you to schedule recurring workflows

```
go run recursive/main.go
```

## Example 8: Workflow Versions
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/version.go) showcases
- how to update your workflow implementation, even while previous versions are still running

```
go run version/main.go
```

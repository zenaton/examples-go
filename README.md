# Zenaton examples for Go
This repository contains examples of workflows built with Zenaton. These examples illustrates how Zenaton orchestrates tasks that are executed on different workers.

## Installation
Download this repo
```
go get github.com/zenaton/examples-go/...
```
and the client library
```
go get github.com/zenaton/zenaton-go/...
```
cd into example directory
```
cd $(go env GOPATH)/src/github.com/zenaton/examples-go
```
then add a]n .env file
```
cp client/.example.env client/.env
```
and populate it with your application id and api token found [here](https://zenaton.com/app/api).
git pu
### Running Locally
Then, you need to install a Zenaton worker
```
curl https://install.zenaton.com | sh
```
and make it listen to your configuration:
```
zenaton listen --env=client/.env --boot=boot/boot.go
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

```
go run parallel/main.go
```

## Example 3: Asynchronous tasks execution
[this example](https://github.com/zenaton/examples-go/tree/master/workflows/asynchronous.go) showcases
- An asynchronous execution of a task A (fire and forget)
- Then a sequential execution of Task B

<p align="center">
    <img src="support/asynchronous_workflow.png" alt="Asynchronous Workflow Diagram" />
</p>

```
go run asynchronous/main.go
```

When a task is dispatched asynchronously, the workflow continues its execution without waiting for the task completion. Consequently, a task asynchronous dispatching does not return a value from the task.

## Example 4: Event
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/event.go) showcases
- how to change a workflow's behaviour based on an external event

<p align="center">
    <img src="support/event_workflow.png" alt="Event Workflow Diagram" />
</p>

```
go run event/main.go
```

## Example 5: Wait
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/wait.go) showcases
- how the provided `Wait` task can be used to pause the workflow for a specified duration

<p align="center">
    <img src="support/wait_workflow.png" alt="Wait Workflow Diagram" />
</p>

```
go run wait/main.go
```

## Example 6: Wait Event
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/wait_event.go) showcases
- how the provided `Wait` task can also be used to pause the workflow up to receiving a specific external event

<p align="center">
    <img src="support/waitEvent_workflow.png" alt="WaitEvent Workflow Diagram" />
</p>

```
go run waitevent/main.go
```

## Example 7: Recursive Workflow
[This example](https://github.com/zenaton/examples-go/tree/master/recursive/recursive.go) showcases
- how launching events or workflows directly from orchestrated tasks allows you to schedule recurring workflows

```
go recursive/main.go
```

## Example 8: Workflow Versions
[This example](https://github.com/zenaton/examples-go/tree/master/workflows/version.go) showcases
- how to update your workflow implementation, even while previous versions are still running

```
go version/main.go
```

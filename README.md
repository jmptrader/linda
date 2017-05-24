# Linda

Linda provides a unified API across different broker (queue) services.

Brokers allow you to defer the processing of a time consuming task.

Inspiration comes from [beanstalkd](https://github.com/kr/beanstalkd) and [goworker](https://github.com/benmanns/goworker) 

## Installation

To install Linda, use
```sh
go get github.com/amlun/linda
```
to get the main package, and then use [glide](https://glide.sh/)
```sh
glide install
```
to install the dependencies

## Getting Started

### Terminology

* Broker
> message transport [MQ]

* poller
> poll job from the broker and send to local job channels
> poller also migrate the expire jobs

* worker
> worker is the main process to work the job

### Worker Type

```
func(job *Job) error
```

### Register Worker
```
linda.Register("MyClass", myFunc)
```

### Examples

```go
package main

import (
	"fmt"
	"github.com/amlun/linda"
)

func init() {
	settings := linda.Settings{
		Queue:         "scheduler",
		Connection:    "redis://localhost:6379/",
		Timeout:       60,
		IntervalFloat: 5.0,
		Concurrency:   1,
		Ack:           true,
	}
	linda.SetSettings(settings)
}
func main() {
	if err := linda.Init(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	broker := linda.GetBroker()
	queue := "scheduler"
	job := &linda.Job{
		Queue: queue,
		Payload: linda.Payload{
			Class: "DispatcherSeed",
			Args:  []interface{}{"seed_url_md5"},
		},
	}

	if err := broker.Push(job, queue); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
```

```go
package main

import (
	"fmt"
	"github.com/amlun/linda"
)

func init() {
	settings := linda.Settings{
		Queue:         "scheduler",
		Connection:      "redis://localhost:6379/",
		Timeout:       60,
		IntervalFloat: 5.0,
		Concurrency:   1,
		Ack:           true,
	}
	linda.SetSettings(settings)
	linda.RegisterWorkers("DispatcherSeed", DispatcherSeed)
}

func main() {
	if err := linda.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func DispatcherSeed(job *linda.Job) error {
	broker := linda.GetBroker()
	// get seed info
	// do seed job
	// delete reserved job and release
	broker.DeleteAndRelease("scheduler", job, 60)
	return nil
}
```

## Features

### Broker List

 - [x] Redis
 - [ ] beanstalkd
 - [ ] NSQ
 - [ ] Kafka
 - [ ] RabbitMQ
 
## Design

![system-design](https://rawgit.com/amlun/linda/master/images/linda.png)
 
## Thanks

* [redigo](https://github.com/garyburd/redigo)
* [goworker](https://github.com/benmanns/goworker)
* [laravel/queue](https://github.com/laravel/framework)
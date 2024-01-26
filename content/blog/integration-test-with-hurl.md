+++
title = "Integration tests with Hurl"
date = "2023-12-23"
draft = true
path = "/2023/12/23/integration-tests-with-hurl"
lang = "en"
[extra]
tags = ["integration-test", "tests", "hurl", "pipeline"]
comment = true
+++
Recently, we've involved in a testing situations which some services didn't have enough tests and we needed to add some tests
on top of services to cover those gaps as we were in hurry to fix the issue on production and also preparing a demo.
[Hurl](https://hurl.dev/) was the option that when I found it in the first place, I was thinking wow, no need any bash script in order to test something!
<!-- more -->

When we're talking about the [integration tests](https://en.wikipedia.org/wiki/Integration_testing), it means we are going to tests some components/services which are working together.
Let's refer to the [Wikipedia](https://en.wikipedia.org/wiki/Integration_testing) about it:
{% quote(type="info") %}
Integration testing (sometimes called integration and testing, abbreviated I&T) is the phase in software testing in which the whole
software module is tested or if it consists of multiple software modules they are combined and then tested as a group.
Integration testing is conducted to evaluate the compliance of a system or component with specified functional requirements.
It occurs after unit testing and before system testing. Integration testing takes as its input modules that have been unit tested,
groups them in larger aggregates, applies tests defined in an integration test plan to those aggregates, and delivers as its output
the integrated system ready for system testing.
{% end %}

Now, we know the meaning of Integration tests, so, it's time to define a scenario and see how can we use Hurl to add the integration tests.

Let's imagin we have 3 services which are:
* Worker (background job)
* API

```
 +--------------+                                 +--------------+
 |     API      |<---------send/receive---------->|    Worker    |
 +--------------+                                 +--------------+
```

So, the API gets the request, sends it to the Worker and worker assigns a process-id on the request and triggers the background process,
and returns the process-id to the API and API gives it to the client. the background process will start calling the Account Management service
in order to create and register a new account for the request which may takes 5 minutes, let's imagine.

In this article, all codes are in [Go](https://go.dev/).

{% quote(type="warning") %}
Go codes are very basic just to put you in the way to use Hurl, therefore, don't expect a perfect Go code in this article.
{% end %}

### API
Let's create a sample HTTP service to serves a request:

```go

```

### Worker

The worker, receives request and what it does bascially is it will store the request in a queue and return a process ID for
future investigation, and when the request is completed (success or failed), then it will call the final webhook if it was defined
in the request.

Input structure:
_/jobs endpoint expected request body_:

```go
type NewJobRequest struct {
	Command     string `json:"command"      binding:"required"`
	CallbackURL string `json:"callback_url" binding:"url"`
}
```

Worker logic:
```go,linenos
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// JobStatus defines data type for job statuses.
type JobStatus string

// List of defined job statues.
const (
	Pending JobStatus = "pending"
	Running           = "running"
	Failed            = "failed"
	Succeed           = "succeed"
)

// Job represents job data structure.
type Job struct {
	Command     string
	CallbackURL string
	Status      JobStatus
}

// NewJobRequest represents user input data structure to create a new job.
type NewJobRequest struct {
	Command     string `json:"command"      binding:"required"`
	CallbackURL string `json:"callback_url"`
}

// GetJobRequest represents user input data structure to get a job info.
type GetJobRequest struct {
	ID string `uri:"id" binding:"required"`
}

func main() {
	ctx := context.Background()
	r := gin.Default()
	jobs := make(map[string]Job)

	r.POST("/jobs", func(ctx *gin.Context) {
		// Parse the request body.
		var input NewJobRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"attributes": gin.H{
					"message": err,
				},
			})
			return
		}

		// Create a new process ID and store the jobs.
		processID := uuid.NewString()
		jobs[processID] = Job{
			Command:     input.Command,
			CallbackURL: input.CallbackURL,
			Status:      Pending,
		}

		// Return the processID as a response.
		ctx.JSON(http.StatusAccepted, gin.H{
			"status": "ok",
			"attributes": gin.H{
				"process_id": processID,
			},
		})
	})

	r.GET("/jobs/:id", func(ctx *gin.Context) {
		// Parse the request URI.
		var input GetJobRequest
		if err := ctx.ShouldBindUri(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "erorr",
				"attributes": gin.H{
					"message": err,
				},
			})
			return
		}

		// Try to find the job.
		job, ok := jobs[input.ID]
		if !ok {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"attributes": gin.H{
					"message": "not found",
				},
			})
			return
		}

		// Return the job info.
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"attributes": gin.H{
				"id":     input.ID,
				"status": job.Status,
			},
		})
	})

	// Background service to run jobs.
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				log.Println("received stop signal")
				return
			default:
				for processID, job := range jobs {
					log.Println(processID, job.Command)

					// Just a short nap in each run.
					time.Sleep(1 * time.Second)
				}
			}
		}
	}(ctx)

	if err := r.Run(); err != nil {
		log.Fatal("running service failed", err)
	}
}
```


### Test

Now, it's time to write an integration tests.
Let's use Hurl functioanlities to test our systems.

First we need to test new command endpoint, and then testing the status of that command.

_cmd.hurl_:
```bash
POST http://localhost:8080/actions
{
    "command": "sleep 10; ls -lha"
}

HTTP 201
[Asserts]
jsonpath "$.id" exists
[Captures]
action_id: jsonpath "$.id"
```

Now, we need to keep trying to get the action status:

```bash
GET http://localhost:8080/actions/{{action_id}}

HTTP 200
[Asserts]
jsonpath "$.status" == "finished"
[Options]
retry: 10
retry-interval: 60000
```

Here we added "Options" to configure our requests policy.

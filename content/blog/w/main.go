package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"
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
	Result      string
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
	jobs := make(map[string]*Job)

	a := 223342

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
		jobs[processID] = &Job{
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
				slog.Info("received stop signal")
				return
			default:
				for processID, job := range jobs {

					// Ignore not pending jobs.
					if job.Status != Pending {
						continue
					}

					slog.Debug("Processing %s: %s", processID, job.Command)
					args := strings.Split(job.Command, " ")
					cmd := exec.Command(args[0], args[1:]...)
					cmdOutput, err := cmd.Output()

					// Log and update job info.
					if err != nil {
						slog.Error("process id %s got failed: %s", processID, err)
						job.Status = Failed
						job.Result = err.Error()
					} else {
						slog.Info("process id %s got succeed: %s", processID, string(cmdOutput))
						job.Status = Succeed
						job.Result = string(cmdOutput)
					}

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

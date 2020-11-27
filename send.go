package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type TaskInfoEvent string

const (
	UpdateTaskInfoEvent        TaskInfoEvent = "UPDATE"
	DeleteTaskInfoEvent        TaskInfoEvent = "DELETE"
	UpdateStreamTitleInfoEvent TaskInfoEvent = "STREAM_TITLE_UPDATE"
)

//StreamTitleInfo dto.
type StreamTitleInfo struct {
	OldStreamTitle string `json:"old_stream_title"`
	NewStreamTitle string `json:"new_stream_title"`
}

//TaskInfo dto.
type TaskInfo struct {
	OrgID           string          `json:"org_id" bson:"org_id"`
	WorkspaceID     string          `json:"workspace_id"`
	TaskID          string          `json:"task_id"`
	StreamTitle     string          `json:"stream_title"`
	TaskName        string          `json:"task_name"`
	TaskStatus      string          `json:"task_status"`
	TaskPriority    int64           `json:"task_priority"`
	TaskDueDate     *time.Time      `json:"task_due_date"`
	StreamTitleInfo StreamTitleInfo `json:"stream_title_info"`
	Event           TaskInfoEvent   `json:"event"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"timesheet-events", // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// task := TaskInfo{
	// 	OrgID:        "332",
	// 	WorkspaceID:  "1246",
	// 	TaskID:       "",
	// 	StreamTitle:  "st1",
	// 	TaskName:     "fix test",
	// 	TaskStatus:   "In Progress",
	// 	TaskPriority: 1,
	// 	TaskDueDate:  nil,
	// 	StreamTitleInfo: StreamTitleInfo{
	// 		NewStreamTitle: "zaz update by me",
	// 		OldStreamTitle: "zaz",
	// 	},
	// 	Event: "STREAM_TITLE_UPDATE",
	// }

	task := TaskInfo{
		OrgID:       "332",
		WorkspaceID: "1246",
		StreamTitleInfo: StreamTitleInfo{
			NewStreamTitle: "zaz update by me",
			OldStreamTitle: "zaz",
		},
		Event: "STREAM_TITLE_UPDATE",
	}

	body, err := json.Marshal(&task)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

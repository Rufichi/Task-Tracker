// task.go
package main

import (
    "time"
)

type Task struct {
    ID          int       `json:"id"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

type TaskManager struct {
    Tasks    []Task `json:"tasks"`
    NextID   int    `json:"nextId"`
    filename string
}

func NewTaskManager() *TaskManager {
    return &TaskManager{
        Tasks:    []Task{},
        NextID:   1,
        filename: "tasks.json",
    }
}

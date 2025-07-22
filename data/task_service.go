package data

import (
	"errors"
	"sync"
	"taskmanager/models"
)

// In-memory database
var (
	tasks  = make(map[int]models.Task)
	nextId = 1
	mutex  = &sync.Mutex{} // handle concurrent requests
)

// GetAllTasks - retrieve all tasks
func GetAllTasks() []models.Task {
	mutex.Lock()
	defer mutex.Unlock()

	allTasks := make([]models.Task, 0, len(tasks))
	for _, task := range tasks {
		allTasks = append(allTasks, task)
	}
	return allTasks
}

// GetTaskById - retrieve a task by its id
func GetTaskById(id int) (models.Task, error) {
	mutex.Lock()
	defer mutex.Unlock()

	task, exists := tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}
	return task, nil
}

// CreateTask - create a new task
func CreateTask(task models.Task) models.Task {
	mutex.Lock()
	defer mutex.Unlock()

	task.Id = nextId
	nextId++
	tasks[task.Id] = task
	return task
}

// UpdateTask - update an existing task
func UpdateTask(id int, updatedTask models.Task) (models.Task, error) {
	mutex.Lock()
	defer mutex.Unlock()

	_, exists := tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}
	updatedTask.Id = id
	tasks[id] = updatedTask
	return updatedTask, nil
}

// DeleteTask - removes a task by its id
func DeleteTask(id int) error {
	mutex.Lock()
	defer mutex.Unlock()

	_, exists := tasks[id]
	if !exists {
		return errors.New("task not found")
	}

	delete(tasks, id)
	return nil
}

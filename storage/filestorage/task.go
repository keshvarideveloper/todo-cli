package filestorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"todo-cli/constance"
	"todo-cli/entity"
)

type TaskStorage struct {
	filePath          string
	serializationMode string
}

var taskStorageInMemory []entity.Task

func NewTaskStorage(filePath, serializationMode string) *TaskStorage {
	return &TaskStorage{
		filePath:          filePath,
		serializationMode: serializationMode,
	}
}

func (ts *TaskStorage) Write(task entity.Task) error {
	file, err := os.OpenFile(ts.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Can't create or open the file", ts.filePath)

		return err
	}

	err = ts.loadTaskToMemory()
	if err != nil {

		return err
	}

	taskId := len(taskStorageInMemory) + 1
	task.ID = taskId

	var data []byte

	switch ts.serializationMode {
	case constance.NormalSerializationMode:
		data = []byte(fmt.Sprintf("id: %d, title: %s, CategoryID: %d, UserID: %d, DouDate: %s, isDone: %v\n",
			task.ID, task.Title, task.CategoryID, task.UserID, task.DouDate, task.GetStatus()))
	case constance.JsonSerializationMode:
		var jErr error
		data, jErr = json.Marshal(task)
		if jErr != nil {
			fmt.Println("can't marshal task struct to json", jErr)

			return jErr
		}
		data = append(data, '\n')
	}

	var wErr error
	_, wErr = file.Write(data)
	if wErr != nil {
		fmt.Printf("can't write to the file %v\n", wErr)

		return wErr
	}

	return nil

}
func (ts *TaskStorage) GetById(ID int) (entity.Task, error) {
	err := ts.loadTaskToMemory()
	if err != nil {

		return entity.Task{}, err
	}

	for _, task := range taskStorageInMemory {
		if task.ID == ID {

			return task, nil
		}
	}

	return entity.Task{}, errors.New("Task not found")
}

func (ts *TaskStorage) List() ([]entity.Task, error) {
	err := ts.loadTaskToMemory()
	if err != nil {

		return []entity.Task{}, err
	}

	return taskStorageInMemory, nil
}
func (ts *TaskStorage) loadTaskToMemory() error {
	dataStr, err := os.ReadFile(ts.filePath)
	if err != nil {
		fmt.Printf("File << %s >> can't open %v\n", ts.filePath, err)

		return err
	}

	taskSlice := strings.Split(string(dataStr), "\n")

	var taskStorage []entity.Task

	for _, taskStr := range taskSlice {
		var taskStruct entity.Task
		switch ts.serializationMode {
		case constance.JsonSerializationMode:
			if taskStr == "" {
				continue
			}

			jErr := json.Unmarshal([]byte(taskStr), &taskStruct)
			if jErr != nil {
				fmt.Println("The record can't convert to Task struct", jErr)

				continue
			}
		case constance.NormalSerializationMode:
			var dErr error
			taskStruct, dErr = entity.DeserializeTaskFromNormal(taskStr)

			if dErr != nil {
				fmt.Println("The record can't convert to Task struct", dErr)

				continue
			}
		default:
			fmt.Println("Serialization mode is not valid")

			return err
		}
		taskStorage = append(taskStorage, taskStruct)
	}
	taskStorageInMemory = taskStorage
	return nil
}

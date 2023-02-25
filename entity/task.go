package entity

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Task struct {
	ID, CategoryID, UserID int
	Title, DouDate         string
	isDone                 bool
}

func (t Task) GetStatus() bool {
	return t.isDone
}
func DeserializeTaskFromNormal(taskStr string) (Task, error) {
	if taskStr == "" || strings.ContainsRune(taskStr, '{') {
		return Task{}, errors.New("Task string can't deserialize to normal")
	}

	var task Task
	taskFields := strings.Split(taskStr, ",")
	for _, field := range taskFields {
		values := strings.Split(field, ": ")
		if len(values) != 2 {
			fmt.Println("Field is not valid, skipping...", len(values))

			continue
		}

		fieldName := strings.ReplaceAll(values[0], " ", "")
		fieldValue := values[1]

		fieldName = strings.ToLower(fieldName)
		switch fieldName {
		case "id":
			id, err := strconv.Atoi(fieldValue)
			if err != nil {

				return Task{}, errors.New("strconv error")
			}
			task.ID = id
		case "title":
			task.Title = fieldValue
		case "categoryid":
			categoryId, err := strconv.Atoi(fieldValue)
			if err != nil {

				return Task{}, errors.New("strconv error")
			}
			task.CategoryID = categoryId
		case "userid":
			userId, err := strconv.Atoi(fieldValue)
			if err != nil {

				return Task{}, errors.New("strconv error")
			}
			task.UserID = userId
		case "isdone":
			isDone, err := strconv.ParseBool(fieldValue)
			if err != nil {

				return Task{}, errors.New("strconv error")
			}
			task.isDone = isDone
		case "doudate":
			task.DouDate = fieldValue
		}

	}
	return task, nil
}

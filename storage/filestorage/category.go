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

type CategoryStorage struct {
	filePath          string
	serializationMode string
}

var categoryStorageInMemory []entity.Category

func NewCategoryStorage(filePath, serializationMode string) *CategoryStorage {
	return &CategoryStorage{
		filePath:          filePath,
		serializationMode: serializationMode,
	}
}

func (cs *CategoryStorage) Write(category entity.Category) error {

	file, err := os.OpenFile(cs.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Can't create or open the file", cs.filePath)

		return err
	}
	err = cs.loadCategoryToMemory()
	if err != nil {

		return err
	}

	categoryId := len(categoryStorageInMemory) + 1
	category.ID = categoryId

	var data []byte

	switch cs.serializationMode {
	case constance.NormalSerializationMode:
		data = []byte(fmt.Sprintf("id: %d, title: %s, color: %s, userId: %d\n", category.ID, category.Title, category.Color, category.UserID))
	case constance.JsonSerializationMode:
		var jErr error
		data, jErr = json.Marshal(category)
		if jErr != nil {
			fmt.Println("can't marshal category struct to json", jErr)

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

func (cs *CategoryStorage) GetById(ID int) (entity.Category, error) {
	err := cs.loadCategoryToMemory()
	if err != nil {
		return entity.Category{}, err
	}

	for _, category := range categoryStorageInMemory {
		if category.ID == ID {

			return category, nil
		}
	}

	return entity.Category{}, errors.New("Category not found")
}

func (cs *CategoryStorage) List() ([]entity.Category, error) {
	err := cs.loadCategoryToMemory()
	if err != nil {

		return []entity.Category{}, err
	}

	return categoryStorageInMemory, nil

}

func (cs *CategoryStorage) loadCategoryToMemory() error {
	dataStr, err := os.ReadFile(cs.filePath)
	if err != nil {
		fmt.Printf("File << %s >> can't open %v\n", cs.filePath, err)

		return err
	}

	categorySlice := strings.Split(string(dataStr), "\n")

	var categoryStorage []entity.Category

	for _, categoryStr := range categorySlice {
		var categoryStruct entity.Category
		switch cs.serializationMode {
		case constance.JsonSerializationMode:
			if categoryStr == "" {
				continue
			}

			jErr := json.Unmarshal([]byte(categoryStr), &categoryStruct)

			if jErr != nil {
				fmt.Println("The record can't convert to categoryt struct", jErr)

				continue
			}
		case constance.NormalSerializationMode:
			var dErr error
			categoryStruct, dErr = entity.DeserializeCategoryFromNormal(categoryStr)

			if dErr != nil {
				fmt.Println("The record can't convert to category struct", dErr)

				continue
			}
		default:
			fmt.Println("Serialization mode is not valid")

			return err
		}
		categoryStorage = append(categoryStorage, categoryStruct)
	}
	categoryStorageInMemory = categoryStorage
	return nil
}

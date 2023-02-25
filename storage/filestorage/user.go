package filestorage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"todo-cli/constance"
	"todo-cli/entity"
)

type UserFileStorage struct {
	filePath          string
	serializationMode string
}

var userStorageInMemory []entity.User

func NewUserFileStorage(filePath, serializationMode string) *UserFileStorage {
	return &UserFileStorage{
		filePath:          filePath,
		serializationMode: serializationMode,
	}
}

func (f *UserFileStorage) Write(user entity.User) error {
	file, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("can't create or open the file")

		return err
	}

	err = f.loadUserToMemory()
	if err != nil {

		return err
	}

	userId := len(userStorageInMemory) + 1
	user.ID = userId

	var data []byte

	switch f.serializationMode {
	case constance.NormalSerializationMode:
		data = []byte(fmt.Sprintf("id: %d, email: %s, password: %s\n", user.ID, user.Email, user.Password))
	case constance.JsonSerializationMode:
		var jErr error
		data, jErr = json.Marshal(user)
		if jErr != nil {
			fmt.Println("can't marshal user struct to json", jErr)

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
func (f *UserFileStorage) GetByEmail(email string) (entity.User, error) {
	err := f.loadUserToMemory()
	if err != nil {

		return entity.User{}, err
	}

	for _, user := range userStorageInMemory {
		if user.Email == email {

			return user, nil
		}
	}

	return entity.User{}, nil
}

func (f *UserFileStorage) List() ([]entity.User, error) {
	err := f.loadUserToMemory()
	if err != nil {

		return []entity.User{}, err
	}

	return userStorageInMemory, nil
}

func (f *UserFileStorage) loadUserToMemory() error {
	dataStr, err := os.ReadFile(f.filePath)

	if err != nil {
		fmt.Printf("File << %s >> can't open %v\n", f.filePath, err)

		return err
	}

	userSlice := strings.Split(string(dataStr), "\n")

	var userStorage []entity.User
	for _, userStr := range userSlice {
		var userStruct entity.User
		switch f.serializationMode {
		case constance.JsonSerializationMode:
			if userStr == "" {
				continue
			}

			jErr := json.Unmarshal([]byte(userStr), &userStruct)

			if jErr != nil {
				fmt.Println("user", userStr)
				fmt.Println("The record can't convert to user struct", jErr)

				continue
			}
		case constance.NormalSerializationMode:
			var dErr error
			userStruct, dErr = entity.DeserializeUserFromNormal(userStr)

			if dErr != nil {
				fmt.Println("The record can't convert to user struct", dErr)

				continue
			}
		default:
			fmt.Println("Serialization mode is not valid")

			return err
		}

		userStorage = append(userStorage, userStruct)
	}
	userStorageInMemory = userStorage
	return nil
}

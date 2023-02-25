package entity

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type User struct {
	ID              int
	Email, Password string
}

func DeserializeUserFromNormal(userStr string) (User, error) {
	if userStr == "" || strings.ContainsRune(userStr, '{') {
		return User{}, errors.New("User string can't deserialize to normal")
	}

	var user User
	userFields := strings.Split(userStr, ",")
	for _, field := range userFields {
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
				return User{}, errors.New("strconv error")
			}
			user.ID = id
		case "email":
			user.Email = fieldValue
		case "password":
			user.Password = fieldValue
		}

	}
	return user, nil
}

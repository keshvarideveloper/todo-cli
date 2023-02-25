package entity

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Category struct {
	ID, UserID   int
	Title, Color string
}

func DeserializeCategoryFromNormal(categoryStr string) (Category, error) {
	if categoryStr == "" || strings.ContainsRune(categoryStr, '{') {
		return Category{}, errors.New("Category string can't deserialize to normal")
	}

	var category Category
	categoryFields := strings.Split(categoryStr, ",")
	for _, field := range categoryFields {
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
				return Category{}, errors.New("strconv error")
			}
			category.ID = id
		case "title":
			category.Title = fieldValue
		case "color":
			category.Color = fieldValue
		case "userid":
			categoryID, err := strconv.Atoi(fieldValue)
			if err != nil {
				fmt.Printf("Category ID is not valid integer %v\n ", err)

				return Category{}, err
			}
			category.UserID = categoryID
		}

	}
	return category, nil
}

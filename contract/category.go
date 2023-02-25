package contract

import "todo-cli/entity"

type WriteCategory interface {
	Write(c entity.Category) error
}

type ReadCategory interface {
	GetById(ID int) (entity.Category, error)
	List() ([]entity.Category, error)
}

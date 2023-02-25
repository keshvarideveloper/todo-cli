package contract

import "todo-cli/entity"

type WriteTask interface {
	Write(c entity.Task) error
}

type ReadTask interface {
	GetById(ID int) (entity.Task, error)
	List() ([]entity.Task, error)
}

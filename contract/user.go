package contract

import "todo-cli/entity"

type WriteUser interface {
	Write(u entity.User) error
}

type ReadUser interface {
	GetByEmail(email string) (entity.User, error)
	List() ([]entity.User, error)
}

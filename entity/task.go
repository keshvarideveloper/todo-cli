package entity

type Task struct {
	ID, CategoryID, UserID int
	Title, DouDate         string
	isDone                 bool
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"todo-cli/text_color"
)

type User struct {
	ID              int
	Email, Password string
}

type Category struct {
	ID, UserID   int
	Title, Color string
}

type Task struct {
	ID, CategoryID, UserID int
	Title, DouDate         string
	isDone                 bool
}

var (
	UserStorage     []User
	CategoryStorage []Category
	TaskStorage     []Task

	authenticatedUser *User

	scanner = bufio.NewScanner(os.Stdin)
)

func main() {
	command := flag.String("command", "login-user", "command to run")
	flag.Parse()

	for {
		runCommand(*command)

		fmt.Println("Enter new command:")
		scanner.Scan()
		*command = scanner.Text()

	}
}

func runCommand(command string) {
	if command != "exit" && command != "register-user" && command != "login-user" && authenticatedUser == nil {
		UserStorageLen := len(UserStorage)
		switch {
		case UserStorageLen == 0:
			fmt.Println("You need to register ... ")
			registerUser()
			loginUser()
		case UserStorageLen > 0:
			fmt.Println("Please Log in ... ")
			loginUser()
		}

		if authenticatedUser == nil {

			return
		}
	}

	switch command {
	case "create-task":
		createTask()
	case "create-category":
		createCategory()
	case "register-user":
		registerUser()
	case "login-user":
		loginUser()
	case "task-list":
		taskList()
	case "exit":
		os.Exit(0)
	default:
		fmt.Printf("command << %s >> is not valid, try another command\n", command)

	}
}

func createTask() {
	fmt.Println(text_color.Green + "Create a task process ..." + text_color.Reset)

	var title, dueDate, category string

	fmt.Println("Enter task title:")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("Enter task category ID:")
	scanner.Scan()
	category = scanner.Text()

	categoryID, err := strconv.Atoi(category)
	if err != nil {
		fmt.Printf("Category ID is not valid integer %v\n ", err)

		return
	}
	isCategoryFind := false
	for _, c := range CategoryStorage {
		if c.UserID == categoryID && c.UserID == authenticatedUser.ID {
			isCategoryFind = true

			break
		}
	}
	if isCategoryFind != true {
		println("The category is not found")

		return
	}
	fmt.Println("Enter task due date :")
	scanner.Scan()
	dueDate = scanner.Text()

	t := Task{
		ID:         len(UserStorage) + 1,
		Title:      title,
		isDone:     false,
		CategoryID: categoryID,
		DouDate:    dueDate,
		UserID:     authenticatedUser.ID,
	}
	TaskStorage = append(TaskStorage, t)

	fmt.Printf("New task added \n %+v\n", TaskStorage)
}
func createCategory() {
	fmt.Println(text_color.Purple + "Create a category process ..." + text_color.Reset)

	var title, color string

	fmt.Println("Enter category title:")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("Enter category color:")
	scanner.Scan()
	color = scanner.Text()

	c := Category{
		ID:     len(CategoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}
	CategoryStorage = append(CategoryStorage, c)

	fmt.Printf("New category added \n %+v\n", CategoryStorage)
}
func taskList() {
	fmt.Println(text_color.Green + "Register process ..." + text_color.Reset)

	for _, task := range TaskStorage {
		if task.UserID == authenticatedUser.ID {
			fmt.Println(task)
		}
	}
}
func registerUser() {
	fmt.Println(text_color.Yellow + "Register process ..." + text_color.Reset)

	var email, password string

	fmt.Println("Enter your email address:")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("Enter your password")
	scanner.Scan()
	password = scanner.Text()

	u := User{
		ID:       len(UserStorage) + 1,
		Email:    email,
		Password: password,
	}
	UserStorage = append(UserStorage, u)

	fmt.Printf("New user added \n %+v\n", UserStorage)

}
func loginUser() {
	fmt.Println(text_color.Cyan + "Login process ..." + text_color.Reset)
	authenticatedUser = nil

	var email, password string

	fmt.Println("Enter email address:")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("Enter password:")
	scanner.Scan()
	password = scanner.Text()

	for i, user := range UserStorage {
		if user.Email == email && user.Password == password {
			authenticatedUser = &UserStorage[i]

			break
		}
	}

	if authenticatedUser == nil {
		fmt.Println("The email or password is wrong")
	} else {
		fmt.Printf("Welcome %s\n", authenticatedUser.Email)
	}

}

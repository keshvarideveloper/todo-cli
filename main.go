package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"todo-cli/constance"
	"todo-cli/contract"
	"todo-cli/entity"
	"todo-cli/storage/filestorage"
	"todo-cli/textcolor"
)

var (
	authenticatedUser   *entity.User
	scanner             = bufio.NewScanner(os.Stdin)
	serializationMode   string
	userFileStorage     *filestorage.UserStorage
	categoryFileStorage *filestorage.CategoryStorage
	taskFileStorage     *filestorage.TaskStorage
)

func main() {
	serializeMode := flag.String("serialize-mode", constance.JsonSerializationMode, "serializationmode to write data to file")
	command := flag.String("command", "login-user", "command to run")
	flag.Parse()

	switch *serializeMode {
	case constance.NormalSerializationMode:
		serializationMode = constance.NormalSerializationMode
	case constance.JsonSerializationMode:
		serializationMode = constance.JsonSerializationMode
	default:
		return
	}

	userFileStorage = filestorage.NewUserStorage("user.txt", serializationMode)
	categoryFileStorage = filestorage.NewCategoryStorage("category.txt", serializationMode)
	taskFileStorage = filestorage.NewTaskStorage("task.txt", serializationMode)

	for {
		runCommand(*command, userFileStorage)

		fmt.Println("Enter new command:")
		scanner.Scan()
		*command = scanner.Text()

	}
}

func runCommand(command string, r contract.ReadUser) {
	if command != "exit" && command != "register-user" && command != "login-user" && authenticatedUser == nil {
		usersList, _ := r.List()
		UserStorageLen := len(usersList)
		switch {
		case UserStorageLen == 0:
			fmt.Println("You need to register ... ")
			registerUser(userFileStorage)
			loginUser(userFileStorage)
		case UserStorageLen > 0:
			fmt.Println("Please Log in ... ")
			loginUser(userFileStorage)
		}

		if authenticatedUser == nil {

			return
		}
	}

	switch command {
	case "create-task":
		createTask(taskFileStorage, categoryFileStorage)
	case "create-category":
		createCategory(categoryFileStorage)
	case "register-user":
		registerUser(userFileStorage)
	case "login-user":
		loginUser(userFileStorage)
	case "list-user":
		userList(userFileStorage)
	case "list-category":
		categoryList(categoryFileStorage)
	case "list-task":
		taskList(taskFileStorage)
	case "exit":
		os.Exit(0)
	default:
		fmt.Printf("command << %s >> is not valid, try another command\n", command)

	}
}

func createTask(w contract.WriteTask, categoryStorage contract.ReadCategory) {
	fmt.Println(textcolor.Green + "Create a task process ..." + textcolor.Reset)

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

	categoryExist, err := categoryStorage.GetById(categoryID)
	if err != nil {
		fmt.Println("Category not found")

		return
	}

	if categoryExist.UserID != authenticatedUser.ID {
		fmt.Println("The category is not yours")

		return
	}

	fmt.Println("Enter task due date :")
	scanner.Scan()
	dueDate = scanner.Text()

	t := entity.Task{
		Title:      title,
		CategoryID: categoryID,
		DouDate:    dueDate,
		UserID:     authenticatedUser.ID,
	}

	err = w.Write(t)
	if err != nil {
		fmt.Println(err)
	}
}
func createCategory(w contract.WriteCategory) {
	fmt.Println(textcolor.Purple + "Create a category process ..." + textcolor.Reset)

	var title, color string

	fmt.Println("Enter category title:")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("Enter category color:")
	scanner.Scan()
	color = scanner.Text()

	c := entity.Category{
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}

	err := w.Write(c)
	if err != nil {
		fmt.Println(err)
	}
}

func taskList(r contract.ReadTask) {
	fmt.Println(textcolor.Green + "Register process ..." + textcolor.Reset)

	tasks, err := r.List()
	if err != nil {
		fmt.Println(err)
	}

	var userTasks []entity.Task
	for _, task := range tasks {
		if task.UserID == authenticatedUser.ID {
			userTasks = append(userTasks, task)
		}
	}
	fmt.Printf("Users list:\n %+v", userTasks)
}

func registerUser(w contract.WriteUser) {
	fmt.Println(textcolor.Yellow + "Register process ..." + textcolor.Reset)

	var email, password string

	fmt.Println("Enter your email address:")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("Enter your password")
	scanner.Scan()
	password = scanner.Text()

	u := entity.User{
		Email:    email,
		Password: password,
	}

	err := w.Write(u)
	if err != nil {
		fmt.Println(err)
	}

}
func loginUser(r contract.ReadUser) {
	fmt.Println(textcolor.Cyan + "Login process ..." + textcolor.Reset)
	authenticatedUser = nil
	var email, password string

	fmt.Println("Enter email address:")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("Enter password:")
	scanner.Scan()
	password = scanner.Text()

	if user, err := r.GetByEmail(email); err == nil {
		if user.Password == password {
			authenticatedUser = &user
		}
	}

	if authenticatedUser == nil {
		fmt.Println("The email or password is wrong")
	} else {
		fmt.Printf("Welcome %s\n", authenticatedUser.Email)
	}

}

func userList(r contract.ReadUser) {
	users, err := r.List()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Users list:\n %+v\n", users)
}

func categoryList(r contract.ReadCategory) {
	categories, err := r.List()
	if err != nil {
		fmt.Println(err)
	}
	var userCategories []entity.Category
	for _, category := range categories {
		if category.UserID == authenticatedUser.ID {
			userCategories = append(userCategories, category)
		}
	}

	fmt.Printf("Categories list:\n %+v\n", userCategories)
}

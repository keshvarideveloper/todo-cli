package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"todo-cli/constance"
	"todo-cli/contract"
	"todo-cli/entity"
	"todo-cli/storage/filestorage"
	"todo-cli/textcolor"
)

var (
	//CategoryStorage []entity.Category
	//TaskStorage     []entity.Task

	authenticatedUser *entity.User

	scanner = bufio.NewScanner(os.Stdin)

	serializationMode string
)

var userFileStorage *filestorage.UserFileStorage

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

	userFileStorage = filestorage.NewUserFileStorage("user.txt", serializationMode)

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
	//case "create-task":
	//	createTask()
	//case "create-category":
	//	createCategory()
	case "register-user":
		registerUser(userFileStorage)
	case "login-user":
		loginUser(userFileStorage)
	case "list-user":
		userList(userFileStorage)
	//case "task-list":
	//	taskList()
	case "exit":
		os.Exit(0)
	default:
		fmt.Printf("command << %s >> is not valid, try another command\n", command)

	}
}

//	func createTask() {
//		fmt.Println(textcolor.Green + "Create a task process ..." + textcolor.Reset)
//
//		var title, dueDate, category string
//
//		fmt.Println("Enter task title:")
//		scanner.Scan()
//		title = scanner.Text()
//
//		fmt.Println("Enter task category ID:")
//		scanner.Scan()
//		category = scanner.Text()
//
//		categoryID, err := strconv.Atoi(category)
//		if err != nil {
//			fmt.Printf("Category ID is not valid integer %v\n ", err)
//
//			return
//		}
//		isCategoryFind := false
//		for _, c := range CategoryStorage {
//			if c.UserID == categoryID && c.UserID == authenticatedUser.ID {
//				isCategoryFind = true
//
//				break
//			}
//		}
//		if isCategoryFind != true {
//			println("The category is not found")
//
//			return
//		}
//		fmt.Println("Enter task due date :")
//		scanner.Scan()
//		dueDate = scanner.Text()
//
//		t := Task{
//			ID:         len(UserStorage) + 1,
//			Title:      title,
//			isDone:     false,
//			CategoryID: categoryID,
//			DouDate:    dueDate,
//			UserID:     authenticatedUser.ID,
//		}
//		TaskStorage = append(TaskStorage, t)
//
//		fmt.Printf("New task added \n %+v\n", TaskStorage)
//	}
//
//	func createCategory() {
//		fmt.Println(textcolor.Purple + "Create a category process ..." + textcolor.Reset)
//
//		var title, color string
//
//		fmt.Println("Enter category title:")
//		scanner.Scan()
//		title = scanner.Text()
//
//		fmt.Println("Enter category color:")
//		scanner.Scan()
//		color = scanner.Text()
//
//		c := Category{
//			ID:     len(CategoryStorage) + 1,
//			Title:  title,
//			Color:  color,
//			UserID: authenticatedUser.ID,
//		}
//		CategoryStorage = append(CategoryStorage, c)
//
//		fmt.Printf("New category added \n %+v\n", CategoryStorage)
//	}
//
//	func taskList() {
//		fmt.Println(textcolor.Green + "Register process ..." + textcolor.Reset)
//
//		for _, task := range TaskStorage {
//			if task.UserID == authenticatedUser.ID {
//				fmt.Println(task)
//			}
//		}
//	}
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

	fmt.Printf("users list:\n %+v", users)
}

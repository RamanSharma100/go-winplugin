package example

import (
	"fmt"
)

type User struct {
	Name string
	Age  int
}

func Execute(a, b int) int {
	return a + b
}

func Version() string {
	return "1.0.0"
}

func Logger(message string) {
	fmt.Println(message)
}

func ProcessUser(
	user *User,
) {
	fmt.Println(
		"User:",
		user.Name,
		user.Age,
	)
}

func ProcessData(
	data interface{},
) {
	fmt.Println(data)
}

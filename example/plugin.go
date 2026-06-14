package example

import (
	"errors"
	"fmt"
)

type User struct {
	Name string
	Age  int
}

func Execute(a, b int) int {
	return AddInternal(a, b)
}

func Version() string {
	return "1.0.0"
}

func Logger(message string) {
	fmt.Println(message)
}

func ProcessUser(user *User) {
	fmt.Println("User:", user.Name, user.Age)
}

func ProcessData(data interface{}) {
	fmt.Println(data)
}

func CreateUser(name string, age int) (*User, error) {
	return &User{Name: name, Age: age}, nil
}

func CreateUserWithError(name string, age int) (*User, error) {
	if age < 0 {
		return nil, errors.New("age cannot be negative")
	}
	return &User{Name: name, Age: age}, nil
}

func IsPositive(n int) bool {
	return n > 0
}

func Divide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func ReadData() []byte {
	return []byte("hello bytes")
}

func Validate(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

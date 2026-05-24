package example

import "fmt"

func Execute(a, b int) int {
	return addInternal(a, b)
}

func Version() string {
	return "1.0.0"
}

func Logger(message string) {
	fmt.Println(message)
}

func addInternal(a, b int) int {
	return a + b
}

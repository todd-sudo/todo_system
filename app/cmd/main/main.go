package main

import "github.com/todd-sudo/todo_system/internal/app"

func main() {
	saveToFile := false
	app.RunApplication(saveToFile)
}

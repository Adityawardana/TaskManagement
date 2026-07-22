package main

import (
	"fmt"
	"runtime/debug"

	service "taskmanagement/internal"
)

const serviceName = "taskmanagement"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\x1b[31;1mFailed to start %s service: %v\x1b[0m\n", serviceName, r)
			fmt.Printf("Stack trace: \n%s\n", debug.Stack())
		}
	}()

	srv := service.NewService(serviceName)
	srv.Run()
}

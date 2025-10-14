package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	// "flag"
)

// Objective: Create a CLI that allows users to add, view, and complete tasks.

type Task struct {
		id int
		title string
		completed bool
		createdAt string
	}

func main() {
	var tasks []Task

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("To edit: /e [id] [text], To quit: /q")

	for {
		fmt.Print("Enter task: ")
		scanner.Scan()
		command := strings.TrimSpace(scanner.Text())
		if command == "" {
			continue
		}
		if command == "/q" {
			break
		}
		tasks = append(tasks, Task{
			id: findId(tasks),
			title: command,
			completed: false,
			createdAt: time.Now().Format(time.RFC3339),
		})

		fmt.Printf("Your tasks: %+v \n", tasks)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func findId(tasks []Task) int {
	if len(tasks) == 0 {
		return 0
	}
	id := tasks[len(tasks)-1].id + 1
	return id
}

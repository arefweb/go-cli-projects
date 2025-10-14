package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/shlex"
)

// Objective: Create a CLI that allows users to add, view, and complete tasks.

type Task struct {
	id        int
	title     string
	completed bool
	createdAt string
}

func main() {
	var tasks []Task

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("To edit: /e [id] [title], To quit: /q")

	for {
		fmt.Print("Task: ")
		scanner.Scan()
		command := strings.TrimSpace(scanner.Text())
		if command == "" {
			continue
		}
		if command == "/q" {
			break
		}

		if strings.HasPrefix(command, "/e") {
			// Use shlex to split input, preserving quoted strings
			args, err := shlex.Split(command)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error splitting input:", err)
				continue
			}

			// Skip the "/e" command itself
			if len(args) < 1 {
				fmt.Fprintln(os.Stderr, "Error: /e command requires arguments")
				continue
			}
			args = args[1:] // Remove "/e" from args

			// Parse flags for edit command
			fs := flag.NewFlagSet("edit", flag.ContinueOnError)
			id := fs.Int("id", -1, "ID of the task to edit")
			text := fs.String("text", "", "New text for the task")

			if err := fs.Parse(args); err != nil {
				fmt.Fprintln(os.Stderr, "Error parsing flags:", err)
				continue
			}

			// Validate flags
			if *id == -1 {
				fmt.Fprintln(os.Stderr, "Error: -id flag is required")
				continue
			}
			if *text == "" {
				fmt.Fprintln(os.Stderr, "Error: -text flag is required")
				continue
			}

			// Edit the task
			if err := editTask(&tasks, *id, *text); err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fmt.Printf("Task %d updated successfully\n", *id)
			printTasks(tasks)
			continue
		}

		tasks = append(tasks, Task{
			id:        findId(tasks),
			title:     command,
			completed: false,
			createdAt: time.Now().Format(time.RFC3339),
		})

		printTasks(tasks)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func findId(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	id := tasks[len(tasks)-1].id + 1
	return id
}

func editTask(tasks *[]Task, id int, newText string) error {
	for i, task := range *tasks {
		if task.id == id {
			(*tasks)[i].title = newText
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func printTasks(tasks []Task) {
    fmt.Println("Your tasks:")
    for _, t := range tasks {
        fmt.Printf("  [%d] %s (done=%v)\n", t.id, t.title, t.completed)
    }
}
// Objective: Create a CLI that allows users to add, view, and complete tasks.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/google/shlex"
)


type Task struct {
	id        int
	title     string
	completed bool
	createdAt string
}

func main() {
	var tasks []Task

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("*********************************")
	fmt.Println("* Welcome, these are the commands: ")
	fmt.Println("* To edit: /e -id=3 -title=\"New Title\" completed=true")
	fmt.Println("* To list tasks: /l")
	fmt.Println("* To Delete a task: /d -id=2")
	fmt.Println("* To quit: /q")
	fmt.Println("*********************************")

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
			title := fs.String("title", "", "New title for the task")
			completed := fs.Bool("completed", false, "Whether task completed")

			if err := fs.Parse(args); err != nil {
				fmt.Fprintln(os.Stderr, "Error parsing flags:", err)
				continue
			}

			// Validate flags
			if *id == -1 {
				fmt.Fprintln(os.Stderr, "Error: -id flag is required")
				continue
			}

			// Edit the task
			if err := editTask(&tasks, *id, title, completed); err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fmt.Printf("Task %d updated successfully\n", *id)
			continue
		}

		if command == "/l" {
			printTasks(tasks)
			continue
		}

		// Delete a task
		if strings.HasPrefix(command, "/d") {
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
			args = args[1:]

			var idToRemove int
			var indexToRemove int
			fs := flag.NewFlagSet("delete", flag.ContinueOnError)
			id := fs.Int("id", -1, "ID of the task to delete")

			if err := fs.Parse(args); err != nil {
				fmt.Fprintln(os.Stderr, "Error parsing flags:", err)
				continue
			}
			
			if *id == -1 {
				fmt.Fprintln(os.Stderr, "Error: -id flag is required")
				continue
			}
			for i, item := range tasks {
				if item.id == *id {
					idToRemove = item.id
					indexToRemove = i
					break
				}
			}
			if idToRemove == 0 {
				fmt.Println("No valid id provided")
				continue
			}
			tasks = slices.Delete(tasks, indexToRemove, indexToRemove+1)
			fmt.Printf("Deleted id %d successfully \n", idToRemove)
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

func editTask(tasks *[]Task, id int, newTitle *string, completed *bool) error {
	for i, task := range *tasks {
		if task.id == id {
			if newTitle != nil && *newTitle != "" {
				(*tasks)[i].title = *newTitle
			}
			if completed != nil {
				(*tasks)[i].completed = *completed
			}
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func printTasks(tasks []Task) {
	fmt.Println("Your tasks:")
	for _, t := range tasks {
		fmt.Printf("  [%d] %s (completed=%v)\n", t.id, t.title, t.completed)
	}
}

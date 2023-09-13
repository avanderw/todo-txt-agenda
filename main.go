package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath" // Import the filepath package
	"sort"
	"strings"
	"time"
)

// Task represents a todo.txt task.
type Task struct {
	Text      string
	DueDate   time.Time
	FilePath  string // Store the base file name
	Completed bool
}

func main() {
	// Read the list of todo.txt files from the text file.
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <files.txt>\n", os.Args[0])
		os.Exit(1)
	}
	filesFilePath := os.Args[1]
	filesFile, err := os.Open(filesFilePath)
	if err != nil {
		fmt.Printf("Failed to open file %s: %v\n", filesFilePath, err)
		os.Exit(1)
	}
	defer filesFile.Close()

	var todoFiles []string
	scanner := bufio.NewScanner(filesFile)
	for scanner.Scan() {
		filePath := expandTilde(scanner.Text()) // Expand ~ to home directory
		todoFiles = append(todoFiles, filePath)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filesFilePath, err)
		os.Exit(1)
	}

	// Parse and extract tasks from todo.txt files, including the base file name.
	var tasks []Task
	for _, filePath := range todoFiles {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Failed to open file %s: %v\n", filePath, err)
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			taskLine := scanner.Text()
			task, err := parseTask(taskLine, filePath) // Pass file path to parseTask
			if err != nil {
				fmt.Printf("Failed to parse task: %v\n", err)
				continue
			}

			// Exclude completed tasks and tasks without due dates.
			if !task.Completed && !task.DueDate.IsZero() {
				tasks = append(tasks, task)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
		}
	}

	// Sort tasks by due date.
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].DueDate.Before(tasks[j].DueDate)
	})

	// Create a weekly agenda view with past due and upcoming tasks.
	currentDate := time.Now()
	endDate := currentDate.AddDate(0, 0, 7)

	fmt.Printf("Weekly Agenda\n")

	// Display past due tasks as a separate heading.
	fmt.Printf("\n[PAST DUE]\n")
	pastDueDisplayed := false
	for _, task := range tasks {
		if task.DueDate.Before(currentDate) {
			fmt.Printf("  - %s (%s)\n", task.Text, task.FilePath)
			pastDueDisplayed = true
		}
	}

	// Continue with the weekly agenda.
	for currentDate.Before(endDate) {
		fmt.Printf("\n%s\n", currentDate.Format("Monday, January 2, 2006"))

		// Display upcoming tasks.
		for _, task := range tasks {
			if task.DueDate.After(currentDate) && task.DueDate.Before(currentDate.AddDate(0, 0, 1)) {
				fmt.Printf("  - %s (%s)\n", task.Text, task.FilePath)
			}
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// If no past due tasks were displayed, remove the [PAST DUE] heading.
	if !pastDueDisplayed {
		fmt.Printf("\n")
	}
}

func parseTask(taskLine string, filePath string) (Task, error) {
	// Parse the task line in todo.txt format.
	parts := strings.Fields(taskLine)

	var task Task
	task.Text = strings.Join(parts[1:], " ")
	task.FilePath = filepath.Base(filePath) // Store only the base file name
	task.FilePath = strings.TrimSuffix(task.FilePath, ".todo.txt")

	for _, part := range parts {
		if strings.HasPrefix(part, "due:") {
			dueDateStr := strings.TrimPrefix(part, "due:")
			dueDate, err := parseDueDate(dueDateStr)
			if err != nil {
				return Task{}, err
			}
			task.DueDate = dueDate
			break
		}
		if part == "x" {
			task.Completed = true
			break
		}
	}

	return task, nil
}

func parseDueDate(dueDateStr string) (time.Time, error) {
	// Parse the due date manually (format: yyyy-mm-dd) and set the time to the end of the day (23:59:59).
	parsedDate, err := time.Parse("2006-01-02", dueDateStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 0, parsedDate.Location()), nil
}

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		homeDir, _ := os.UserHomeDir()
		return strings.Replace(path, "~", homeDir, 1)
	}
	return path
}

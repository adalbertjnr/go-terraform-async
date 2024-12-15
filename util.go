package main

import (
	"log/slog"
	"strings"
)

func logErrors(errors map[string]error) {
	for err, msg := range errors {
		slog.Error("validate errors", err, msg)
	}
}

func validateVerb(inputVerb string) bool {
	verbs := map[string]bool{
		Plan:    true,
		Apply:   true,
		Destroy: true,
	}
	return verbs[strings.ToLower(inputVerb)]
}

func fetchOptions(data string) []string {
	var tasks []string
	var sanitizedTasks []string

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		tasks = append(tasks, strings.Split(line, ",")...)
	}

	for _, task := range tasks {
		if task != "" {
			sanitizedTasks = append(sanitizedTasks, task)
		}
	}
	return sanitizedTasks
}

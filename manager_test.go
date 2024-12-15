package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTerraformManager(t *testing.T) {
	var inputUser string
	var inputVerb string
	var inputTasks string
	var inputVersion string
	var inputWorker int
	var enableDebug bool

	flag.StringVar(&inputUser, "user", "terraform", "the username who started the pipeline")
	flag.StringVar(&inputVerb, "verb", "plan", "set the verb - plan - apply - destroy")
	flag.StringVar(&inputTasks, "tasks", "", "tasks that will be executed by terraform")
	flag.StringVar(&inputVersion, "version", "1.9.5", "terraform version")
	flag.BoolVar(&enableDebug, "DEBUG", false, "enable debug option")
	flag.IntVar(&inputWorker, "workers", 3, "set the number of workers to run concurrently")
	flag.Parse()

	tasks := []string{"testCase1", "testCase2", "testCase3", "testCase4", "testCase5", "testCase6", "testCase7", "testCase8", "testCase9", "testCase10"}

	svc, err := NewTerraformService(terraformInstaller, inputVersion)
	if err != nil {
		t.Fatalf("terraform service failed: err %v", err)
	}

	enableDebug = true

	var logLevel slog.Level
	logLevel = slog.LevelInfo

	if enableDebug {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	var locations []string
	for _, taskLocation := range tasks {
		temp, err := createTemp(taskLocation)
		if err != nil {
			t.Fatalf("create temp failed: %v", err)
		}
		locations = append(locations, temp)
	}

	inputTasks = strings.Join(locations, ",")

	task := NewTaskManager(
		context.Background(),
		inputUser,
		inputVerb,
		inputTasks,
		inputVersion,
		inputWorker,
		svc,
	)

	task.start()

	cleanupTasks(t, locations)
}

func cleanupTasks(t *testing.T, locations []string) {
	matches, err := filepath.Glob("/tmp/terraform*")
	if err != nil {
		t.Fatal("matches error", err)
	}

	for _, match := range matches {
		if err := os.RemoveAll(match); err != nil {
			t.Fatal(err)
		}
	}

	for _, location := range locations {
		if err := os.RemoveAll(location); err != nil {
			t.Fatal(err)
		}
	}
}

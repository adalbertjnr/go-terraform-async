package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
)

func main() {
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
	flag.BoolVar(&enableDebug, "debug", false, "enable debug option")
	flag.IntVar(&inputWorker, "workers", 2, "set the number of workers to run concurrently")
	flag.Parse()

	var logLevel slog.Level
	logLevel = slog.LevelInfo

	if enableDebug {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	slog.Debug("init input variables", "user", inputUser)
	slog.Debug("init input variables", "verb", inputVerb)
	slog.Debug("init input variables", "tasks", inputTasks)
	slog.Debug("init input variables", "version", inputVersion)
	slog.Debug("init input variables", "workers", inputWorker)
	slog.Debug("init input variables", "debug", enableDebug)

	terraform, err := NewTerraformService(terraformInstaller, inputVersion)
	if err != nil {
		log.Fatal(err)
	}

	task := NewTaskManager(
		context.Background(),
		inputUser,
		inputVerb,
		inputTasks,
		inputVersion,
		inputWorker,
		terraform,
	)

	task.start()
}

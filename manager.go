package main

import (
	"context"
	"log/slog"
	"strings"
)

const (
	Apply   = "apply"
	Destroy = "destroy"
	Plan    = "plan"
)

type TaskManager struct {
	ctx   context.Context
	tfsvc *TFService
	args  input
	tasks []string
}

type input struct {
	inputUser    string
	inputVerb    string
	inputTasks   string
	inputVersion string
	inputWorker  int
}

func NewTaskManager(ctx context.Context, inputUser, inputVerb, inputTasks, inputVersion string, inputWorker int, svc *TFService) *TaskManager {
	args := input{
		inputUser:    inputUser,
		inputVerb:    inputVerb,
		inputTasks:   inputTasks,
		inputVersion: inputVersion,
		inputWorker:  inputWorker,
	}

	return &TaskManager{
		ctx:   ctx,
		args:  args,
		tfsvc: svc,
	}
}

func (w *TaskManager) retrieveTasks() error {
	data := w.args.inputTasks
	slog.Debug("task manager", "fetching options from", data)

	if len(fetchOptions(data)) > 0 {
		w.tasks = fetchOptions(data)
		slog.Debug("task manager", "fetching options from", data, "status", "ok")
		return nil
	}

	return ErrEmptyTask
}

func (w *TaskManager) validateArgs() map[string]error {
	slog.Debug("task manager", "validating arguments", w.args.inputTasks)

	errors := make(map[string]error)
	if w.args.inputTasks == "" {
		errors["TASK"] = ErrEmptyTask
	}
	if !validateVerb(w.args.inputVerb) {
		errors["VERB"] = ErrVerbNotFound
	}

	return errors
}

func (w *TaskManager) start() {
	if errors := w.validateArgs(); len(errors) > 0 {
		logErrors(errors)
		return
	}

	if err := w.retrieveTasks(); err != nil {
		slog.Error("retrieve tasks", "error", err)
		return
	}

	slog.Info("task manager", "started by user", w.args.inputUser, "verb", strings.ToUpper(w.args.inputVerb))

	doneChannel := make(chan struct{})
	taskChannel := make(chan string, len(w.tasks))

	fillTaskChannel(taskChannel, w.tasks)
	close(taskChannel)

	for i := 0; i < w.args.inputWorker; i++ {
		go w.terraformWorker(i, w.args.inputVerb, taskChannel, doneChannel)
	}

	w.wait(w.args.inputWorker, doneChannel)
}

func (w *TaskManager) wait(workers int, doneChannel chan struct{}) {
	for i := 0; i < workers; i++ {
		slog.Debug("task manager", "waiting", i)
		<-doneChannel
	}
}
func fillTaskChannel(taskChannel chan<- string, terraformTasks []string) {
	for _, terraformTask := range terraformTasks {
		slog.Debug("task manager", "filling channel with", terraformTask)
		taskChannel <- terraformTask
	}
}

func (t *TaskManager) terraformWorker(workerId int, verb string, taskChannel chan string, doneChannel chan struct{}) {
	defer func() {
		doneChannel <- struct{}{}
		slog.Debug("task manager", "id", workerId, "status", "exiting")
	}()

	for task := range taskChannel {
		slog.Info("task manager", "aws account", task, "worker", workerId, "verb", strings.ToUpper(verb))

		switch verb {
		case Apply:
			slog.Debug("task manager", "stage", "terraform task apply")

			err := t.tfsvc.terraformTaskCreate(t.ctx, task, t.tfsvc.execPath)
			if err != nil {
				slog.Error("task manager", "apply stage error", err)
			}
		case Destroy:
			slog.Debug("task manager", "stage", "terraform task destroy")

			err := t.tfsvc.terraformTaskDestroy(t.ctx, task, t.tfsvc.execPath)
			if err != nil {
				slog.Error("task manager", "destroy stage error", err)
			}
		case Plan:
			slog.Debug("task manager", "stage", "terraform task plan")

			err := t.tfsvc.terraformTaskPlan(t.ctx, task, t.tfsvc.execPath)
			if err != nil {
				slog.Error("task manager", "plan stage error", err)
			}
		default:
			slog.Warn("task manager", "verb", strings.ToUpper(verb), "status", "not found")
		}
	}
}

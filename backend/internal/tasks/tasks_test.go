package tasks

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func Test(t *testing.T) {
	manager := NewManager()
	defer manager.ShutdownWatcher()

	time.Sleep(200 * time.Millisecond) // нужно чтобы watchTasks успел запустится

	task := TaskData{
		Id:            uuid.Must(uuid.NewRandom()),
		Arg1:          10,
		Arg2:          5,
		Operation:     '+',
		OperationTime: 5000,
	}
	const except = 15
	const exceptSuccess = true

	t.Run("Add task", func(t *testing.T) {
		manager.AddTask(task, 15*time.Second)
		time.Sleep(400 * time.Millisecond) // Чтобы watchTasks успел это обработать
	})

	t.Run("Get task", func(t *testing.T) {
		taskResponse, ok := manager.GetTask()

		if !ok {
			t.Fatal(`manager.GetTask !ok`)
		}

		if taskResponse != task {
			t.Fatalf(`manager.GetTask = %v, want %v`, task, taskResponse)
		}
	})

	t.Run("Get empty task", func(t *testing.T) {
		_, ok := manager.GetTask()

		if ok {
			t.Fatal(`manager.GetTask is ok, want !ok`)
		}
	})

	t.Run("Add result", func(t *testing.T) {
		err := manager.AddResultToTask(task.Id, except, exceptSuccess)

		if err != nil {
			t.Fatalf(`manager.AddResultToTask has error: %s`, err.Error())
		}
	})

	t.Run("Add result to unknown task", func(t *testing.T) {
		err := manager.AddResultToTask(uuid.Must(uuid.NewRandom()), 10, true)

		if err == nil {
			t.Fatal(`manager.GetTaskResult does not have an error but it is necessary`)
		}
	})

	t.Run("Get result", func(t *testing.T) {
		result, err := manager.GetTaskResult(task.Id)

		if err != nil {
			t.Fatalf(`manager.GetTaskResult has error: %s`, err.Error())
		}

		if result.Successful != exceptSuccess {
			t.Fatalf(`manager.GetTaskResult.Successful = %t, want %t`, result.Successful, exceptSuccess)
		}

		if result.Result != except {
			t.Fatalf(`manager.GetTaskResult.Result = %t, want %t`, result.Successful, exceptSuccess)
		}
	})

	t.Run("Get unknown result", func(t *testing.T) {
		_, err := manager.GetTaskResult(task.Id)

		if err == nil {
			t.Fatal(`manager.GetTaskResult does not have an error but it is necessary`)
		}
	})

	t.Run("Failed run watcher", func(t *testing.T) {
		err := manager.RunWatcher()

		if err == nil {
			t.Fatal(`manager.RunWatcher does not have an error but it is necessary`)
		}
	})

	t.Run("Shutdown watcher", func(t *testing.T) {
		err := manager.ShutdownWatcher()

		if err != nil {
			t.Fatalf(`manager.ShutdownWatcher has error: %s`, err.Error())
		}

		time.Sleep(400 * time.Millisecond) // Чтобы watchTasks успел это обработать
	})

	t.Run("Failed shutdown watcher", func(t *testing.T) {
		err := manager.ShutdownWatcher()

		if err == nil {
			t.Fatalf(`manager.ShutdownWatcher does not have an error but it is necessary`)
		}
	})

	t.Run("Run watcher", func(t *testing.T) {
		err := manager.RunWatcher()

		if err != nil {
			t.Fatalf(`manager.RunWatcher has error: %s`, err.Error())
		}
	})

	t.Run("Add a task with a small timeout", func(t *testing.T) {
		task := TaskData{
			Id:        uuid.Must(uuid.NewRandom()),
			Operation: '*',
		}
		manager.AddTask(task, 100*time.Millisecond)
	})

	t.Run("Get task with a small timeout", func(t *testing.T) {
		_, ok := manager.GetTask()

		if !ok {
			t.Fatal(`manager.GetTask !ok`)
		}

		time.Sleep(400 * time.Millisecond) // Чтобы watchTasks успел это обработать)
	})
}

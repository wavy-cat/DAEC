package tasks

import (
	"backend/internal/storage"
	"backend/pkg/queue"
	"errors"
	"github.com/google/uuid"
	"time"
)

// Manager Это верхнеуровневый менеджер тасок между оркестратором и агентом
type Manager struct {
	queue       *queue.Queue[TaskWrapper] // Очередь нерешённых задач для агента
	database    *storage.Storage[Task]    // База данных, где хранятся задачи
	watchStatus chan interface{}          // Канал, который регулирует работу watchTasks
}

func NewManager() *Manager {
	manager := Manager{
		queue:       queue.NewQueue[TaskWrapper](),
		database:    storage.NewStorage[Task](),
		watchStatus: make(chan interface{}),
	}

	go manager.watchTasks()

	return &manager
}

// Эта функция будет смотреть за тасками.
// Если какая-то из них долго выполняется (выходит за таймаут),
// то она будет заново добавлена в очередь.
func (m *Manager) watchTasks() {
	working := true
	for {
		select {
		case <-m.watchStatus:
			// Обработка завершения работы
			working = false
			break
		default:
			// Работаем
			tasks := m.database.GetAll()
			for _, task := range tasks {
				if task.Status != "processing" {
					continue
				}
				if time.Now().After(task.CompleteBefore) {
					task.Status = "queue"
					m.database.Set(task.Data.Id, task)
					m.queue.Enqueue(TaskWrapper{Task: task.Data})
				}
			}
		}

		if !working {
			break
		}

		time.Sleep(200 * time.Millisecond)
	}
}

// RunWatcher Запускает watchTasks
func (m *Manager) RunWatcher() {
	m.watchStatus = make(chan interface{})
	go m.watchTasks()
}

// ShutdownWatcher Останавливает работу watchTasks
func (m *Manager) ShutdownWatcher() {
	close(m.watchStatus)
}

// AddTask добавляет новую задачу
func (m *Manager) AddTask(task TaskData, timeout time.Duration) {
	m.database.Set(task.Id, Task{
		Data:    task,
		Timeout: timeout,
		Status:  "queue",
	}) // Добавляем задачу в БД
	m.queue.Enqueue(TaskWrapper{Task: task}) // Добавляем задачу в очередь обработки
}

// GetTask возвращает задачу для обработки
func (m *Manager) GetTask() (TaskWrapper, bool) {
	task, ok := m.queue.Dequeue()
	if !ok {
		return task, false
	}

	t, ok := m.database.Get(task.Task.Id)
	if !ok {
		return task, false
	}

	t.Status = "processing"
	t.CompleteBefore = time.Now().Add(t.Timeout)
	m.database.Set(task.Task.Id, t)
	return task, true
}

// AddResultToTask добавляет решение к задаче
func (m *Manager) AddResultToTask(id uuid.UUID, result float64, successful bool) error {
	task, ok := m.database.Get(id)
	if !ok {
		return errors.New("task not found")
	}
	task.Result = result
	task.Successful = successful
	task.Status = "done"
	m.database.Set(id, task)
	return nil
}

// GetTaskResult Отдаёт результат задачи
func (m *Manager) GetTaskResult(id uuid.UUID) (TaskResult, error) {
	task, ok := m.database.Get(id)
	if !ok {
		return TaskResult{}, errors.New("task not found")
	}
	var done bool

	if task.Status == "done" {
		done = true
		m.database.Del(id)
	}

	return TaskResult{
		IsDone:     done,
		Result:     task.Result,
		Successful: task.Successful,
	}, nil
}

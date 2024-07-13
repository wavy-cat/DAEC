package tasks

import (
	"errors"
	"github.com/google/uuid"
	"github.com/wavy-cat/DAEC/backend/internal/storage"
	"github.com/wavy-cat/DAEC/backend/pkg/queue"
	"time"
)

// Manager Это верхнеуровневый менеджер тасок между оркестратором и агентом
type Manager struct {
	queue       *queue.Queue[TaskData] // Очередь нерешённых задач для агента
	database    *storage.Storage[Task] // База данных, где хранятся задачи
	watchStatus chan interface{}       // Канал, который регулирует работу watchTasks
}

func NewManager() *Manager {
	manager := Manager{
		queue:       queue.NewQueue[TaskData](),
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
	for {
		select {
		case <-m.watchStatus:
			// Обработка завершения работы
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
					m.queue.Enqueue(task.Data)
				}
			}
		}

		time.Sleep(200 * time.Millisecond)
	}
}

// RunWatcher Запускает watchTasks
func (m *Manager) RunWatcher() error {
	select {
	case <-m.watchStatus:
		m.watchStatus = make(chan interface{})
		go m.watchTasks()
	default:
		return errors.New("watcher is already running")
	}

	return nil
}

// ShutdownWatcher Останавливает работу watchTasks
func (m *Manager) ShutdownWatcher() error {
	select {
	case <-m.watchStatus:
		return errors.New("watcher no longer works")
	default:
		close(m.watchStatus)
	}

	return nil
}

// AddTask добавляет новую задачу
func (m *Manager) AddTask(task TaskData, timeout time.Duration) {
	m.database.Set(task.Id, Task{
		Data:    task,
		Timeout: timeout,
		Status:  "queue",
	}) // Добавляем задачу в БД
	m.queue.Enqueue(task) // Добавляем задачу в очередь обработки
}

// GetTask возвращает задачу для обработки
func (m *Manager) GetTask() (TaskData, bool) {
	task, ok := m.queue.Dequeue()
	if !ok {
		return task, false
	}

	t, _ := m.database.Get(task.Id) // По идее нет ситуации когда задачи в базе не окажется

	t.Status = "processing"
	t.CompleteBefore = time.Now().Add(t.Timeout)
	m.database.Set(task.Id, t)
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

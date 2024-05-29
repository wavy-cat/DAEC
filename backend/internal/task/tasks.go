package task

import (
	"errors"
	"sync"
)

// Task структура задачи (вычисления арифметического выражения)
type Task struct {
	Id         int    // Идентификатор выражения
	Expression string // Строка с выражением
	Status     string // Статус задачи
	Result     int    // Результат
}

// Tasks структура, содержащая список задач в работе и готовых.
type Tasks struct {
	pendingTasks []Task       // Задачи в работе
	doneTasks    []Task       // Готовые задачи
	mutex        sync.RWMutex // Mutex для синхронизации
}

// Функция для проверки есть ли задача с переданным Id в pendingTasks.
// Если задача найдена возвращается ещё её и индекс.
func (t *Tasks) containsPending(taskId int) (int, bool) {
	for i, tk := range t.pendingTasks {
		if tk.Id == taskId {
			return i, true
		}
	}
	return 0, false
}

// Функция для проверки есть ли задача с переданным Id в doneTasks.
// Если задача найдена возвращается ещё её и индекс.
func (t *Tasks) containsDone(taskId int) (int, bool) {
	for i, tk := range t.doneTasks {
		if tk.Id == taskId {
			return i, true
		}
	}
	return 0, false
}

// AddTask добавляет задачу в работу
func (t *Tasks) AddTask(task Task) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Проверяем выполняется ли уже задача
	if _, ok := t.containsPending(task.Id); ok {
		return errors.New("the task is already being processed")
	}
	if _, ok := t.containsDone(task.Id); ok {
		return errors.New("the task has already been completed")
	}

	t.pendingTasks = append(t.pendingTasks, task)

	// TODO: Логика обратной польской записи
	// TODO: Отправка агенту (в другой горутине)

	return nil
}

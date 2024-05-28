package work

import (
	"go.uber.org/zap"
	"sync"
)

// Pool Пул для выполнения
type Pool struct {
	tasks  chan Expression // Канал, из которого будут браться задачи для обработки
	wg     sync.WaitGroup
	logger *zap.Logger
}

// NewPool Конструктор структуры Pool
func NewPool(maxGoroutines int, logger *zap.Logger) *Pool {
	p := Pool{tasks: make(chan Expression), logger: logger}

	p.wg.Add(maxGoroutines)

	for i := 0; i < maxGoroutines; i++ {
		// создадим горутины по указанному количеству maxGoroutines
		go func() {
			// забираем задачи из канала
			for w := range p.tasks {
				// выполняем
				w.Execute()

				// отправляем результат
				result := ResultData{Id: w.Id, Result: w.Result}
				err := SendResult(result, 3)
				if err != nil {
					logger.Error("Error sending result: " + err.Error())
				}
			}
			// после закрытия канала нужно оповестить наш пул
			p.wg.Done()
		}()
	}

	return &p
}

// Run Добавляет новую задачу в пул
func (p *Pool) Run(w Expression) {
	p.tasks <- w
}

// Shutdown Остановка работы пула
func (p *Pool) Shutdown() {
	close(p.tasks)
	p.wg.Wait()
}

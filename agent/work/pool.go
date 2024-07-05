package work

import (
	"context"
	pb "github.com/wavy-cat/DAEC/agent/proto"
	"go.uber.org/zap"
	"sync"
)

// Pool Пул для выполнения
type Pool struct {
	tasks  chan Expression // Канал, из которого будут браться задачи для обработки
	wg     sync.WaitGroup
	logger *zap.Logger
	client pb.TasksServiceClient
}

// NewPool Конструктор структуры Pool
func NewPool(maxGoroutines int, logger *zap.Logger, client pb.TasksServiceClient) *Pool {
	p := Pool{tasks: make(chan Expression), logger: logger, client: client}

	p.wg.Add(maxGoroutines)

	for i := 0; i < maxGoroutines; i++ {
		// создадим горутины по указанному количеству maxGoroutines
		go func() {
			// забираем задачи из канала
			for w := range p.tasks {
				// выполняем
				w.Execute()

				// отправляем результат
				var request pb.PushTaskRequest

				switch w.Successful {
				case true:
					request = pb.PushTaskRequest{
						Id:         w.Id.String(),
						Result:     w.Result,
						Successful: true,
					}
				default:
					request = pb.PushTaskRequest{
						Id:         w.Id.String(),
						Successful: false,
					}
				}

				_, err := client.Push(context.TODO(), &request)
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

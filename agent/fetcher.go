package main

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/wavy-cat/DAEC/agent/proto"
	"github.com/wavy-cat/DAEC/agent/work"
	"go.uber.org/zap"
	"strings"
	"time"
)

// Функция для получения новых задач от оркестратора
func fetcher(pool *work.Pool, logger *zap.Logger, client pb.TasksServiceClient) {
	for {
		// Отправляем запрос оркестратору
		task, err := client.Pull(context.TODO(), &pb.Empty{})
		if err != nil {
			if strings.Contains(err.Error(), "no task yet") {
				// Задачи пока нет, ждём-с
				time.Sleep(200 * time.Millisecond)
			} else {
				// Какая-то ошибка, логируем
				logger.Error(err.Error())
			}
			continue
		}

		// Создаём новую задачу и отправляем её в пул
		id, err := uuid.Parse(task.Id)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		expression := work.Expression{
			Id:            id,
			Num1:          task.Arg1,
			Num2:          task.Arg2,
			Operator:      ConvertOperationToRune(task.Operation),
			OperationTime: time.Duration(task.OperationTime) * time.Millisecond,
		}
		pool.Run(expression)
	}
}

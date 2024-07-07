package evaluate

import (
	"github.com/google/uuid"
	"github.com/wavy-cat/DAEC/backend/internal/config"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"github.com/wavy-cat/DAEC/backend/pkg/postfix"
	"time"
)

type solver struct {
	manager *tasks.Manager
}

func (s *solver) solve(num1, num2 float64, operator byte, sleepTime int) (tasks.TaskResult, bool) {
	id := uuid.Must(uuid.NewRandom())
	sleepTimeDuration := time.Duration(sleepTime) * time.Millisecond
	s.manager.AddTask(tasks.TaskData{
		Id:            id,
		Arg1:          num1,
		Arg2:          num2,
		Operation:     operator,
		OperationTime: uint32(sleepTime),
	}, sleepTimeDuration*3+sleepTimeDuration/2)

	var r tasks.TaskResult
	var err error

	for r, err = s.manager.GetTaskResult(id); err != nil || r.IsDone == false; r, err = s.manager.GetTaskResult(id) {
		if err != nil {
			return r, false
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !r.Successful {
		return r, false
	}

	return r, true
}

func (s *solver) generalSolve(num1, num2 *postfix.LazyFloat, operator byte, sleepTime int) *postfix.LazyFloat {
	var float postfix.LazyFloat
	go func() {
		// Получаем значения двух операндов
		n1, n2 := num1.GetValue(), num2.GetValue()
		if num1.IsFail || num2.IsFail {
			float.SetFail()
			return
		}

		// Высчитываем результат
		r, ok := s.solve(n1, n2, operator, sleepTime)
		if !ok {
			float.SetFail()
			return
		}

		// Устанавливаем значение
		float.SetValue(r.Result)
	}()
	return &float
}

func (s *solver) Addition(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	return s.generalSolve(num1, num2, '+', config.TimeAdditionMs)
}

func (s *solver) Subtraction(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	return s.generalSolve(num1, num2, '-', config.TimeSubtractionMs)
}

func (s *solver) Division(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	return s.generalSolve(num1, num2, '/', config.TimeDivisionsMs)
}

func (s *solver) Multiplication(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	return s.generalSolve(num1, num2, '*', config.TimeMultiplicationsMs)
}

func (s *solver) Exponentiation(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	return s.generalSolve(num1, num2, '^', config.TimeExponentiationMs)
}

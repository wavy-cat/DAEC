package evaluate

import (
	"backend/internal/config"
	"backend/internal/tasks"
	"backend/pkg/postfix"
	"github.com/google/uuid"
	"time"
)

type solver struct {
	manager *tasks.Manager
}

func (s *solver) solve(num1, num2 float64, operator string, sleepTime int) (tasks.TaskResult, bool) {
	id := uuid.Must(uuid.NewRandom())
	sleepTimeDuration := time.Duration(sleepTime) * time.Millisecond
	s.manager.AddTask(tasks.TaskData{
		Id:            id,
		Arg1:          num1,
		Arg2:          num2,
		Operation:     operator,
		OperationTime: sleepTime,
	}, sleepTimeDuration*3+sleepTimeDuration/2)

	var r tasks.TaskResult
	var err error

	for r, err = s.manager.GetTaskResult(id); r.IsDone == false; r, err = s.manager.GetTaskResult(id) {
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

func (s *solver) Addition(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		// Получаем значения двух операндов
		n1, n2 := num1.GetValue(), num2.GetValue()
		if num1.IsFail || num2.IsFail {
			float.SetFail()
			return
		}

		// Складываем их
		r, ok := s.solve(n1, n2, "+", config.TimeAdditionMs)
		if !ok {
			float.SetFail()
			return
		}

		// Устанавливаем значение
		float.SetValue(r.Result)
	}()

	return &float
}

func (s *solver) Subtraction(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		n1, n2 := num1.GetValue(), num2.GetValue()
		if num1.IsFail || num2.IsFail {
			float.SetFail()
			return
		}

		r, ok := s.solve(n1, n2, "-", config.TimeSubtractionMs)
		if !ok {
			float.SetFail()
			return
		}

		float.SetValue(r.Result)
	}()

	return &float
}

func (s *solver) Division(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		n1, n2 := num1.GetValue(), num2.GetValue()
		if num1.IsFail || num2.IsFail {
			float.SetFail()
			return
		}

		r, ok := s.solve(n1, n2, "/", config.TimeDivisionsMs)
		if !ok {
			float.SetFail()
			return
		}

		float.SetValue(r.Result)
	}()

	return &float
}

func (s *solver) Multiplication(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		n1, n2 := num1.GetValue(), num2.GetValue()
		if num1.IsFail || num2.IsFail {
			float.SetFail()
			return
		}

		r, ok := s.solve(n1, n2, "*", config.TimeMultiplicationsMs)
		if !ok {
			float.SetFail()
			return
		}

		float.SetValue(r.Result)
	}()

	return &float
}

func (s *solver) Exponentiation(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		n1, n2 := num1.GetValue(), num2.GetValue()
		if num1.IsFail || num2.IsFail {
			float.SetFail()
			return
		}

		r, ok := s.solve(n1, n2, "^", config.TimeExponentiationMs)
		if !ok {
			float.SetFail()
			return
		}

		float.SetValue(r.Result)
	}()

	return &float
}

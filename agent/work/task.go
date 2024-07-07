package work

import (
	"github.com/google/uuid"
	"math"
	"time"
)

// Expression структура выражения
type Expression struct {
	Id            uuid.UUID
	Num1          float64
	Num2          float64
	Operator      byte
	OperationTime time.Duration // Задержка в миллисекундах
	Result        float64
	Successful    bool
}

func (t *Expression) Execute() {
	time.Sleep(t.OperationTime)
	switch t.Operator {
	case '+':
		t.Result = t.Num1 + t.Num2
	case '-':
		t.Result = t.Num1 - t.Num2
	case '*':
		t.Result = t.Num1 * t.Num2
	case '/':
		t.Result = t.Num1 / t.Num2
	case '^':
		t.Result = math.Pow(t.Num1, t.Num2)
	}
	t.Successful = !(math.IsInf(t.Result, 1) || math.IsInf(t.Result, -1))
}

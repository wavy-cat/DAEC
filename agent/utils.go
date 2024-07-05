package main

import "github.com/wavy-cat/DAEC/agent/proto"

// ConvertOperationToRune converts the proto.Operation enum value to its corresponding rune symbol.
func ConvertOperationToRune(operation proto.Operation) rune {
	switch operation {
	case proto.Operation_ADDITION:
		return '+'
	case proto.Operation_SUBTRACTION:
		return '-'
	case proto.Operation_MULTIPLICATION:
		return '*'
	case proto.Operation_DIVISION:
		return '/'
	case proto.Operation_EXPONENTIATION:
		return '^'
	}
	panic("unknown operation")
}

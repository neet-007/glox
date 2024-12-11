package runtime

import (
	"fmt"

	"github.com/neet-007/glox/pkg/scanner"
)

type RuntimeError struct {
	Token   scanner.Token
	Message string
}

func NewRuntimeError(token scanner.Token, message string) *RuntimeError {
	return &RuntimeError{
		Token:   token,
		Message: message,
	}
}

func (r *RuntimeError) Error() string {
	return fmt.Sprintf("%v %v\n", r.Token, r.Message)
}

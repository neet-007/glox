package runtime

import "fmt"

type RuntimeError struct {
	Message string
}

func NewRuntimeError(message string) *RuntimeError {
	fmt.Printf("Creating RuntimeError with message: %s\n", message)
	return &RuntimeError{
		Message: message,
	}
}

func (r *RuntimeError) Error() string {
	return fmt.Sprintf("%v\n", r.Message)
}

package errors

import "fmt"

type DecodeError struct {
	Err error
}

func (err DecodeError) Error() string {
	wrap := fmt.Errorf("Failed to decode model\n%w", err.Err)
	return fmt.Sprintln(wrap.Error())
}

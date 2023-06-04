package errors

import "fmt"

type NotFoundError struct {
	Id int
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("User with id %d not found", err.Id)
}

package errors

import "fmt"

type NotFoundError struct {
	UserId int
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("User with id %d not found", err.UserId)
}

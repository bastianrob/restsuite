package exception

import "fmt"

//Exception interface
type Exception interface {
	Code() int
	String() string
}

//exception data model
type exception struct {
	Code    int
	Message string
}

func (e *exception) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

//New exception
func New(code int, message string, params ...interface{}) error {
	if params != nil && len(params) > 0 {
		message = fmt.Sprintf(message, params)
	}
	return &exception{code, message}
}

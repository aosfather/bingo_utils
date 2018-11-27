package bingo_utils

type MethodError struct {
	code int
	msg  string
}

func (this MethodError) Code() int {
	return this.code
}

func (this MethodError) Error() string {
	return this.msg
}

func CreateError(c int, text string) MethodError {
	var err MethodError
	err.code = c
	err.msg = text
	return err
}

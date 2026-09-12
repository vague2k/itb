package services

type UserError struct {
	msg string
}

func NewUserError(msg string) UserError {
	return UserError{msg: msg}
}

func (e UserError) Error() string { return e.msg }

type InternalError struct {
	msg string
	err error
}

func NewInternalError(err error, msg string) InternalError {
	return InternalError{msg: msg, err: err}
}

func (e InternalError) Error() string { return e.msg }

func (e InternalError) Unwrap() error { return e.err }

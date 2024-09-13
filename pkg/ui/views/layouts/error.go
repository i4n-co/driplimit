package layouts

type Error struct {
	Code    int
	Message string
	orig    error
}

func NewError(code int, message string, e ...error) *Error {
	if len(e) == 0 {
		return &Error{
			Code:    code,
			Message: message,
		}
	}
	return &Error{
		Code:    code,
		Message: message,
		orig:    e[0],
	}
}

func (e *Error) Error() string {
	if e.orig != nil {
		return e.orig.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.orig
}

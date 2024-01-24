package errors

type HTTPError struct {
	Message string
	Code    int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(message string, code int) *HTTPError {
	return &HTTPError{
		Message: message,
		Code:    code,
	}
}

func (e *HTTPError) StatusCode() int {
	return e.Code
}

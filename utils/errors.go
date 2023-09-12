package utils

//HandledError struct for generating errors can be identified as known
type HandledError struct {
	CapturedError error
	Code          string
	Message       string
}

func (e *HandledError) Error() string {
	return e.Message
}

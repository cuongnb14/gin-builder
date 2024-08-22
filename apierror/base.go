package apierror

type APIError struct {
	Status    int
	ErrorCode string
	Message   string
	Err       error
}

func (e APIError) Error() string {
	if e.Err != nil {
		return e.Message + ":\n	" + e.Err.Error()
	}
	return e.Message
}

func (e APIError) WithErr(err error) *APIError {
	e.Err = err
	return &e
}

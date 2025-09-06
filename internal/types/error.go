package types

import "net/http"

type HTTPErrorSwaggerWrapper HTTPError

type HTTPError struct {
	Code int   `json:"code"`
	Err  error `json:"error"`
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

var (
	ErrBadRequest          = func(err error) HTTPError { return HTTPError{Code: http.StatusBadRequest, Err: err} }
	ErrNotFound            = func(err error) HTTPError { return HTTPError{Code: http.StatusNotFound, Err: err} }
	ErrInternalServerError = func(err error) HTTPError { return HTTPError{Code: http.StatusInternalServerError, Err: err} }
)

package main

type SCAError struct {
	ErrorCode    uint8
	ErrorMessage string
}

const (
	// Server-side errors

	ErrServerFailure   = 0x01
	ErrUpstreamTimeout = 0x02

	// Client side timeouts

	ErrInvalidQuery      = 0x40
	ErrInvalidPagination = 0x41
	ErrInvalidComment    = 0x42
)

func (api simpleCommentsAPI) constructError(errorCode uint8, errorMessage string) SCAError {
	return SCAError{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
}

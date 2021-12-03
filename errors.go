package main

type SCAError struct {
	ErrorCode    uint8
	ErrorMessage string
}

const (
	// Server-side errors

	ERR_SERVER_FAILURE   = 0x01
	ERR_UPSTREAM_TIMEOUT = 0x02

	// Client side timeouts

	ERR_INVALID_QUERY      = 0x40
	ERR_INVALID_PAGINATION = 0x41
)

func (api simpleCommentsAPI) constructError(errorCode uint8, errorMessage string) SCAError {
	return SCAError{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
}

package consts

const (
	ErrnoSuccess      = 0
	ErrnoUnknownError = 1

	// param error 1XXX
	ErrnoBindRequestError     = 1000
	ErrnoRequestValidateError = 1001

	// mySQL error 2XXX

)

var ErrMsg = map[int]string{
	ErrnoSuccess:      "success",
	ErrnoUnknownError: "unknown error",

	ErrnoBindRequestError:     "bind request error",
	ErrnoRequestValidateError: "validate request error",
}

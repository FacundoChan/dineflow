package consts

const (
	ErrnoSuccess      = 0
	ErrnoUnknownError = 1

	// param error 1XXX
	ErrnoBindRequestError     = 1000
	ErrnoRequestValidateError = 1001
	ErrnoRequestNilItemsError = 1002

	// mySQL error 2XXX

	// Stripe error 3XXX
	ErrnoStripeUnknownError         = 3000
	ErrnoStripeResourceMissingError = 3001
	ErrnoStripeRateLimitError       = 3002 // Stripe rate limit exceeded (HTTP 429)
)

var ErrMsg = map[int]string{
	ErrnoSuccess:      "success",
	ErrnoUnknownError: "unknown error",

	ErrnoBindRequestError:     "bind request error",
	ErrnoRequestValidateError: "validate request error",

	// Stripe
	ErrnoStripeRateLimitError: "stripe rate limit exceeded",
}

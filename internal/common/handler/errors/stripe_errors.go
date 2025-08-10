package errors

import (
	"encoding/json"
	"errors"

	"github.com/FacundoChan/dineflow/common/consts"
	"google.golang.org/grpc/status"
)

type StripeError struct {
	Code      string `json:"code"`
	DocURL    string `json:"doc_url"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Param     string `json:"param"`
	RequestID string `json:"request_id"`
	LogURL    string `json:"request_log_url"`
	Type      string `json:"type"`
}

func ParseStripeError(err error) (int, error) {
	if err == nil {
		return consts.ErrnoSuccess, nil
	}

	statusErr, ok := status.FromError(err)
	if !ok {
		return consts.ErrnoUnknownError, err
	}

	desc := statusErr.Message()
	var stripeErr StripeError
	if jsonErr := json.Unmarshal([]byte(desc), &stripeErr); jsonErr != nil {
		return consts.ErrnoUnknownError, err
	}

	// Stripe Error Handlers
	switch stripeErr.Code {
	case "resource_missing":
		return consts.ErrnoStripeResourceMissingError, errors.New("stripe: resource missing")
	case "card_declined":
		return consts.ErrnoRequestValidateError, errors.New("stripe: card declined")
	case "rate_limit":
		// HTTP 429: Request rate limit exceeded
		return consts.ErrnoStripeRateLimitError, errors.New("stripe: rate limit exceeded")
	// more mapping...
	default:
		return consts.ErrnoUnknownError, err
	}
}

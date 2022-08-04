package valueObjects

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

const (
	HttpCodeSuccess       = 200
	HttpCodeBadRequest    = 400
	HttpCodeUnauthorized  = 401
	HttpCodeForbidden     = 403
	HttpCodeNotFound      = 404
	HttpCodeInternalError = 500
)

type HttpCodeReturn uint16

var clientErrorHashMap = map[message.ClientErrorType]HttpCodeReturn{
	message.ClientErrorBadRequest:   HttpCodeBadRequest,
	message.ClientErrorUnauthorized: HttpCodeUnauthorized,
	message.ClientErrorForbidden:    HttpCodeForbidden,
	message.ClientErrorNotFound:     HttpCodeNotFound,
}

func (h *HttpCodeReturn) MapClientErrorToHttpCode(errorType message.ClientErrorType) (HttpCodeReturn, error) {
	v, ok := clientErrorHashMap[errorType]
	if !ok {
		return 0, fmt.Errorf("invalid Error Type")
	}

	return v, nil
}

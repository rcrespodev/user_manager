package valueObjects

const (
	HttpCodeSuccess       = 200
	HttpCodeBadRequest    = 400
	HttpCodeUnauthorized  = 401
	HttpCodeFordibben     = 403
	HttpCodeNotFound      = 404
	HttpCodeInternalError = 500
)

type HttpCodeReturn uint16

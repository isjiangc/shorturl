package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	ErrEmailAlreadyUse     = newError(1001, "The email is already in use.")
	ErrDestUrlAlreadyExist = newError(1002, "The dest_url is already in exist.")
	ErrDestUrlIllegal      = newError(1003, "The dest_url is illegal")
	ErrOpenTypeIllegal     = newError(1004, "the open_type is illegal")
	ErrMemoIsEmpty         = newError(1005, "the memo is empty")
	ErrDestUrlNotExist     = newError(1006, "The dest_url is not exist")
	ErrShortUrlEmpty       = newError(1007, "The short_url is not empty")
)

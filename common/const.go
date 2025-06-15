package common

const (
	UserStatusUnknown = iota
	UserStatusDisabled
	UserStatusNormal
)

const (
	UserRoleNormal = "normal"
	UserRoleAdmin  = "admin"
)

const (
	RespCodeOK = iota
	RespCodeNotAuthed
	RespCodeInvalidRequest
	RespCodeUserNotAdmin
	RespCodeInternalError
	RespCodeUserAlreadyExists
	RespCodeAuthErr
	RespCodeDBErr
	RespCodeMethodNotAllowed
	RespCodeServiceNotFound
	RespCodeNotFound
)

const (
	RespMsgOK                = "OK"
	RespMsgNotAuthed         = "Not Authed"
	RespMsgInvalidRequest    = "Invalid Request"
	RespMsgUserNotAdmin      = "User Not Admin"
	RespMsgInternalError     = "Internal Error"
	RespMsgUserAlreadyExists = "User Already Exists"
	RespMsgAuthErr           = "Auth Err"
	RespMsgDBErr             = "DB Err"
	RespMsgMethodNotAllowed  = "Method Not Allowed"
	RespMsgServiceNotFound   = "Service Not Found"
)

const (
	UIDKey                 = "uid"
	AuthorizationKey       = "authorization"
	AuthorizationHeaderKey = "X-Authorization-Token"
)

const (
	ErrMsgNotAuthed         = "not authed"
	ErrMsgInvalidRequest    = "invalid request"
	ErrMsgUserNotAdmin      = "user not admin"
	ErrMsgInternalError     = "internal error"
	ErrMsgUserAlreadyExists = "user already exists"
	ErrMsgAuthErr           = "auth err"
	ErrMsgDBErr             = "db err"
	ErrMsgMethodNotAllowed  = "method not allowed"
)

const (
	ServiceLitefs = "litefs"
)

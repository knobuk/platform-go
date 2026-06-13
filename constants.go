package platform

const (
	CodeSuccess      = 2100
	CodeCreated      = 2101
	CodeBadRequest   = 4100
	CodeValidation   = 4101
	CodeUnauthorized = 4010
	CodeForbidden    = 4030
	CodeNotFound     = 4040
	CodeConflict     = 4090
	CodeServerError  = 5000
)

const (
	RoleReader    = "READER"
	RoleAuthor    = "AUTHOR"
	RoleAdmin     = "ADMIN"
	RoleBookstore = "BOOKSTORE"
)

const (
	MetadataUserID    = "x-user-id"
	MetadataUserRoles = "x-user-roles"
	MetadataRequestID = "x-request-id"
	MetadataTraceID   = "x-trace-id"
)

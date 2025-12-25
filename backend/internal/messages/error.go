package messages

const (
	NotFound                 = "API not found"
	FailedDataFromBody       = "Failed to get data from body"
	FailedInputFormat        = "Invalid input format"
	FailedValidationOccurred = "Validation error occurred"
	InvalidToken             = "Invalid token"
	InvalidRequestData       = "Invalid request data"
	InternalServerError      = "Internal server error"
	TooManyRequests          = "Too many requests, please try again later"
	Unauthorized             = "Unauthorized access"
	ErrInvalidIDFormat       = "Invalid ID format"
	ConflictError            = "Resource conflict"
	StatusNotFound           = "Resource not found"
)

const (
	ErrUserAlreadyExists      = "User with this email already exists"
	ErrCheckUser              = "Failed to check user"
	ErrInvalidEmailOrPassword = "Invalid email or password"
	ErrProcessPassword        = "Failed to process password"
	ErrGenerateAccessToken    = "Failed to generate access token"
	ErrCreateUser             = "Failed to create user"
	ErrCommitTransaction      = "Failed to commit transaction"
	ErrUserNotFound           = "User not found"
)

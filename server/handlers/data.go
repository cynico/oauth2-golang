package handlers

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	clientIDKey struct{}
)

const (
	maxInsertRetries = 10
	internalSEMsg = "internal_server_error"
	internalSEDesc = "Internal server error. Try again later."
)

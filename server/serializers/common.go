package serializers

// Error struct to store error codes and error message.
type Error struct {
	Code    int
	Message string
}

type SuccessResponse struct {
	Message string `json:"message"`
}

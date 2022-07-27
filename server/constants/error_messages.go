package constants

const (
	// Generic messages
	// TODO: all these messages are to be verified from Mike at the end
	GenericErrorMessage    = "something went wrong, please try again later"
	ConnectAccount         = "[Click here to link your Azure DevOps account](%s%s?channel_id=%s)"
	ConnectAccountFirst    = "You do not have any Azure Devops account connected, kindly link the account first"
	UserConnected          = "Your Azure Devops account is succesfully connected!"
	UserAlreadyConnected   = "Your Azure Devops account is already connected"
	UserDisconnected       = "Your Azure Devops account is now disconnected"
	UnableToDisconnectUser = "Unable to disconnect user"
	AuthAttemptExpired     = "authentication attempt expired, please try again"
	InvalidAuthState       = "invalid oauth state, please try again"
)

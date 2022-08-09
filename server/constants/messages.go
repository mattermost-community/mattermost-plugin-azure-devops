package constants

// TODO: all these messages are to be verified from Mike at the end
const (
	// Generic
	ConnectAccount       = "[Click here to link your Azure DevOps account](%s%s)"
	ConnectAccountFirst  = "You do not have any Azure Devops account connected. Kindly link the account first"
	UserConnected        = "Your Azure Devops account is succesfully connected!"
	UserAlreadyConnected = "Your Azure Devops account is already connected"
	UserDisconnected     = "Your Azure Devops account is now disconnected"
	CreatedTask          = "Link for newly created task: %s"
	TaskTitle            = "[%s #%d: %s](%s)"
	TaskPreviewMessage   = "State: %s\nAssigned To: %s\nDescription: %s"
	AlreadyLinkedProject = "This project is already linked."
	NoProjectLinked      = "No project is linked, please link a project."

	// Validations Errors
	OrganizationRequired   = "organization is required"
	ProjectRequired        = "project is required"
	TaskTypeRequired       = "task type is required"
	TaskTitleRequired      = "task title is required"
	EventTypeRequired      = "event type is required"
	ChannelNameRequired    = "channel name is required"
	ChannelIDRequired      = "channel ID is required"
	SubscriptionIDRequired = "subscription ID is required"
)

const (
	// Error messages
	Error                           = "error"
	NotAuthorized                   = "not authorized"
	GenericErrorMessage             = "Something went wrong, please try again later"
	ErrorFetchProjectList           = "Error in fetching project list"
	ErrorUnlinkProject              = "Error in unlinking the project"
	ErrorDecodingBody               = "Error in decoding body"
	GetProjectListError             = "Error getting Project List"
	ProjectNotFound                 = "Requested project does not exists"
	UnableToDisconnectUser          = "Unable to disconnect user"
	UnableToStoreOauthState         = "Unable to store oAuth state for the userID %s"
	InvalidAuthState                = "Invalid oauth state, please try again"
	UnableToCheckIfAlreadyConnected = "Unable to check if user account is already connected"
	AuthAttemptExpired              = "Authentication attempt expired, please try again"
	ErrorFetchSubscriptionList      = "Error in fetching subscription list"
	ErrorCreateSubscription         = "Error in creating subscription"
	ProjectNotLinked                = "Requested project is not linked"
	GetSubscriptionListError        = "Error getting Subscription List"
	SubscriptionAlreadyPresent      = "Subscription is already present"
	SubscriptionNotFound            = "Requested subscription does not exists"
)

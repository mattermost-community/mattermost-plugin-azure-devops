package constants

const (
	// TODO: all these messages are to be verified from Mike at the end
	// Generic
	GenericErrorMessage  = "Something went wrong, please try again later"
	ConnectAccount       = "[Click here to connect your Azure DevOps account](%s%s)"
	ConnectAccountFirst  = "Your Azure DevOps account is not connected \n%s"
	UserConnected        = "Your Azure DevOps account is successfully connected!"
	UserAlreadyConnected = "Your Azure DevOps account is already connected"
	UserDisconnected     = "Your Azure DevOps account is now disconnected"
	CreatedTask          = "Link for newly created task: %s"
	TaskTitle            = "[%s #%d: %s](%s)"
	TaskPreviewMessage   = "State: %s\nAssigned To: %s\nDescription: %s"
	AlreadyLinkedProject = "This project is already linked."
	NoProjectLinked      = "No project is linked, please link a project."

	// Validations Errors
	OrganizationRequired = "organization is required"
	ProjectRequired      = "project is required"
	TaskTypeRequired     = "task type is required"
	TaskTitleRequired    = "task title is required"
	EventTypeRequired    = "event type is required"
	ChannelIDRequired    = "channel ID is required"
)

const (
	// Error messages
	Error                           = "error"
	NotAuthorized                   = "not authorized"
	UnableToDisconnectUser          = "Unable to disconnect user"
	UnableToCheckIfAlreadyConnected = "Unable to check if user account is already connected"
	UnableToStoreOauthState         = "Unable to store oAuth state for the userID %s"
	AuthAttemptExpired              = "Authentication attempt expired, please try again"
	InvalidAuthState                = "Invalid oauth state, please try again"
	GetProjectListError             = "Error getting Project List"
	ErrorFetchProjectList           = "Error in fetching project list"
	ErrorDecodingBody               = "Error in decoding body"
	ProjectNotFound                 = "Requested project does not exists"
	ErrorUnlinkProject              = "Error in unlinking the project"
	FetchSubscriptionListError      = "Error in fetching subscription list"
	CreateSubscriptionError         = "Error in creating subscription"
	ProjectNotLinked                = "Requested project is not linked"
	GetSubscriptionListError        = "Error getting Subscription List"
	SubscriptionAlreadyPresent      = "Requested subscription already exists"
	SubscriptionNotFound            = "Requested subscription does not exists"
)

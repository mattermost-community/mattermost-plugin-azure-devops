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

	// Errors
	Error                           = "error"
	GenericErrorMessage             = "Something went wrong, please try again later"
	ErrorFetchProjectList           = "Error in fetching project list"
	ErrorUnlinkProject              = "Error in unlinking the project"
	NotAuthorized                   = "not authorized"
	ErrorDecodingBody               = "Error in decoding body"
	GetProjectListError             = "Error getting Project List"
	ProjectNotFound                 = "Requested project does not exists"
	UnableToDisconnectUser          = "Unable to disconnect user"
	UnableToStoreOauthState         = "Unable to store oAuth state for the userID %s"
	InvalidAuthState                = "Invalid oauth state, please try again"
	UnableToCheckIfAlreadyConnected = "Unable to check if user account is already connected"
	AuthAttemptExpired              = "Authentication attempt expired, please try again"

	// Validations
	OrganizationRequired = "Organization is required"
	ProjectRequired      = "Project is required"
	TaskTypeRequired     = "Task type is required"
	TaskTitleRequired    = "Task title is required"
)

export const error = {
    generic: 'Something went wrong, please try again later',

    // Forms
    errorFetchingChannelsList: 'Error occurred while fetching the channel list',
    errorFetchingOrganizationAndProjectsList: 'Error occurred while fetching the organization and project list',

    // Subscription
    subscriptionAlreadyExists: 'Requested subscription already exists',
    accessDenied: 'Access Denied',
    subscriptionNotFound: 'Requested subscription does not exist',
    adminAccessError: 'Looks like you do not have access to add/delete a subscription for this project. Please make sure you are a project or team administrator for this project',
    failedToGetSubscriptions: 'Failed to get the subscription filter values',

    // Link
    notAccessibleError: 'Looks like this project/organization does not exist or you do not have permissions to access it',
    adminAccessErrorForUnlinking: 'You do not have sufficient permissions to delete subscriptions for this project but you can still unlink the project',
    projectAlreadyLinkedError: 'This project is already linked.',
    errorExpectedForOAuthNotEnabled: 'failed to link Project: status: 401 Unauthorized',
    errorMessageOAuthNotEnabled: 'Looks like "third-party application access via OAuth" setting is not enabled for the organization or you do not have sufficient permissions to access it',
};

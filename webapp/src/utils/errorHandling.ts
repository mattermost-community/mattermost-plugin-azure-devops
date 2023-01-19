import pluginConstants from 'pluginConstants';

const getErrorMessage = (
    isError: boolean,
    component: ErrorComponents,
    errorState?: ApiErrorResponse,
): string => {
    if (!isError) {
        return '';
    }

    switch (component) {
    case 'SubscribeModal': // Create subscription modal
        if (errorState?.status === 400 && errorState?.data.Error === pluginConstants.messages.error.subscriptionAlreadyExists) {
            return pluginConstants.messages.error.subscriptionAlreadyExists;
        }
        if (errorState?.status === 403 && errorState?.data.Error.includes(pluginConstants.messages.error.accessDenied)) {
            return pluginConstants.messages.error.adminAccessError;
        }
        if (errorState?.status === 500 && errorState?.data.Error.includes(pluginConstants.messages.error.failedToGetSubscriptions)) {
            return pluginConstants.messages.error.failedToGetSubscriptions;
        }
        return errorState?.data.Error ?? pluginConstants.messages.error.generic;

    case 'LinkProjectModal':
        if (errorState?.status === 404 || errorState?.status === 401) {
            return pluginConstants.messages.error.notAccessibleError;
        }
        if (errorState?.status === 500 && errorState?.data.Error.includes(pluginConstants.messages.error.errorExpectedForOAuthNotEnabled)) {
            return pluginConstants.messages.error.errorMessageOAuthNotEnabled;
        }
        return errorState?.data.Error ?? pluginConstants.messages.error.generic;

    case 'ConfirmationModal':
        if (errorState?.status === 403 && errorState?.data.Error.includes(pluginConstants.messages.error.accessDenied)) {
            return pluginConstants.messages.error.adminAccessError;
        }
        if (errorState?.status === 404 && errorState?.data.Error.includes(pluginConstants.messages.error.subscriptionNotFound)) {
            return pluginConstants.messages.error.subscriptionNotFound;
        }
        return pluginConstants.messages.error.generic;

    default:
        return errorState?.data.Error ?? pluginConstants.messages.error.generic;
    }
};

export default getErrorMessage;

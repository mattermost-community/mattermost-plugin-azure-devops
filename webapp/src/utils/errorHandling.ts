import plugin_constants from 'plugin_constants';

const getErrorMessage = (
    isError: boolean,
    component: ErrorComponents,
    errorState: ApiErrorResponse,
): string => {
    if (!isError) {
        return '';
    }

    switch (component) {
    case 'SubscribeModal': // Create subscription modal
        if (errorState.status === 400 && errorState.data.error === plugin_constants.messages.error.subscriptionAlreadyExists) {
            return errorState.data.error;
        }
        if (errorState.status === 403 && errorState.data.error.includes(plugin_constants.messages.error.accessDenied)) {
            return plugin_constants.messages.error.adminAccessError;
        }
        return plugin_constants.messages.error.generic;

    case 'LinkProjectModal':
        if (errorState.status === 404 || errorState.status === 401) {
            return plugin_constants.messages.error.notAccessibleError;
        }
        return plugin_constants.messages.error.generic;

    case 'ConfirmationModal':
        if (errorState.status === 403 && errorState.data.error.includes(plugin_constants.messages.error.accessDenied)) {
            return plugin_constants.messages.error.adminAccessError;
        }
        return plugin_constants.messages.error.generic;

    default:
        return plugin_constants.messages.error.generic;
    }
};

export default getErrorMessage;

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
        return plugin_constants.messages.error.generic;

    default:
        return plugin_constants.messages.error.generic;
    }
};

export default getErrorMessage;

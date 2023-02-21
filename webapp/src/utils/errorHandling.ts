import pluginConstants from 'pluginConstants';

const getSubscribeModalCustomError = (errorState?: ApiErrorResponse) => {
    if (errorState?.status === 400 && errorState?.data.Error === pluginConstants.messages.error.subscriptionAlreadyExists) {
        return pluginConstants.messages.error.subscriptionAlreadyExists;
    }

    if (errorState?.status === 403 && errorState?.data.Error.includes(pluginConstants.messages.error.accessDenied)) {
        return pluginConstants.messages.error.adminAccessError;
    }

    if (errorState?.status === 500 && errorState?.data.Error.includes(pluginConstants.messages.error.failedToGetSubscriptions)) {
        return pluginConstants.messages.error.failedToGetSubscriptions;
    }

    return null;
};

const getLinkProjectModalCustomError = (errorState?: ApiErrorResponse) => {
    if (errorState?.status === 404 || errorState?.status === 401) {
        return pluginConstants.messages.error.notAccessibleError;
    }

    if (errorState?.status === 500 && errorState?.data.Error.includes(pluginConstants.messages.error.errorExpectedForOAuthNotEnabled)) {
        return pluginConstants.messages.error.errorMessageOAuthNotEnabled;
    }

    return null;
};

const confirmationModalCustomError = (errorState?: ApiErrorResponse) => {
    if (errorState?.status === 403 && errorState?.data.Error.includes(pluginConstants.messages.error.accessDenied)) {
        return pluginConstants.messages.error.adminAccessError;
    }

    if (errorState?.status === 404 && errorState?.data.Error.includes(pluginConstants.messages.error.subscriptionNotFound)) {
        return pluginConstants.messages.error.subscriptionNotFound;
    }

    return null;
};

const customErrorAndComponentMap: Partial<Record<ErrorComponents, (errorState?: ApiErrorResponse) => string | null>> = {
    SubscribeModal: (errorState?: ApiErrorResponse) => getSubscribeModalCustomError(errorState),
    LinkProjectModal: (errorState?: ApiErrorResponse) => getLinkProjectModalCustomError(errorState),
    ConfirmationModal: (errorState?: ApiErrorResponse) => confirmationModalCustomError(errorState),
};

const getErrorMessage = (
    isError: boolean,
    component: ErrorComponents,
    errorState?: ApiErrorResponse,
): string => {
    if (!isError) {
        return '';
    }

    // Return custom or API error for different components
    return customErrorAndComponentMap[component]?.(errorState) ?? errorState?.data.Error ?? pluginConstants.messages.error.generic;
};

export default getErrorMessage;

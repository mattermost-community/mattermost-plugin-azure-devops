import {useSelector, useDispatch} from 'react-redux';

import {setApiRequestCompletionState} from 'reducers/apiRequest';

import {azureDevOpsPluginApi} from 'services';

function usePluginApi() {
    const state = useSelector((reduxState: ReduxState) => reduxState);
    const dispatch = useDispatch();

    // Pass payload in POST requests only. For GET requests, there is no need to pass a payload argument
    const makeApiRequest = async (serviceName: PluginApiServiceName, payload: APIRequestPayload): Promise<any> => {
        return dispatch(azureDevOpsPluginApi.endpoints[serviceName].initiate(payload)); //TODO: add proper type here
    };

    const makeApiRequestWithCompletionStatus = async (serviceName: PluginApiServiceName, payload: APIRequestPayload) => {
        const apiRequest = await makeApiRequest(serviceName, payload);

        if (apiRequest) {
            dispatch(setApiRequestCompletionState(serviceName));
        }
    };

    // Pass payload in POST requests only. For GET requests, there is no need to pass a payload argument
    const getApiState = (serviceName: PluginApiServiceName, payload: APIRequestPayload) => {
        const {data, isError, isLoading, isSuccess, error, isUninitialized} = azureDevOpsPluginApi.endpoints[serviceName].select(payload)(state['plugins-mattermost-plugin-azure-devops']);
        return {data, isError, isLoading, isSuccess, error, isUninitialized};
    };

    return {makeApiRequest, makeApiRequestWithCompletionStatus, getApiState, state};
}

export default usePluginApi;

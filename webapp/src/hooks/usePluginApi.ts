import {useSelector, useDispatch} from 'react-redux';

import {setApiRequestCompletionState} from 'reducers/apiRequest';

import services from 'services';

function usePluginApi() {
    const state = useSelector((pluginState: ReduxState) => pluginState);
    const dispatch = useDispatch();

    // Pass payload only in POST rquests for GET requests there is no need to pass payload argument
    const makeApiRequest = async (serviceName: ApiServiceName, payload: APIRequestPayload): Promise<any> => {
        return dispatch(services.endpoints[serviceName].initiate(payload)); //TODO: add proper type here
    };

    const makeApiRequestWithCompletionStatus = async (serviceName: ApiServiceName, payload: APIRequestPayload) => {
        const apiRequest = await makeApiRequest(serviceName, payload);

        if (apiRequest) {
            dispatch(setApiRequestCompletionState(serviceName));
        }
    };

    // Pass payload only in POST rquests for GET requests there is no need to pass payload argument
    const getApiState = (serviceName: ApiServiceName, payload: APIRequestPayload) => {
        const {data, isError, isLoading, isSuccess, error, isUninitialized} = services.endpoints[serviceName].select(payload)(state['plugins-mattermost-plugin-azure-devops']);
        return {data, isError, isLoading, isSuccess, error, isUninitialized};
    };

    const isUserAccountConnected = (): boolean => {
        return state['plugins-mattermost-plugin-azure-devops'].userConnectedSlice.isConnected;
    };

    return {makeApiRequest, makeApiRequestWithCompletionStatus, getApiState, state, isUserAccountConnected};
}

export default usePluginApi;

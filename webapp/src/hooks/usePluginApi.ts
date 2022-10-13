import {useSelector, useDispatch} from 'react-redux';

import {setApiRequestCompletionState} from 'reducers/apiRequest';

import services from 'services';

function usePluginApi() {
    const state = useSelector((reduxState: ReduxState) => reduxState);
    const dispatch = useDispatch();

    // Pass payload only in POST rquests for GET requests there is no need to pass payload argument
    const makeApiRequest = async (serviceName: ApiServiceName, requestData: APIRequestData): Promise<any> => {
        return dispatch(services.endpoints[serviceName].initiate(requestData)); //TODO: add proper type here
    };

    const makeApiRequestWithCompletionStatus = async (serviceName: ApiServiceName, requestData: APIRequestData) => {
        const apiRequest = await makeApiRequest(serviceName, requestData);

        if (apiRequest) {
            dispatch(setApiRequestCompletionState(serviceName));
        }
    };

    // Pass payload only in POST rquests for GET requests there is no need to pass payload argument
    const getApiState = (serviceName: ApiServiceName, requestData: APIRequestData | void) => {
        const {data, isError, isLoading, isSuccess, error, isUninitialized} = services.endpoints[serviceName].select(requestData)(state['plugins-mattermost-plugin-azure-devops']);
        return {data, isError, isLoading, isSuccess, error, isUninitialized};
    };

    return {makeApiRequest, makeApiRequestWithCompletionStatus, getApiState, state};
}

export default usePluginApi;

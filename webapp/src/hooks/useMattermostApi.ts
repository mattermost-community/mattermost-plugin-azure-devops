import {useSelector, useDispatch} from 'react-redux';

import {mattermostServerApi} from 'services';

function useMattermostApi() {
    const state = useSelector((reduxState: ReduxState) => reduxState);
    const dispatch = useDispatch();

    // Pass payload in POST requests only. For GET requests, there is no need to pass a payload argument
    const makeMattermostApiRequest = async (serviceName: MattermostApiServiceName, payload: FetchChannelParams): Promise<any> => {
        return dispatch(mattermostServerApi.endpoints[serviceName].initiate(payload)); //TODO: add proper type here
    };

    // Pass payload in POST requests only. For GET requests, there is no need to pass a payload argument
    const getMattermostApiState = (serviceName: MattermostApiServiceName, payload: FetchChannelParams) => {
        const {data, isError, isLoading, isSuccess, error, isUninitialized} = mattermostServerApi.endpoints[serviceName].select(payload)(state['plugins-mattermost-plugin-azure-devops']);
        return {data, isError, isLoading, isSuccess, error, isUninitialized};
    };

    return {makeMattermostApiRequest, getMattermostApiState, state};
}

export default useMattermostApi;

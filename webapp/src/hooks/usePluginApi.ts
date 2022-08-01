import {useSelector, useDispatch} from 'react-redux';

import services from 'services';

function usePluginApi() {
    const state = useSelector((pluginState: PluginState) => pluginState);
    const dispatch = useDispatch();

    // Pass payload only in POST rquests for GET requests there is no need to pass payload argument
    const makeApiRequest = (serviceName: ApiServiceName, payload?: APIRequestPayload) => {
        dispatch(services.endpoints[serviceName].initiate(payload));
    };

    // Pass payload only in POST rquests for GET requests there is no need to pass payload argument
    const getApiState = (serviceName: ApiServiceName, payload?: APIRequestPayload) => {
        const {data, isError, isLoading, isSuccess} = services.endpoints[serviceName].select(payload)(state['plugins-mattermost-plugin-azure-devops']);
        return {data, isError, isLoading, isSuccess};
    };

    return {makeApiRequest, getApiState, state};
}

export default usePluginApi;

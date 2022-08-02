import {useSelector, useDispatch} from 'react-redux';
import {AnyAction} from 'redux';

import plugin_constants from 'plugin_constants';

import services from 'services';

function usePluginApi() {
    const state = useSelector((pluginState: PluginState) => pluginState);
    const dispatch = useDispatch();

    // Pass payload only in POST rquests for GET requests there is no need to pass payload argument
    const makeApiRequest = (serviceName: ApiServiceName, payload?: APIRequestPayload): Promise<AnyAction> => {
        return dispatch(services.endpoints[serviceName].initiate(payload)); //TODO: add proper type here
    };

    // Pass payload only in POST rquests for GET requests there is no need to pass payload argument
    const getApiState = (serviceName: ApiServiceName, payload?: APIRequestPayload) => {
        const {data, isError, isLoading, isSuccess} = services.endpoints[serviceName].select(payload)(state['plugins-mattermost-plugin-azure-devops']);
        return {data, isError, isLoading, isSuccess};
    };

    const getUserAccountConnectionState = () => {
        const {isLoading, isError, isSuccess, data} = getApiState(plugin_constants.pluginApiServiceConfigs.getUserDetails.apiServiceName);
        return {isLoading, data, isError, isSuccess};
    };

    return {makeApiRequest, getApiState, state, getUserAccountConnectionState};
}

export default usePluginApi;

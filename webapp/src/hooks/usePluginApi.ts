import {useSelector, useDispatch} from 'react-redux';

import services from 'services';

function usePluginApi() {
    const state = useSelector((pluginState: PluginState) => pluginState);
    const dispatch = useDispatch();

    const makeApiRequest = (apiServiceName: string) => {
        dispatch(services.endpoints[apiServiceName].initiate());
    };

    const getApiState = (apiServiceName: string) => {
        const {data, isError, isLoading, isSuccess} = services.endpoints[apiServiceName].select()(state['plugins-mattermost-plugin-wellsite-witsml']);
        return {data, isError, isLoading, isSuccess};
    };

    return {makeApiRequest, getApiState, state};
}

export default usePluginApi;

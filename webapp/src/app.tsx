import React, {useEffect} from 'react';
import {useDispatch} from 'react-redux';

import plugin_constants from 'plugin_constants';
import usePluginApi from 'hooks/usePluginApi';

import {toggleConnectionTriggered, toggleIsDisconnected} from 'reducers/userAcountDetails';
import {getLinkModalState, getTaskModalState, getRhsState, getUserConnectionState, getSubscribeModalState} from 'selectors';

// Global styles
import 'styles/main.scss';

/**
 * Mattermost plugin allows registering only one component in RHS
 * So, we would be grouping all the different components inside "Rhs" component to generate one final component for registration
 */
const App = (): JSX.Element => {
    const usePlugin = usePluginApi();
    const dispatch = useDispatch();

    useEffect(() => {
        usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.getUserDetails.apiServiceName);
    }, []);

    useEffect(() => {
        if (
            (
                getLinkModalState(usePlugin.state).visibility ||
                getTaskModalState(usePlugin.state).visibility ||
                getSubscribeModalState(usePlugin.state).visibility ||
                getRhsState(usePlugin.state).isSidebarOpen
            ) && getUserConnectionState(usePlugin.state).isConnectionTriggered
        ) {
            dispatch(toggleConnectionTriggered(false));
            usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.getUserDetails.apiServiceName);
        }
    }, [usePlugin.state]);

    useEffect(() => {
        if (getUserConnectionState(usePlugin.state).isUserDisconnected) {
            dispatch(toggleConnectionTriggered(false));
            dispatch(toggleIsDisconnected(true));
        }
    }, [usePlugin.state]);

    return (
        <></>
    );
};

export default App;
